package config

import (
	"errors"
)

// Include include structure
type Include struct {
	*Directive
	IncludePath string
	Configs     []*Config
	Parent      IBlock
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

func (c *Include) GetParent() IBlock {
	return c.Parent
}

func (c *Include) SetParent(parent IBlock) {
	c.Parent = parent
}

// GetDirectives return all directives inside the included file
func (c *Include) GetDirectives() []IDirective {
	directives := make([]IDirective, 0)
	for _, config := range c.Configs {
		directives = append(directives, config.GetDirectives()...)
	}

	return directives
}

// FindDirectives find a specific directive in included file
func (c *Include) FindDirectives(directiveName string) []IDirective {
	directives := make([]IDirective, 0)
	for _, config := range c.Configs {
		directives = append(directives, config.FindDirectives(directiveName)...)
	}

	return directives
}

// GetName get directive name
func (c *Include) GetName() string {
	return c.Directive.Name
}

// SetComment set comment of include directive
func (c *Include) SetComment(comment []string) {
	c.Comment = comment
}

func NewInclude(dir IDirective) (*Include, error) {
	directive, ok := dir.(*Directive)
	if !ok {
		return nil, errors.New("type error")
	}
	include := &Include{
		Directive:   directive,
		IncludePath: directive.Parameters[0],
	}

	if len(directive.Parameters) > 1 {
		panic("include directive can not have multiple parameters")
	}

	if directive.Block != nil {
		panic("include can not have a block, or missing semicolon at the end of include statement")
	}
	return include, nil
}
