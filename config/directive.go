package config

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/tufanbarisyildirim/gonginx/dumper"
)

//Directive represents any nginx directive
type Directive struct {
	*Block
	Name       string
	Parameters []string //TODO: Save parameters with their type
}

//ToString string repre of a directive
func (d *Directive) ToString(style *dumper.Style) string {
	var buf bytes.Buffer
	if style.SpaceBeforeBlocks && d.Block != nil {
		buf.WriteString("\n")
	}
	buf.WriteString(fmt.Sprintf("%s%s", strings.Repeat(" ", style.StartIndent), d.Name))
	if len(d.Parameters) > 0 {
		buf.WriteString(fmt.Sprintf(" %s", strings.Join(d.Parameters, " ")))
	}
	if d.Block == nil {
		buf.WriteRune(';')
	} else {
		buf.WriteString(" {\n")
		buf.WriteString(d.Block.ToString(style.Iterate()))
		buf.WriteString(fmt.Sprintf("\n%s}", strings.Repeat(" ", style.StartIndent)))
	}
	return buf.String()
}

//GetName get directive name
func (d *Directive) GetName() string {
	return d.Name
}
