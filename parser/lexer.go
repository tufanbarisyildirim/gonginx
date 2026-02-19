package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/tufanbarisyildirim/gonginx/parser/token"
)

// lexer is the main tokenizer
type lexer struct {
	reader     *bufio.Reader
	file       string
	line       int
	column     int
	inLuaBlock bool
	Latest     token.Token
	Err        error
}

// lex initializes a lexer from string conetnt
func lex(content string) *lexer {
	return newLexer(bytes.NewBuffer([]byte(content)))
}

// newLexer initilizes a lexer from a reader
func newLexer(r io.Reader) *lexer {
	return &lexer{
		line:   1,
		reader: bufio.NewReader(r),
	}
}

// Scan gives you next token
func (s *lexer) scan() token.Token {
	if s.Err != nil {
		s.Latest = s.NewToken(token.EOF).Lit("")
		return s.Latest
	}

	s.Latest = s.getNextToken()
	return s.Latest
}

// All scans all token and returns them as a slice
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
		if isLuaBlock(s.Latest) {
			s.inLuaBlock = true
		}
		return s.NewToken(token.BlockStart).Lit(string(s.read()))
	case ch == '}':
		return s.NewToken(token.BlockEnd).Lit(string(s.read()))
	case ch == '#':
		return s.scanComment()
	case isEndOfLine(ch):
		return s.NewToken(token.EndOfLine).Lit(string(s.read()))
	case isQuote(ch):
		return s.scanQuotedString(ch)
	default:
		return s.scanKeyword()
	}
}

// Peek returns nexr rune without consuming it
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
		ch := s.peek()
		if isEOF(ch) {
			break
		} else if until(ch) {
			// Check if this is a $ followed by a variable like ${var}
			// If so, don't break - this is part of the token
			if ch == '}' && s.maybeLookingAtVariableClose() {
				buf.WriteRune(s.read())
				continue
			}
			break
		} else {
			buf.WriteRune(s.read())
		}
	}

	return buf.String()
}

// maybeLookingAtVariableClose checks if we're possibly inside a variable reference
// This is a heuristic to determine if a closing brace is part of a variable reference
func (s *lexer) maybeLookingAtVariableClose() bool {
	// Try to peek ahead to see if this looks like a variable pattern

	// Read and immediately unread to preserve our position
	s.reader.ReadRune()
	ch := s.peek()
	s.reader.UnreadRune()

	// If the next character after } is $, 1-9, or a letter, we might be in a variable context
	return ch == '$' || (ch >= '0' && ch <= '9') ||
		(ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

// NewToken creates a new Token with its line and column
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

// TODO: support unpaired bracket in string
func (s *lexer) scanLuaCode() token.Token {
	// used to save the real line and column
	ret := s.NewToken(token.LuaCode)
	stack := make([]rune, 0, 50)
	code := strings.Builder{}

	for {
		ch := s.read()
		if ch == rune(token.EOF) {
			s.setErrOnce("unexpected end of file while scanning lua code starting at line %d, column %d", ret.Line, ret.Column)
			return s.NewToken(token.EOF).Lit("")
		}
		if ch == '#' {
			code.WriteRune(ch)
			code.WriteString(s.readUntil(isEndOfLine))
			continue
		} else if ch == '}' {
			if len(stack) == 0 {
				// the end of block
				_ = s.reader.UnreadRune()
				return ret.Lit(code.String())
			}
			// maybe it's lua table end, pop stack
			if stack[len(stack)-1] == '{' {
				stack = stack[0 : len(stack)-1]
			}
		} else if ch == '{' {
			// maybe it's lua table start, push stack
			stack = append(stack, ch)
		}
		code.WriteRune(ch)
	}
}

/*
*
\” – To escape " within double quoted string.
\\ – To escape the backslash.
\n – To add line breaks between string.
\t – To add tab space.
\r – For carriage return.
*/
func (s *lexer) scanQuotedString(delimiter rune) token.Token {
	var buf bytes.Buffer
	tok := s.NewToken(token.QuotedString)
	_, _ = buf.WriteRune(s.read()) //consume delimiter
	for {
		ch := s.read()

		if ch == rune(token.EOF) {
			s.setErrOnce("unexpected end of file while scanning quoted string starting at line %d, column %d", tok.Line, tok.Column)
			return s.NewToken(token.EOF).Lit("")
		}

		if ch == '\\' && (s.peek() == delimiter) {
			buf.WriteRune(ch)       // the backslash
			buf.WriteRune(s.read()) // the char needed escaping
			continue
		}

		_, _ = buf.WriteRune(ch)
		if ch == delimiter {
			break
		}
	}

	return tok.Lit(buf.String())
}

func (s *lexer) scanKeyword() token.Token {
	var buf bytes.Buffer
	tok := s.NewToken(token.Keyword)
	prev := s.read()
	buf.WriteRune(prev)
	inVarRef := false

	for {
		ch := s.peek()

		// Space, semicolon, and file ending definitely end the keyword
		if isSpace(ch) || isEOF(ch) || ch == ';' || isEndOfLine(ch) {
			break
		}

		// Block start character ends the keyword unless we're in a variable reference
		if ch == '{' {
			if prev == '$' {
				// Starting a ${var} variable reference
				inVarRef = true
				buf.WriteRune(s.read()) // consume '{'
			} else if !inVarRef {
				// This is a real block start, end the keyword
				break
			} else {
				// Otherwise, just consume it as part of the keyword
				buf.WriteRune(s.read())
			}
		} else if ch == '}' {
			if inVarRef {
				// End of the variable reference
				inVarRef = false
				buf.WriteRune(s.read()) // consume '}'
			} else {
				// This is a real block end, end the keyword
				break
			}
		} else {
			// Any other character is part of the keyword
			prev = s.read()
			buf.WriteRune(prev)
		}
	}

	return tok.Lit(buf.String())
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

func isSpace(ch rune) bool {
	return ch == ' ' || ch == '\t'
}

func isEOF(ch rune) bool {
	return ch == rune(token.EOF)
}

func isEndOfLine(ch rune) bool {
	return ch == '\r' || ch == '\n'
}

func isLuaBlock(t token.Token) bool {
	return t.Type == token.Keyword && strings.HasSuffix(t.Literal, "_by_lua_block")
}

func (s *lexer) setErrOnce(format string, args ...any) {
	if s.Err != nil {
		return
	}

	s.Err = fmt.Errorf(format, args...)
}
