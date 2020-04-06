package token

import (
	"fmt"
)

type TokenType int

const (
	EOF TokenType = iota
	EOL
	KEYWORD
	QuotedString
	OpenBrace
	CloseBrace
	SEMICOLON
	COMMENT
	UNKNOWN
	REGEX
)

var (
	tokenName = map[TokenType]string{
		QuotedString: "QUOTED_STRING",
		EOF:          "EOF",
		KEYWORD:      "KEYWORD",
		OpenBrace:    "OPEN_BRACE",
		CloseBrace:   "CLOSE_BRACE",
		SEMICOLON:    "SEMI_COLON",
		COMMENT:      "COMMENT",
		UNKNOWN:      "UNKNOWN",
		REGEX:        "REGEX",
	}
)

func (tt TokenType) String() string {
	return tokenName[tt]
}

type Token struct {
	Type    TokenType
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
