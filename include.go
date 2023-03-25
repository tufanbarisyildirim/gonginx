package gonginx

// Include include structure
type Include struct {
	*Directive
	IncludePath string
	Configs     []*Config
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
