package config

import (
	"errors"
)

// Include represents an include directive.
type Include struct {
	*Directive
	IncludePath string
	Configs     []*Config
	Parent      IDirective
}

//TODO(tufan): move that part into dumper package
//SaveToFile saves include to its own file
//func (i *Include) SaveToFile(style *dumper.Style) error {
//	if len(i.Configs) == 0 {
//		return fmt.Errorf("included empty file %s", i.IncludePath)
//	}
//	for _, c := range i.Configs {
//		err := c.SaveToFile(style)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

// SetLine sets the line number.
func (c *Include) SetLine(line int) {
	c.Line = line
}

// GetLine returns the line number.
func (c *Include) GetLine() int {
	return c.Line
}

// GetParent returns the parent directive.
func (c *Include) GetParent() IDirective {
	return c.Parent
}

// SetParent sets the parent directive.
func (c *Include) SetParent(parent IDirective) {
	c.Parent = parent
}

// GetDirectives returns all directives inside the included file.
func (c *Include) GetDirectives() []IDirective {
	directives := make([]IDirective, 0)
	for _, config := range c.Configs {
		directives = append(directives, config.GetDirectives()...)
	}

	return directives
}

// FindDirectives finds a specific directive in the included file.
func (c *Include) FindDirectives(directiveName string) []IDirective {
	directives := make([]IDirective, 0)
	for _, config := range c.Configs {
		directives = append(directives, config.FindDirectives(directiveName)...)
	}

	return directives
}

// GetName returns the directive name.
func (c *Include) GetName() string {
	return c.Directive.Name
}

// SetComment sets the comment of the include directive.
func (c *Include) SetComment(comment []string) {
	c.Comment = comment
}

// NewInclude initializes an Include from a directive.
func NewInclude(dir IDirective) (*Include, error) {
	directive, ok := dir.(*Directive)
	if !ok {
		return nil, errors.New("type error")
	}
	include := &Include{
		Directive:   directive,
		IncludePath: directive.Parameters[0].GetValue(),
	}

	if len(directive.Parameters) > 1 {
		panic("include directive can not have multiple parameters")
	}

	if directive.Block != nil {
		panic("include can not have a block, or missing semicolon at the end of include statement")
	}
	return include, nil
}
