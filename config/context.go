package config

//Context represents nginx context
type Context interface {
	Node
	ParentContext() Context
	Directives() []Directive
}
