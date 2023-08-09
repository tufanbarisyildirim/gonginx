package parser

import (
	"bufio"
	"bytes"
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
	prev       token.Token
	latest     token.Token
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
	prev := s.latest
      	s.latest = s.getNextToken()
      	s.prev = prev
        return s.latest
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
		if isLuaBlock(s.latest, s.prev) {
			s.inLuaBlock = true
		}
		return s.NewToken(token.BlockStart).Lit(string(s.read()))
	case ch == '}':
		return s.NewToken(token.BlockEnd).Lit(string(s.read()))
	case ch == '#':
		return s.scanComment()
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
			panic("unexpected end of file while scanning a string, maybe an unclosed lua code?")
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
\” – To escape “ within double quoted string.
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
			panic("unexpected end of file while scanning a string, maybe an unclosed quote?")
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
	for {
		ch := s.peek()

		//space, ;  and file ending definitely ends the keyword.
		if isSpace(ch) || isEOF(ch) || ch == ';' {
			break
		}

		//the keyword could contain a variable with like ${var}
		if ch == '{' {
			if prev == '$' {
				buf.WriteString(s.readUntil(func(r rune) bool {
					return r == '}'
				}))
				buf.WriteRune(s.read()) //consume latest '}'
			} else {
				break
			}
		}
		buf.WriteRune(s.read())
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
	return ch == ' ' || ch == '\t' || isEndOfLine(ch)
}

func isEOF(ch rune) bool {
	return ch == rune(token.EOF)
}

func isEndOfLine(ch rune) bool {
	return ch == '\r' || ch == '\n'
}

func isLuaBlock(latest, prev token.Token) bool {
        // *_by_lua_block {}
        if latest.Type == token.Keyword && strings.HasSuffix(latest.Literal, "_by_lua_block") {
               return true
        }
        // set_by_lua_block $var {}
        if (latest.Type == token.Keyword && strings.HasPrefix(latest.Literal, "$")) &&
               (prev.Type == token.Keyword && prev.Literal == "set_by_lua_block") {
               return true
        }
        return false
}

