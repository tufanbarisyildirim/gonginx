package config

import (
	"fmt"
	"strings"
)

//Directive represents any nginx directive
type Directive struct {
	*Block
	Name       string
	Parameters []string
}

//ToString string repre of a directive
func (d *Directive) ToString() string {
	if d.Block == nil {
		return fmt.Sprintf("%s %s;", d.Name, strings.Join(d.Parameters, " "))
	}
	return fmt.Sprintf("%s {\n%s\n}", strings.Join(append([]string{d.Name}, d.Parameters...), " "), d.Block.ToString())
}
