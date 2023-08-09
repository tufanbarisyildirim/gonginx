package gonginx

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
func DumpDirective(d IDirective, style *Style) string {
	if d == nil {
		return ""
	}

	var buf bytes.Buffer

	if style.SpaceBeforeBlocks && d.GetBlock() != nil {
		buf.WriteString("\n")
	}
	if len(d.GetComment()) > 0 {
		for _, comment := range d.GetComment() {
			buf.WriteString(fmt.Sprintf("%s%s\n", strings.Repeat(" ", style.StartIndent), comment))
		}
	}
	buf.WriteString(fmt.Sprintf("%s%s", strings.Repeat(" ", style.StartIndent), d.GetName()))
	if len(d.GetParameters()) > 0 {
		buf.WriteString(fmt.Sprintf(" %s", strings.Join(d.GetParameters(), " ")))
	}
	if d.GetBlock() == nil {
		buf.WriteRune(';')
	} else {
		buf.WriteString(" {\n")
		buf.WriteString(DumpBlock(d.GetBlock(), style.Iterate()))
		buf.WriteString(fmt.Sprintf("\n%s}", strings.Repeat(" ", style.StartIndent)))
	}
	return buf.String()
}

// DumpBlock convert a directive to a string
func DumpBlock(b IBlock, style *Style) string {
	var buf bytes.Buffer

	if b.GetCodeBlock() != "" {
		luaLines := strings.Split(b.GetCodeBlock(), "\n")
		for i, line := range luaLines {
			buf.WriteString(fmt.Sprintf("%s%s", strings.Repeat(" ", style.StartIndent), line))
			if i != len(luaLines)-1 {
				buf.WriteString("\n")
			}
		}

		return buf.String()
	}

	directives := b.GetDirectives()
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
func DumpConfig(c *Config, style *Style) string {
	return DumpBlock(c.Block, style)
}

// DumpInclude dump(stringify) the included AST
func DumpInclude(i *Include, style *Style) map[string]string {
	mp := make(map[string]string)
	for _, cfg := range i.Configs {
		mp[cfg.FilePath] = DumpConfig(cfg, style)
	}
	return mp
}

// WriteConfig writes config
func WriteConfig(c *Config, style *Style, writeInclude bool) error {
	if writeInclude {
		includes := c.FindDirectives("include")
		for _, include := range includes {
			i, ok := include.(*Include)
			if !ok {
				panic("bug in FindDirective")
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
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	return os.WriteFile(c.FilePath, []byte(DumpConfig(c, style)), 0644)
}
