package config

import (
	"fmt"

	"github.com/tufanbarisyildirim/gonginx/token"
)

//Include include structure
type Include struct {
	Token       token.Token
	IncludePath string
	Config
}

func (i *Include) includeStatement() {}

//ToString returns include statement as string
func (i *Include) ToString() string {
	return fmt.Sprintf("include %s;", i.IncludePath)
}

//TokenLiteral return "include"
func (i *Include) TokenLiteral() string {
	return i.Token.Literal
}

//SaveToFile saves include to its own file
func (i *Include) SaveToFile() error {
	return i.Config.SaveToFile()
}
