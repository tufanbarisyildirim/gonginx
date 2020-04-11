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

func (d *Directive) directiveStatement() {}

//ToString string repre of a directive
func (d *Directive) ToString() string {
	if d.Block == nil {
		return fmt.Sprintf("%s %s;", d.Name, strings.Join(d.Parameters, " "))
	} else {
		return fmt.Sprintf("%s %s {\n%s\n}", d.Name, strings.Join(d.Parameters, " "), d.Block.ToString())
	}
}
