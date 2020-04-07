package config

//Config represents nginx conf file
type Config struct {
	Directives []Directive
	Context    Context
}
