package ast

type Node interface {
	TokenLiteral() string
	String() string
}

type Directive struct {
	Node
	Name       string
	Parameters []string
}

type Context interface {
	Node
	ParentContext() Context
	Directives() []Directive
}

type Config struct {
	Directives []Directive
}

type ServerDirective struct {
	//ports, binding ip addresses, ip versions
	//ssl conf
	//locations
	//servername
}

type LocationDirective struct {
	MatchModifier string
	LocationMatch string
}
