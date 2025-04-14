package token

import (
	"fmt"
)

// Type Token.Type
type Type int

const (
	//EOF end of file
	EOF Type = iota
	//Eol end of line
	Eol
	//Keyword any keyword
	Keyword
	//QuotedString Quoted String
	QuotedString
	//Variable any $variabl
	Variable
	//BlockStart {
	BlockStart
	//BlockEnd }
	BlockEnd
	//Semicolon ;
	Semicolon
	//Comment #comment
	Comment
	//EndOfLine \n or \r
	EndOfLine
	//Illegal a token that should never happen
	Illegal
	//Regex any reg expression
	Regex
	// LuaCode lua block
	LuaCode
)

var (
	tokenName = map[Type]string{
		QuotedString: "QuotedString",
		EOF:          "Eof",
		Keyword:      "Keyword",
		Variable:     "Variable",
		BlockStart:   "BlockStart",
		BlockEnd:     "BlockEnd",
		Semicolon:    "Semicolon",
		Comment:      "Comment",
		EndOfLine:    "EndOfLine",
		Illegal:      "Illegal",
		Regex:        "Regex",
	}
)

// String convert a token to string as it should be written
func (tt Type) String() string {
	return tokenName[tt]
}

// Token represents a config token
type Token struct {
	Type    Type
	Literal string
	Line    int
	Column  int
}

func (t Token) String() string {
	return fmt.Sprintf("{Type:%s,Literal:\"%s\",Line:%d,Column:%d}", t.Type, t.Literal, t.Line, t.Column)
}

// Lit set literal string
func (t Token) Lit(literal string) Token {
	t.Literal = literal
	return t
}

// EqualTo checks equality
func (t Token) EqualTo(t2 Token) bool {
	return t.Type == t2.Type && t.Literal == t2.Literal
}

// Tokens list of token
type Tokens []Token

// EqualTo check Tokens equality of token list
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

// Diff find what is difference between two Token collection
func (ts Tokens) Diff(tokens Tokens) error {
	if len(ts) != len(tokens) {
		return fmt.Errorf("different token count %d vs %d", len(ts), len(tokens))
	}
	for i, t := range ts {
		if t.Type != tokens[i].Type {
			return fmt.Errorf("i=%d,  Type[%s]!=Type[%s]", i, t.Type.String(), tokens[i].Type.String())
		}

		if !t.EqualTo(tokens[i]) {
			return fmt.Errorf("tokens are not equal: i=%d, %v!=%v", i, t, tokens[i])
		}
	}
	return nil
}

// Is check type of a token
func (t Token) Is(typ Type) bool {
	return t.Type == typ
}

// IsParameterEligible checks if token is directive parameter eligible
func (t Token) IsParameterEligible() bool {
	return t.Is(Keyword) || t.Is(QuotedString) || t.Is(Variable) || t.Is(Regex)
}
