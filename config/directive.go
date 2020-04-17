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
	Parameters []string
}

//ToString string repre of a directive
func (d *Directive) ToString(style *dumper.Style) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s%s", strings.Repeat(" ", style.StartIndent), d.Name))
	if len(d.Parameters) > 0 {
		buf.WriteString(fmt.Sprintf(" %s", strings.Join(d.Parameters, " ")))
	}
	if d.Block == nil {
		buf.WriteRune(';')
	} else {
		buf.WriteString(fmt.Sprintf(" {\n"))
		buf.WriteString(d.Block.ToString(style.Iterate()))
		buf.WriteString(fmt.Sprintf("\n%s}", strings.Repeat(" ", style.StartIndent)))
	}
	return buf.String()
}
