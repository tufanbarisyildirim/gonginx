package parser

import (
	"bufio"
	"os"

	"github.com/tufanbarisyildirim/gonginx/config"
	"github.com/tufanbarisyildirim/gonginx/token"
)

//Parser is an nginx config parser
type Parser struct {
	lexer          *lexer
	currentToken   token.Token
	followingToken token.Token
}

//NewStringParser parses nginx conf from string
func NewStringParser(str string) *Parser {
	return NewParserFromLexer(lex(str))
}

//NewParser create new parser
func NewParser(filePath string) (*Parser, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return NewParserFromLexer(newLexer(bufio.NewReader(f))), nil
}

//NewParserFromLexer initilizes a new Parser
func NewParserFromLexer(lexer *lexer) *Parser {
	parser := &Parser{
		lexer: lexer,
	}
	parser.nextToken()
	parser.nextToken()

	return parser
}

func (p *Parser) nextToken() {
	p.currentToken = p.followingToken
	p.followingToken = p.lexer.scan()
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.currentToken.Type == t
}

func (p *Parser) followingTokenIs(t token.Type) bool {
	return p.followingToken.Type == t
}

//ParseBlock parse a block statement
func (p *Parser) ParseBlock(name string) *config.Block {
	if name == "" {
		name = "main"
	}
	context := &config.Block{
		Context:    name,
		Statements: make([]config.Statement, 0),
	}

parsingloop:
	for {
		switch {
		case p.curTokenIs(token.Eof):
			break parsingloop
		case p.curTokenIs(token.Keyword):
			if p.followingTokenIs(token.BlockStart) {
				context.Statements = append(context.Statements, p.ParseBlock(p.currentToken.Literal))
			} else {
				context.Statements = append(context.Statements, p.parseDirective())
			}
			break
		}
	}

	return context
}

func (p *Parser) parseInclude() *config.Include {
	include := &config.Include{}
	p.nextToken() //path
	include.IncludePath = p.currentToken.Literal
	// path
	if !p.curTokenIs(token.Semicolon) {
		panic("expected semicolon after include path")
	}

	return include
}

func (p *Parser) parseDirective() *config.Directive {
	d := &config.Directive{
		Name: p.currentToken.Literal,
	}
	for p.nextToken(); !p.curTokenIs(token.Eof) && !p.curTokenIs(token.Semicolon); p.nextToken() {
		d.Parameters = append(d.Parameters, p.currentToken.Literal)
	}

	if !p.curTokenIs(token.Semicolon) {
		panic("expected semicolon after include path")
	}

	return d
}
