package token

import (
	"fmt"
	"gotest.tools/v3/assert/cmp"
)

type TokenType int

const (
	EOF TokenType = iota
	EOL
	KEYWORD
	QUOTED_STRING
	OPEN_BRACE
	CLOSE_BRACE
	SEMICOLON
	COMMENT
	UNKNOWN
	REGEX
)

var (
	tokenName = map[TokenType]string{
		QUOTED_STRING: "QUOTED_STRING",
		EOF:           "EOF",
		KEYWORD:       "KEYWORD",
		OPEN_BRACE:    "OPEN_BRACE",
		CLOSE_BRACE:   "CLOSE_BRACE",
		SEMICOLON:     "SEMI_COLON",
		COMMENT:       "comment",
		UNKNOWN:       "unknown",
		REGEX:         "regex",
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

func (t Token) EqualTo(t2 interface{}, a cmp.Comparison) bool {
	return t.Type == t2.(Token).Type || t.Literal == t2.(Token).Literal
}

type Tokens []Token

func (ts Tokens) EqualTo(ts2 interface{}) bool {
	ts22 := ts2.(Tokens)
	if len(ts) != len(ts22) {
		return false
	}
	for i, t := range ts {
		if t.Type != ts22[i].Type || t.Literal != ts22[i].Literal {
			return false
		}
	}
	return true
}
