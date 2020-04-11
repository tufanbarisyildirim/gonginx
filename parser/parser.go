package parser

import (
	"bufio"
	"fmt"
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

//Parse the config.
func (p *Parser) Parse() *config.Config {
	return &config.Config{
		FilePath: "nil", //TODO: set filepath here,
		Block:    p.parseBlock(),
	}
}

//ParseBlock parse a block statement
func (p *Parser) parseBlock() *config.Block {

	context := &config.Block{
		Statements: make([]config.Statement, 0),
	}

parsingloop:
	for {
		switch {
		case p.curTokenIs(token.Eof):
			break parsingloop
		case p.curTokenIs(token.Keyword):
			context.Statements = append(context.Statements, p.parseStatement())
			break
		}
		p.nextToken()
	}

	return context
}

func (p *Parser) parseInclude() *config.Include {
	include := &config.Include{}
	p.nextToken() //include
	include.IncludePath = p.currentToken.Literal
	p.nextToken() //path

	// path
	if !p.curTokenIs(token.Semicolon) {
		panic(fmt.Errorf("expected semicolon after include path but got %s", p.currentToken.Literal))
	}

	//TODO: start sub parsing here, detect all files from include path
	//		support wildcards, include all matching files

	return include
}

func (p *Parser) parseStatement() config.Statement {
	d := &config.Directive{
		Name: p.currentToken.Literal,
	}

	switch d.Name {
	case "include":
		return p.parseInclude()
	}

	for p.nextToken(); !p.curTokenIs(token.Eof) && p.curTokenIs(token.Keyword); p.nextToken() {
		d.Parameters = append(d.Parameters, p.currentToken.Literal)
	}

	if p.curTokenIs(token.Semicolon) {
		return d
	}

	if p.curTokenIs(token.BlockStart) {
		d.Block = p.parseBlock()
		return d
	}

	panic(fmt.Errorf("unexpected token %s : %s ", p.currentToken.Type.String(), p.currentToken.Literal))
}
