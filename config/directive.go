package config

// Directive represents any nginx directive
type Directive struct {
	Block      IBlock
	Name       string
	Parameters []Parameter //TODO: Save parameters with their type
	Comment    []string
	DefaultInlineComment
	Parent IBlock
	Line   int
}

// SetLine Set line number
func (d *Directive) SetLine(line int) {
	d.Line = line
}

// GetLine Get the line number
func (d *Directive) GetLine() int {
	return d.Line
}

// SetParent  the parent block
func (d *Directive) SetParent(parent IBlock) {
	d.Parent = parent
}

// GetParent change the parent block
func (d *Directive) GetParent() IBlock {
	return d.Parent
}

// SetComment sets comment tied to this directive
func (d *Directive) SetComment(comment []string) {
	d.Comment = comment
}

// GetName get directive name
func (d *Directive) GetName() string {
	return d.Name
}

// GetParameters get all parameters of a directive
func (d *Directive) GetParameters() []Parameter {
	return d.Parameters
}

// GetBlock get block if it has
func (d *Directive) GetBlock() IBlock {
	return d.Block
}

// GetComment get directive comment
func (d *Directive) GetComment() []string {
	return d.Comment
}
