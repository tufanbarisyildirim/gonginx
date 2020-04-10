package parser

import (
	"github.com/tufanbarisyildirim/gonginx/config"
	"github.com/tufanbarisyildirim/gonginx/token"
)

//Parser is an nginx config parser
type Parser struct {
	lexer          *lexer
	currentToken   token.Token
	followingToken token.Token
}

//NewParser initilizes a new Parser
func NewParser(lexer *lexer) *Parser {
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
