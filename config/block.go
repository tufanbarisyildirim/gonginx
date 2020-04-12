package config

import (
	"bytes"
)

//Block a block statement
type Block struct {
	Statements []Statement
}

//ToString return config as string
func (b *Block) ToString() string {
	return string(b.ToByteArray())
}

//ToByteArray return config as byte array
func (b *Block) ToByteArray() []byte {
	var buf bytes.Buffer

	for i, statement := range b.Statements {
		buf.WriteString(statement.ToString())
		if i != len(b.Statements)-1 {
			buf.WriteString("\n")
		}
	}

	return buf.Bytes()
}
