package dumper

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/tufanbarisyildirim/gonginx/config"
)

var (
	//NoIndentStyle default style
	NoIndentStyle = &Style{
		SortDirectives: false,
		StartIndent:    0,
		Indent:         0,
		Debug:          false,
	}

	//IndentedStyle default style
	IndentedStyle = &Style{
		SortDirectives: false,
		StartIndent:    0,
		Indent:         4,
		Debug:          false,
	}

	//NoIndentSortedStyle default style
	NoIndentSortedStyle = &Style{
		SortDirectives: true,
		StartIndent:    0,
		Indent:         0,
		Debug:          false,
	}

	//NoIndentSortedSpaceStyle default style
	NoIndentSortedSpaceStyle = &Style{
		SortDirectives:    true,
		SpaceBeforeBlocks: true,
		StartIndent:       0,
		Indent:            0,
		Debug:             false,
	}
)

// Style dumping style
type Style struct {
	SortDirectives    bool
	SpaceBeforeBlocks bool
	StartIndent       int
	Indent            int
	Debug             bool
}

// NewStyle create new style
func NewStyle() *Style {
	style := &Style{
		SortDirectives: false,
		StartIndent:    0,
		Indent:         4,
		Debug:          false,
	}
	return style
}

// Iterate interate the indentation for sub blocks
func (s *Style) Iterate() *Style {
	newStyle := &Style{
		SortDirectives:    s.SortDirectives,
		SpaceBeforeBlocks: s.SpaceBeforeBlocks,
		StartIndent:       s.StartIndent + s.Indent,
		Indent:            s.Indent,
	}
	return newStyle
}

// DumpDirective convert a directive to a string
func DumpDirective(d config.IDirective, style *Style) string {
	if d == nil {
		return ""
	}

	var buf bytes.Buffer

	if style.SpaceBeforeBlocks && d.GetBlock() != nil {
		buf.WriteString("\n")
	}
	// outline comment
	if len(d.GetComment()) > 0 {
		for _, comment := range d.GetComment() {
			buf.WriteString(fmt.Sprintf("%s%s\n", strings.Repeat(" ", style.StartIndent), comment))
		}
	}
	buf.WriteString(fmt.Sprintf("%s%s", strings.Repeat(" ", style.StartIndent), d.GetName()))

	inlineComments := make(map[int]config.InlineComment)
	for _, comment := range d.GetInlineComment() {
		inlineComments[comment.RelativeLineIndex] = comment
	}

	// Use relative line index to handle different line number arrangements of instruction parameters
	relativeLineIndex := 0
	for _, parameter := range d.GetParameters() {
		// If the parameter line index is not the same as the previous one, add a new line
		if parameter.GetRelativeLineIndex() != relativeLineIndex {
			// write param comment
			if comment, ok := inlineComments[relativeLineIndex]; ok {
				buf.WriteString(fmt.Sprintf(" %s\n", comment.Value))
			} else {
				buf.WriteString("\n")
			}
			buf.WriteString(fmt.Sprintf("%s%s", strings.Repeat(" ", style.StartIndent+style.Indent), parameter.GetValue()))
			relativeLineIndex = parameter.GetRelativeLineIndex()
		} else {
			buf.WriteString(fmt.Sprintf(" %s", parameter.GetValue()))
		}
	}

	if d.GetBlock() == nil {
		if d.GetName() != "" {
			buf.WriteRune(';')
		}
		// the last inline comment
		if comment, ok := inlineComments[relativeLineIndex]; ok {
			buf.WriteString(comment.Value)
		}
	} else {
		buf.WriteString(" {\n")
		buf.WriteString(DumpBlock(d.GetBlock(), style.Iterate()))
		buf.WriteString(fmt.Sprintf("\n%s}", strings.Repeat(" ", style.StartIndent)))
	}
	return buf.String()
}

// DumpBlock convert a directive to a string
func DumpBlock(b config.IBlock, style *Style) string {
	if b.GetCodeBlock() != "" {
		return DumpLuaBlock(b, style)
	}

	var buf bytes.Buffer
	directives := append([]config.IDirective(nil), b.GetDirectives()...)
	if style.SortDirectives {
		sort.SliceStable(directives, func(i, j int) bool {
			return directives[i].GetName() < directives[j].GetName()
		})
	}

	for i, directive := range directives {
		if style.Debug {
			buf.WriteString("#")
			buf.WriteString(directive.GetName())
			buf.WriteString(fmt.Sprintf("%t", directive.GetBlock()))
			buf.WriteString("\n")
		}
		buf.WriteString(DumpDirective(directive, style))
		if i != len(directives)-1 {
			buf.WriteString("\n")
		}
	}

	return buf.String()
}

// DumpConfig dump whole config
func DumpConfig(c *config.Config, style *Style) string {
	return DumpBlock(c.Block, style)
}

// DumpInclude dump(stringify) the included AST
func DumpInclude(i *config.Include, style *Style) map[string]string {
	mp := make(map[string]string)
	for _, cfg := range i.Configs {
		mp[cfg.FilePath] = DumpConfig(cfg, style)
	}
	return mp
}

// WriteConfig writes config
func WriteConfig(c *config.Config, style *Style, writeInclude bool) error {
	if writeInclude {
		includes := c.FindDirectives("include")
		for _, include := range includes {
			i, ok := include.(*config.Include)
			if !ok {
				return fmt.Errorf("include directive type mismatch: %T", include)
			}

			// no config parsed
			if len(i.Configs) == 0 {
				continue
			}

			mp := DumpInclude(i, style)
			for path, config := range mp {
				// create parent directories, if not exit
				dir, _ := filepath.Split(path)
				err := os.MkdirAll(dir, 0755)
				if err != nil {
					return err
				}
				err = os.WriteFile(path, []byte(config), 0644)
				if err != nil {
					return err
				}
			}
		}
	}
	// create parent directories, if not exit
	dir, _ := filepath.Split(c.FilePath)
	if dir != "" {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return os.WriteFile(c.FilePath, []byte(DumpConfig(c, style)), 0644)
}
