package parser

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"github.com/tufanbarisyildirim/gonginx/parser/token"
)

//lexer is the main tokenizer
type lexer struct {
	reader     *bufio.Reader
	file       string
	line       int
	column     int
	inLuaBlock bool
	Latest     token.Token
}

//lex initializes a lexer from string conetnt
func lex(content string) *lexer {
	return newLexer(bytes.NewBuffer([]byte(content)))
}

//newLexer initilizes a lexer from a reader
func newLexer(r io.Reader) *lexer {
	return &lexer{
		line:   1,
		reader: bufio.NewReader(r),
	}
}

//Scan gives you next token
func (s *lexer) scan() token.Token {
	s.Latest = s.getNextToken()
	return s.Latest
}

//All scans all token and returns them as a slice
func (s *lexer) all() token.Tokens {
	tokens := make([]token.Token, 0)
	for {
		v := s.scan()
		if v.Type == token.EOF || v.Type == -1 {
			break
		}
		tokens = append(tokens, v)
	}
	return tokens
}

func (s *lexer) getNextToken() token.Token {
	if s.inLuaBlock {
		s.inLuaBlock = false
		return s.scanLuaCode()
	}
reToken:
	ch := s.peek()
	switch {
	case isSpace(ch):
		s.skipWhitespace()
		goto reToken
	case isEOF(ch):
		return s.NewToken(token.EOF).Lit(string(s.read()))
	case ch == ';':
		return s.NewToken(token.Semicolon).Lit(string(s.read()))
	case ch == '{':
		if s.Latest.Type == token.Keyword && strings.HasSuffix(s.Latest.Literal, "_by_lua_block") {
			s.inLuaBlock = true
		}
		return s.NewToken(token.BlockStart).Lit(string(s.read()))
	case ch == '}':
		return s.NewToken(token.BlockEnd).Lit(string(s.read()))
	case ch == '#':
		return s.scanComment()
	case ch == '$':
		return s.scanVariable()
	case isQuote(ch):
		return s.scanQuotedString(ch)
	default:
		return s.scanKeyword()
	}
}

//Peek returns nexr rune without consuming it
func (s *lexer) peek() rune {
	r, _, _ := s.reader.ReadRune()
	_ = s.reader.UnreadRune()
	return r
}

type runeCheck func(rune) bool

func (s *lexer) readUntil(until runeCheck) string {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.peek(); isEOF(ch) {
			break
		} else if until(ch) {
			break
		} else {
			buf.WriteRune(s.read())
		}
	}

	return buf.String()
}

//NewToken creates a new Token with its line and column
func (s *lexer) NewToken(tokenType token.Type) token.Token {
	return token.Token{
		Type:   tokenType,
		Line:   s.line,
		Column: s.column,
	}
}

func (s *lexer) readWhile(while runeCheck) string {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.peek(); while(ch) {
			buf.WriteRune(s.read())
		} else {
			break
		}
	}
	// unread the latest char we consume.
	return buf.String()
}

func (s *lexer) skipWhitespace() {
	s.readWhile(isSpace)
}

func (s *lexer) scanComment() token.Token {
	return s.NewToken(token.Comment).Lit(s.readUntil(isEndOfLine))
}

// TODO: support unpaired bracket in comment and string
func (s *lexer) scanLuaCode() token.Token {
	// used to save the real line and column
	ret := s.NewToken(token.Keyword)
	stack := make([]rune, 0, 50)
	code := make([]rune, 0, 100)

	for {
		ch := s.read()
		if ch == rune(token.EOF) {
			panic("unexpected end of file while scanning a string, maybe an unclosed lua code?")
		}

		if ch == '}' {
			if len(stack) == 0 {
				// the end of block
				_ = s.reader.UnreadRune()
				return ret.Lit(string(code))
			}
			// maybe it's lua table end, pop stack
			if stack[len(stack)-1] == '{' {
				stack = stack[0 : len(stack)-1]
			}
		} else if ch == '{' {
			// maybe it's lua table start, push stack
			stack = append(stack, ch)
		}
		code = append(code, ch)
	}
}

/**
\” – To escape “ within double quoted string.
\\ – To escape the backslash.
\n – To add line breaks between string.
\t – To add tab space.
\r – For carriage return.
*/
func (s *lexer) scanQuotedString(delimiter rune) token.Token {
	var buf bytes.Buffer
	tok := s.NewToken(token.QuotedString)
	buf.WriteRune(s.read()) //consume delimiter
	for {
		ch := s.read()

		if ch == rune(token.EOF) {
			panic("unexpected end of file while scanning a string, maybe an unclosed quote?")
		}

		if ch == '\\' {
			if needsEscape(s.peek(), delimiter) {
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
		buf.WriteRune(ch)
		if ch == delimiter {
			break
		}
	}

	return tok.Lit(buf.String())
}

func (s *lexer) scanKeyword() token.Token {
	return s.NewToken(token.Keyword).Lit(s.readUntil(isKeywordTerminator))
}

func (s *lexer) scanVariable() token.Token {
	return s.NewToken(token.Variable).Lit(s.readUntil(isKeywordTerminator))
}

func (s *lexer) read() rune {
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

func isKeywordTerminator(ch rune) bool {
	return isSpace(ch) || isEndOfLine(ch) || ch == '{' || ch == ';'
}

func needsEscape(ch, delimiter rune) bool {
	return ch == delimiter || ch == 'n' || ch == 't' || ch == '\\' || ch == 'r'
}

func isSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || isEndOfLine(ch)
}

func isEOF(ch rune) bool {
	return ch == rune(token.EOF)
}

func isEndOfLine(ch rune) bool {
	return ch == '\r' || ch == '\n'
}
