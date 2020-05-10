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
