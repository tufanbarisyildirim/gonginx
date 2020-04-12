package config

import (
	"fmt"
)

//Include include structure
type Include struct {
	IncludePath string
	*Config
}

//ToString returns include statement as string
func (i *Include) ToString() string {
	return fmt.Sprintf("include %s;", i.IncludePath)
}

//SaveToFile saves include to its own file
func (i *Include) SaveToFile() error {
	if i.Config == nil {
		return fmt.Errorf("included empty file %s", i.IncludePath)
	}
	return i.Config.SaveToFile()
}
