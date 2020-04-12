package parser

import (
	"bufio"
	"fmt"
	"os"

	"github.com/tufanbarisyildirim/gonginx/config"
	"github.com/tufanbarisyildirim/gonginx/parser/token"
)

//Parser is an nginx config parser
type Parser struct {
	lexer             *lexer
	currentToken      token.Token
	followingToken    token.Token
	statementParsers  map[string]func() config.Statement
	blockWrappers     map[string]func(*config.Directive) config.Statement
	directiveWrappers map[string]func(*config.Directive) config.Statement
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

	parser.statementParsers = map[string]func() config.Statement{
		"include": func() config.Statement {
			return parser.parseInclude()
		},
	}

	parser.blockWrappers = map[string]func(*config.Directive) config.Statement{
		"server": func(directive *config.Directive) config.Statement {
			return parser.wrapServer(directive)
		},
		"location": func(directive *config.Directive) config.Statement {
			return parser.wrapLocation(directive)
		},
	}

	parser.directiveWrappers = map[string]func(*config.Directive) config.Statement{
		"server": func(directive *config.Directive) config.Statement {
			return parser.parseUpstreamServer(directive)
		},
	}

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
		case p.curTokenIs(token.EOF):
			break parsingloop
		case p.curTokenIs(token.Keyword):
			context.Statements = append(context.Statements, p.parseStatement())
			break
		}
		p.nextToken()
	}

	return context
}

func (p *Parser) parseStatement() config.Statement {
	d := &config.Directive{
		Name: p.currentToken.Literal,
	}

	//if we have a special parser for the directive, we use it.
	if sp, ok := p.statementParsers[d.Name]; ok {
		return sp()
	}

	//parse parameters until the end.
	for p.nextToken(); p.curTokenIs(token.Keyword); p.nextToken() {
		d.Parameters = append(d.Parameters, p.currentToken.Literal)
	}

	//if we find a semicolon it is a directive, we will check directive converters
	if p.curTokenIs(token.Semicolon) {
		if dw, ok := p.directiveWrappers[d.Name]; ok {
			return dw(d)
		}
		return d
	}

	//ok, it does not end with a semicolon but a block starts, we will convert that block if we have a converter
	if p.curTokenIs(token.BlockStart) {
		d.Block = p.parseBlock()
		if bw, ok := p.blockWrappers[d.Name]; ok {
			return bw(d)
		}
		return d
	}

	panic(fmt.Errorf("unexpected token %s (%s) on line %d, column %d", p.currentToken.Type.String(), p.currentToken.Literal, p.currentToken.Line, p.currentToken.Column))
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

func (p *Parser) wrapLocation(directive *config.Directive) *config.Location {
	location := &config.Location{
		Modifier: "",
		Match:    "",
		Block:    directive.Block,
	}

	if len(directive.Parameters) == 0 {
		panic("no enough parameter for location")
	}

	if len(directive.Parameters) == 1 {
		location.Match = directive.Parameters[0]
		return location
	} else if len(directive.Parameters) == 2 {
		location.Modifier = directive.Parameters[0]
		location.Match = directive.Parameters[1]
		return location
	}

	panic("too many arguments for location directive")
}

func (p *Parser) wrapServer(directive *config.Directive) *config.Server {
	server := &config.Server{
		Directive: directive,
	}

	return server
}

func (p *Parser) parseUpstreamServer(directive *config.Directive) *config.UpstreamServer {
	upstreamServer := &config.UpstreamServer{
		Directive: directive,
	}

	//TODO: param 1 should be the server, with port.
	//others should be parsed as key=value
	//some of them line down, etc are sub directives or arguments.

	return upstreamServer
}
