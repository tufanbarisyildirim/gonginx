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

//Parse parses a config and returns an AST
func (p *Parser) Parse() *config.Config {
	c := &config.Config{
		FilePath: p.lexer.file,
	}

parsingloop:
	for {
		switch {
		case p.curTokenIs(token.Eof):
			break parsingloop
		case p.curTokenIs(token.Keyword):
			switch p.currentToken.Literal {
			case "events":
				//parse event context
				p.nextToken()
				//c.Statements = append(c.Statements, p.parseEvents())
				break
			case "http":
				//parse http context
				break
			case "server":
				//parser server context
				break
			case "location":
				//parse location context
				break
			case "types":
				//parse mime types
				break
			case "upstream":
				//parse upstream context
				break
			case "include":
				//parse all files match include statement
				c.Statements = append(c.Statements, p.parseInclude())
				break
			default:
				//parse unknown directive
				break
			}
			break
		}

	}

	return c
}

func (p *Parser) parseInclude() *config.Include {
	include := &config.Include{}
	include.Token = p.currentToken
	include.IncludePath = p.followingToken.String()

	p.nextToken() // read include
	p.nextToken() // read path
	if !p.curTokenIs(token.Semicolon) {
		panic("expected semicolon after include path")
	}

	return include
}
