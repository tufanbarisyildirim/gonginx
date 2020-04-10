package config

import (
	"fmt"
	"strings"
)

//Directive represents any nginx directive
type Directive struct {
	Name       string
	Parameters []string
}

func (d *Directive) directiveStatement() {}

//ToString string repre of a directive
func (d *Directive) ToString() string {
	return fmt.Sprintf("%s %s;", d.Name, strings.Join(d.Parameters, " "))
}
