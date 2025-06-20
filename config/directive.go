package config

// Directive represents any nginx directive.
type Directive struct {
	Block      IBlock
	Name       string
	Parameters []Parameter //TODO: Save parameters with their type
	Comment    []string
	DefaultInlineComment
	Parent IDirective
	Line   int
}

// SetLine sets the line number.
func (d *Directive) SetLine(line int) {
	d.Line = line
}

// GetLine returns the line number.
func (d *Directive) GetLine() int {
	return d.Line
}

// SetParent sets the parent directive.
func (d *Directive) SetParent(parent IDirective) {
	d.Parent = parent
}

// GetParent returns the parent directive.
func (d *Directive) GetParent() IDirective {
	return d.Parent
}

// SetComment sets the directive comment.
func (d *Directive) SetComment(comment []string) {
	d.Comment = comment
}

// GetName returns the directive name.
func (d *Directive) GetName() string {
	return d.Name
}

// GetParameters returns all parameters of the directive.
func (d *Directive) GetParameters() []Parameter {
	return d.Parameters
}

// GetBlock returns the directive block if it exists.
func (d *Directive) GetBlock() IBlock {
	return d.Block
}

// GetComment returns the directive comment.
func (d *Directive) GetComment() []string {
	return d.Comment
}
