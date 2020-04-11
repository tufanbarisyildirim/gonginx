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

	for _, statement := range b.Statements {
		buf.WriteString(statement.ToString())
		buf.WriteString("\n")
	}

	return buf.Bytes()
}

func (b *Block) contextStatement() {}
