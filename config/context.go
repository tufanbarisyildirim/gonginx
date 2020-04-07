package config

type Context interface {
	Node
	ParentContext() Context
	Directives() []Directive
}
