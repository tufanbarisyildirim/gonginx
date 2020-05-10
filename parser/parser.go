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
	statementParsers  map[string]func() config.IDirective
	blockWrappers     map[string]func(*config.Directive) config.IDirective
	directiveWrappers map[string]func(*config.Directive) config.IDirective
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
	l := newLexer(bufio.NewReader(f))
	l.file = filePath
	p := NewParserFromLexer(l)
	return p, nil
}

//NewParserFromLexer initilizes a new Parser
func NewParserFromLexer(lexer *lexer) *Parser {
	parser := &Parser{
		lexer: lexer,
	}
	parser.nextToken()
	parser.nextToken()

	parser.blockWrappers = map[string]func(*config.Directive) config.IDirective{
		"http": func(directive *config.Directive) config.IDirective {
			return parser.wrapHttp(directive)
		},
		"server": func(directive *config.Directive) config.IDirective {
			return parser.wrapServer(directive)
		},
		"location": func(directive *config.Directive) config.IDirective {
			return parser.wrapLocation(directive)
		},
	}

	parser.directiveWrappers = map[string]func(*config.Directive) config.IDirective{
		"server": func(directive *config.Directive) config.IDirective {
			return parser.parseUpstreamServer(directive)
		},
		"include": func(directive *config.Directive) config.IDirective {
			return parser.parseInclude(directive)
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
		FilePath: p.lexer.file, //TODO: set filepath here,
		Block:    p.parseBlock(),
	}
}

//ParseBlock parse a block statement
func (p *Parser) parseBlock() *config.Block {

	context := &config.Block{
		Directives: make([]config.IDirective, 0),
	}

parsingloop:
	for {
		switch {
		case p.curTokenIs(token.EOF) || p.curTokenIs(token.BlockEnd):
			break parsingloop
		case p.curTokenIs(token.Keyword):
			context.Directives = append(context.Directives, p.parseStatement())
			break
		}
		p.nextToken()
	}

	return context
}

func (p *Parser) parseStatement() config.IDirective {
	d := &config.Directive{
		Name: p.currentToken.Literal,
	}

	//if we have a special parser for the directive, we use it.
	if sp, ok := p.statementParsers[d.Name]; ok {
		return sp()
	}

	//parse parameters until the end.
	for p.nextToken(); p.currentToken.IsParameterEligible(); p.nextToken() {
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

//TODO: move this into config.Include
func (p *Parser) parseInclude(directive *config.Directive) *config.Include {
	include := &config.Include{
		Directive:   directive,
		IncludePath: directive.Parameters[0],
	}

	if len(directive.Parameters) > 1 {
		panic("include directive can not have multiple parameters")
	}

	if directive.Block != nil {
		panic("include can not have a block, or missing semicolon at the end of include statement")
	}

	return include
}

//TODO: move this into config.Location
func (p *Parser) wrapLocation(directive *config.Directive) *config.Location {
	location := &config.Location{
		Modifier:  "",
		Match:     "",
		Directive: directive,
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
	s, _ := config.NewServer(directive)
	return s
}

func (p *Parser) wrapHttp(directive *config.Directive) *config.Http {
	h, _ := config.NewHttp(directive)
	return h
}

func (p *Parser) parseUpstreamServer(directive *config.Directive) *config.UpstreamServer {
	return config.NewUpstreamServer(directive)
}
