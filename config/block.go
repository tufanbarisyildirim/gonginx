package config

import (
	"bytes"

	"github.com/tufanbarisyildirim/gonginx/dumper"
)

//Block a block statement
type Block struct {
	Statements []Statement
}

//ToString return config as string
func (b *Block) ToString(style *dumper.Style) string {
	return string(b.ToByteArray(style))
}

//ToByteArray return config as byte array
func (b *Block) ToByteArray(style *dumper.Style) []byte {
	var buf bytes.Buffer

	for i, statement := range b.Statements {
		buf.WriteString(statement.ToString(style))
		if i != len(b.Statements)-1 {
			buf.WriteString("\n")
		}
	}

	return buf.Bytes()
}
