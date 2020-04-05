package parser

import (
	"bufio"
	"bytes"
	"fmt"
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
)

var (
	tokenName = map[TokenType]string{
		quotedString: "quoted string",
		eof:          "end of file",
		keyword:      "keyword",
		openBrace:    "open brace",
		closeBrace:   "close brace",
		semiColon:    "semi colon",
		comment:      "comment",
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
	return fmt.Sprintf("{%s:%s}", t.Type, t.Literal)
}

var (
	tokenEof        = Token{Type: eof}
	tokenEol        = Token{Type: eol}
	tokenOpenBrace  = Token{Type: openBrace, Literal: "{"}
	tokenCloseBrace = Token{Type: closeBrace, Literal: "}"}
	tokenSemiColon  = Token{Type: semiColon, Literal: ";"}
)

type Scanner struct {
	reader *bufio.Reader
	line   int
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
		return tokenEof
	case ch == ';':
		s.read()
		return tokenSemiColon
	case ch == '{':
		s.read()
		return tokenOpenBrace
	case ch == '}':
		s.read()
		return tokenCloseBrace
	case ch == '#':
		return s.scanComment()
	case isLetter(ch):
		return s.scanKeyword()
	case isQuote(ch):
		return s.scanQuotedString(ch)
	}

	return Token{
		Type:    TokenType(ch),
		Literal: string(s.read()),
	}
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
	return Token{
		Type:    comment,
		Literal: s.readUntil(isEndOfLine),
		Line:    s.line,
	}
}

/**
\” – To escape “ within double quoted string.
\\ – To escape the backslash.
\n – To add line breaks between string.
\t – To add tab space.
\r – For carriage return.
*/
func (s *Scanner) scanQuotedString(delimiter rune) (tok Token) {
	var buf bytes.Buffer
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

	tok = Token{
		Type:    quotedString,
		Literal: buf.String(),
	}
	return tok
}

func (s *Scanner) scanKeyword() (tok Token) {
	word := s.readUntil(isKeywordTerminator)
	return Token{
		Type:    keyword,
		Literal: word,
	}
}

func (s *Scanner) unread() {
	_ = s.reader.UnreadRune()
}

func (s *Scanner) read() rune {
	ch, _, err := s.reader.ReadRune()
	if err != nil {
		return rune(eof)
	}
	if ch == '\n' {
		s.line++
	}
	return ch
}

func isQuote(ch rune) bool {
	return ch == '"' || ch == '\'' || ch == '`'
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
