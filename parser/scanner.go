package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"gotest.tools/v3/assert/cmp"
	"io"
	"unicode"
)

type TokenType int

const (
	eof TokenType = iota
	eol
	keyword
	quotedString
	openBrace
	closeBrace
	semiColon
	comment
	symbol
	regex
)

var (
	tokenName = map[TokenType]string{
		quotedString: "quotedString",
		eof:          "eof",
		keyword:      "keyword",
		openBrace:    "openBrace",
		closeBrace:   "closeBrace",
		semiColon:    "semiColon",
		comment:      "comment",
		symbol:       "symbol",
		regex:        "regex",
	}
)

func (tt TokenType) String() string {
	return tokenName[tt]
}

type Token struct {
	Type      TokenType
	Literal   string
	Line      int
	LineEnd   int
	Column    int
	ColumnEnd int
}

func (t Token) String() string {
	return fmt.Sprintf("{Type: %s,Literal:\"%s\", Line:%d,LineEnd:%d,Column:%d,ColumnEnd:%d}", t.Type, t.Literal, t.Line, t.LineEnd, t.Column, t.ColumnEnd)
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
	ts22 := ts2.([]Token)
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

type Scanner struct {
	reader *bufio.Reader
	file   string
	line   int
	column int
	Latest Token
}

func Parse(content string) *Scanner {
	return NewScanner(bytes.NewBuffer([]byte(content)))
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		reader: bufio.NewReader(r),
	}
}

func (s *Scanner) Scan() Token {
	s.Latest = s.getNextToken()
	return s.Latest
}

func (s *Scanner) All() []Token {
	tokens := make([]Token, 0)
	for {
		v := s.Scan()
		if v.Type == eof || v.Type == -1 {
			break
		}
		tokens = append(tokens, v)
	}
	return tokens
}

func (s *Scanner) getNextToken() Token {
reToken:
	ch := s.Peek()
	switch {
	case isSpace(ch):
		s.skipWhitespace()
		goto reToken
	case isEof(ch):
		return s.NewToken(eof).Lit(string(s.read()))
	case ch == ';':
		return s.NewToken(semiColon).Lit(string(s.read()))
	case ch == '{':
		return s.NewToken(openBrace).Lit(string(s.read()))
	case ch == '}':
		return s.NewToken(closeBrace).Lit(string(s.read()))
	case ch == '#':
		return s.scanComment()
	case isQuote(ch):
		return s.scanQuotedString(ch)
	case isNotSpace(ch):
		return s.scanKeyword()
	}

	return s.NewToken(symbol).Lit(string(s.read()))
}

func (s *Scanner) Peek() rune {
	r := s.read()
	s.unread()
	return r
}

func (s *Scanner) PeekPrev() rune {
	s.unread()
	r := s.read()
	return r
}

type runeCheck func(rune) bool

func (s *Scanner) readUntil(until runeCheck) string {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); isEof(ch) {
			break
		} else if until(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return buf.String()
}

func (s *Scanner) NewToken(tokenType TokenType) Token {
	return Token{
		Type:   tokenType,
		Line:   s.line,
		Column: s.column,
	}
}

func (s *Scanner) readUntilWith(until runeCheck) string {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); isEof(ch) {
			break
		} else if until(ch) {
			buf.WriteRune(ch)
			break
		} else {
			buf.WriteRune(ch)
		}
	}
	return buf.String()
}

func (s *Scanner) readWhile(while runeCheck) string {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		ch := s.read()
		if while(ch) {
			buf.WriteRune(ch)
		} else {
			s.unread()
			break
		}
	}
	// unread the latest char we consume.
	return buf.String()
}

func (s *Scanner) skipWhitespace() {
	s.readWhile(isSpace)
}

func (s *Scanner) skipEndOfLine() {
	s.readUntilWith(isEndOfLine)
}

func (s *Scanner) scanComment() Token {
	return s.NewToken(comment).Lit(s.readUntil(isEndOfLine))
}

func (s *Scanner) scanRegex() Token {
	return s.NewToken(regex).Lit(s.readUntil(isSpace))
}

/**
\” – To escape “ within double quoted string.
\\ – To escape the backslash.
\n – To add line breaks between string.
\t – To add tab space.
\r – For carriage return.
*/
func (s *Scanner) scanQuotedString(delimiter rune) Token {
	var buf bytes.Buffer
	tok := s.NewToken(quotedString)
	s.read() //consume delimiter
	for {
		ch := s.read()

		if ch == rune(eof) {
			panic("unexpected end of file while scanning a string, maybe an unclosed quote?")
		}

		if ch == '\\' {
			if needsEscape(s.Peek(), delimiter) {
				switch s.read() {
				case 'n':
					buf.WriteRune('\n')
				case 'r':
					buf.WriteRune('\r')
				case 't':
					buf.WriteRune('\t')
				case '\\':
					buf.WriteRune('\\')
				case delimiter:
					buf.WriteRune(delimiter)
				}
				continue
			}
		}
		if ch == delimiter {
			break
		}
		buf.WriteRune(ch)
	}

	return tok.Lit(buf.String())
}

func (s *Scanner) scanKeyword() Token {
	return s.NewToken(keyword).Lit(s.readUntil(isKeywordTerminator))
}

func (s *Scanner) unread() {
	_ = s.reader.UnreadRune()
	s.column--
}

func (s *Scanner) read() rune {
	ch, _, err := s.reader.ReadRune()
	if err != nil {
		return rune(eof)
	}

	if ch == '\n' {
		s.column = 1
		s.line++
	} else {
		s.column++
	}
	return ch
}

func isQuote(ch rune) bool {
	return ch == '"' || ch == '\'' || ch == '`'
}

func isRegexDelimiter(ch rune) bool {
	return ch == '/'
}

func isNotSpace(ch rune) bool {
	return !isSpace(ch)
}

func isKeywordTerminator(ch rune) bool {
	return isSpace(ch) || isEndOfLine(ch) || ch == '{' || ch == ';'
}

func needsEscape(ch, delimiter rune) bool {
	return ch == delimiter || ch == 'n' || ch == 't' || ch == '\\' || ch == 'r'
}

func isSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || isEndOfLine(ch)
}

func isEof(ch rune) bool {
	return ch == rune(eof)
}

func isEndOfLine(ch rune) bool {
	return ch == '\r' || ch == '\n'
}

func isLetter(ch rune) bool {
	return ch == '_' || unicode.IsLetter(ch)
}

func isWordStart(ch rune) bool {
	return isLetter(ch) || unicode.IsDigit(ch)
}
