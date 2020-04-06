package parser

import (
	"bufio"
	"bytes"
	"github.com/tufanbarisyildirim/gonginx/token"
	"io"
	"unicode"
)

type Scanner struct {
	reader *bufio.Reader
	file   string
	line   int
	column int
	Latest token.Token
}

func Parse(content string) *Scanner {
	return NewScanner(bytes.NewBuffer([]byte(content)))
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		reader: bufio.NewReader(r),
	}
}

func (s *Scanner) Scan() token.Token {
	s.Latest = s.getNextToken()
	return s.Latest
}

func (s *Scanner) All() token.Tokens {
	tokens := make([]token.Token, 0)
	for {
		v := s.Scan()
		if v.Type == token.EOF || v.Type == -1 {
			break
		}
		tokens = append(tokens, v)
	}
	return tokens
}

func (s *Scanner) getNextToken() token.Token {
reToken:
	ch := s.Peek()
	switch {
	case isSpace(ch):
		s.skipWhitespace()
		goto reToken
	case isEof(ch):
		return s.NewToken(token.EOF).Lit(string(s.read()))
	case ch == ';':
		return s.NewToken(token.SEMICOLON).Lit(string(s.read()))
	case ch == '{':
		return s.NewToken(token.OpenBrace).Lit(string(s.read()))
	case ch == '}':
		return s.NewToken(token.CloseBrace).Lit(string(s.read()))
	case ch == '#':
		return s.scanComment()
	case isQuote(ch):
		return s.scanQuotedString(ch)
	case isNotSpace(ch):
		return s.scanKeyword()
	}

	return s.NewToken(token.UNKNOWN).Lit(string(s.read())) //that should never happen :)
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

func (s *Scanner) NewToken(tokenType token.TokenType) token.Token {
	return token.Token{
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

func (s *Scanner) scanComment() token.Token {
	return s.NewToken(token.COMMENT).Lit(s.readUntil(isEndOfLine))
}

func (s *Scanner) scanRegex() token.Token {
	return s.NewToken(token.REGEX).Lit(s.readUntil(isSpace))
}

/**
\” – To escape “ within double quoted string.
\\ – To escape the backslash.
\n – To add line breaks between string.
\t – To add tab space.
\r – For carriage return.
*/
func (s *Scanner) scanQuotedString(delimiter rune) token.Token {
	var buf bytes.Buffer
	tok := s.NewToken(token.QuotedString)
	s.read() //consume delimiter
	for {
		ch := s.read()

		if ch == rune(token.EOF) {
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

func (s *Scanner) scanKeyword() token.Token {
	return s.NewToken(token.KEYWORD).Lit(s.readUntil(isKeywordTerminator))
}

func (s *Scanner) unread() {
	_ = s.reader.UnreadRune()
	s.column--
}

func (s *Scanner) read() rune {
	ch, _, err := s.reader.ReadRune()
	if err != nil {
		return rune(token.EOF)
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
	return ch == rune(token.EOF)
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
