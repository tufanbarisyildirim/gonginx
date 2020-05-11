package config

import (
	"bytes"
	"sort"

	"github.com/tufanbarisyildirim/gonginx/dumper"
)

//Block a block statement
type Block struct {
	Directives []IDirective
}

//ToString return config as string
func (b *Block) ToString(style *dumper.Style) string {
	return string(b.ToByteArray(style))
}

//ToByteArray return config as byte array
func (b *Block) ToByteArray(style *dumper.Style) []byte {
	var buf bytes.Buffer

	if style.SortDirectives {
		sort.SliceStable(b.Directives, func(i, j int) bool {
			return b.Directives[i].GetName() < b.Directives[j].GetName()
		})
	}

	for i, statement := range b.Directives {
		buf.WriteString(statement.ToString(style))
		if i != len(b.Directives)-1 {
			buf.WriteString("\n")
		}
	}

	return buf.Bytes()
}

//FindDirectives find directives in block recursively
func (b *Block) FindDirectives(directiveName string) []IDirective {
	directives := make([]IDirective, 0)
	for _, directive := range b.Directives {
		if directive.GetName() == directiveName {
			directives = append(directives, directive)
		}
		if directive.GetBlock() != nil {
			directives = append(directives, directive.GetBlock().FindDirectives(directiveName)...)
		}
	}

	return directives
}
