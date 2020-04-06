package token

import (
	"fmt"
)

type Type int

const (
	Eof Type = iota
	Eol
	Keyword
	QuotedString
	OpenBrace
	CloseBrace
	Semicolon
	Comment
	Illegal
	Regex
)

var (
	tokenName = map[Type]string{
		QuotedString: "QuotedString",
		Eof:          "Eof",
		Keyword:      "Keyword",
		OpenBrace:    "OPEN_BRACE",
		CloseBrace:   "CLOSE_BRACE",
		Semicolon:    "SEMI_COLON",
		Comment:      "Comment",
		Illegal:      "Illegal",
		Regex:        "Regex",
	}
)

func (tt Type) String() string {
	return tokenName[tt]
}

type Token struct {
	Type    Type
	Literal string
	Line    int
	Column  int
}

func (t Token) String() string {
	return fmt.Sprintf("{Type: %s,Literal:\"%s\", Line:%d,Column:%d}", t.Type, t.Literal, t.Line, t.Column)
}

func (t Token) Lit(literal string) Token {
	t.Literal = literal
	return t
}

func (t Token) EqualTo(t2 Token) bool {
	return t.Type == t2.Type || t.Literal == t2.Literal
}

type Tokens []Token

func (ts Tokens) EqualTo(tokens Tokens) bool {
	if len(ts) != len(tokens) {
		return false
	}
	for i, t := range ts {
		if !t.EqualTo(tokens[i]) {
			return false
		}
	}
	return true
}

func (t Token) is(typ Type) bool {
	return t.Type == typ
}
