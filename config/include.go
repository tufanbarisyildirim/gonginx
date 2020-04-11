package config

import (
	"fmt"
)

//Include include structure
type Include struct {
	IncludePath string
	*Config
}

func (i *Include) includeStatement() {}

//ToString returns include statement as string
func (i *Include) ToString() string {
	return fmt.Sprintf("include %s;", i.IncludePath)
}

//SaveToFile saves include to its own file
func (i *Include) SaveToFile() error {
	return i.Config.SaveToFile()
}
