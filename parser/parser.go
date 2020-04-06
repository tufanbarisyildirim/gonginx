package parser

import "github.com/tufanbarisyildirim/gonginx/token"

type Parser struct {
	lexer          *Lexer
	currentToken   token.Token
	followingToken token.Token
}

func NewParser(lexer *Lexer) *Parser {
	parser := &Parser{
		lexer: lexer,
	}
	parser.nextToken()
	parser.nextToken()

	return parser
}

func (p *Parser) nextToken() {
	p.currentToken = p.followingToken
	p.followingToken = p.lexer.Scan()
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.currentToken.Type == t
}

func (p *Parser) followingTokenIs(t token.Type) bool {
	return p.followingToken.Type == t
}