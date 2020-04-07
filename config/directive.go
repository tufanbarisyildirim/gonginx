package config

//Directive represents nginx directive
type Directive struct {
	Node
	Name       string
	Parameters []string
}
