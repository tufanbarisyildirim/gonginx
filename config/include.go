package config

import (
	"fmt"

	"github.com/tufanbarisyildirim/gonginx/dumper"
)

//Include include structure
type Include struct {
	IncludePath string
	Configs     []*Config
}

//ToString returns include statement as string
func (i *Include) ToString(style *dumper.Style) string {
	return fmt.Sprintf("include %s;", i.IncludePath)
}

//SaveToFile saves include to its own file
func (i *Include) SaveToFile(style *dumper.Style) error {
	if len(i.Configs) == 0 {
		return fmt.Errorf("included empty file %s", i.IncludePath)
	}
	for _, c := range i.Configs {
		err := c.SaveToFile(style)
		if err != nil {
			return err
		}
	}
	return nil
}
