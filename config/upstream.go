package config

import "github.com/tufanbarisyildirim/gonginx/dumper"

//Upstream represents `upstream{}` block
type Upstream struct {
	*Directive
	Name string
}

//ToString convert it to a string
func (us *Upstream) ToString(style *dumper.Style) string {
	return us.Directive.ToString(style)
}
