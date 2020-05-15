package gonginx

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

)

var (
	//NoIndentStyle default style
	NoIndentStyle = &Style{
		SortDirectives: false,
		StartIndent:    0,
		Indent:         0,
	}

	//IndentedStyle default style
	IndentedStyle = &Style{
		SortDirectives: false,
		StartIndent:    0,
		Indent:         4,
	}

	//NoIndentSortedStyle default style
	NoIndentSortedStyle = &Style{
		SortDirectives: true,
		StartIndent:    0,
		Indent:         0,
	}

	//NoIndentSortedSpaceStyle default style
	NoIndentSortedSpaceStyle = &Style{
		SortDirectives:    true,
		SpaceBeforeBlocks: true,
		StartIndent:       0,
		Indent:            0,
	}
)

//Style dumping style
type Style struct {
	SortDirectives    bool
	SpaceBeforeBlocks bool
	StartIndent       int
	Indent            int
}

//NewStyle create new style
func NewStyle() *Style {
	style := &Style{
		SortDirectives: false,
		StartIndent:    0,
		Indent:         4,
	}
	return style
}

//Iterate interate the indentation for sub blocks
func (s *Style) Iterate() *Style {
	newStyle := &Style{
		SortDirectives:    s.SortDirectives,
		SpaceBeforeBlocks: s.SpaceBeforeBlocks,
		StartIndent:       s.StartIndent + s.Indent,
		Indent:            s.Indent,
	}
	return newStyle
}


//DumpDirective convert a directive to a string
func DumpDirective(d IDirective, style *Style) string {
	var buf bytes.Buffer

	if style.SpaceBeforeBlocks && d.GetBlock() != nil {
		buf.WriteString("\n")
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

//DumpBlock convert a directive to a string
func DumpBlock(b IBlock, style *Style) string {
	var buf bytes.Buffer

	directives := b.GetDirectives()
	if style.SortDirectives {
		sort.SliceStable(directives, func(i, j int) bool {
			return directives[i].GetName() < directives[j].GetName()
		})
	}

	for i, directive := range directives {
		buf.WriteString(DumpDirective(directive, style))
		if i != len(directives)-1 {
			buf.WriteString("\n")
		}
	}

	return string(buf.Bytes())
}
