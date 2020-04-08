package parser

import "github.com/tufanbarisyildirim/gonginx/token"

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

func (p *Parser) parseBlock()