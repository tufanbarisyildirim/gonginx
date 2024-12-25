package config

// IBlock represents any directive block
type IBlock interface {
	GetDirectives() []IDirective
	FindDirectives(directiveName string) []IDirective
	GetCodeBlock() string
	SetParent(IDirective)
	GetParent() IDirective
}

// IDirective represents any directive
type IDirective interface {
	GetName() string //the directive name.
	GetParameters() []Parameter
	GetBlock() IBlock
	GetComment() []string
	SetComment(comment []string)
	SetParent(IDirective)
	GetParent() IDirective
	GetLine() int
	SetLine(int)
	InlineCommenter
}

// InlineCommenter represents the inline comment holder
type InlineCommenter interface {
	GetInlineComment() string
	SetInlineComment(comment string)
}

// DefaultInlineComment represents the default inline comment holder
type DefaultInlineComment struct {
	InlineComment string
}

// GetInlineComment returns the inline comment
func (d *DefaultInlineComment) GetInlineComment() string {
	return d.InlineComment
}

// SetInlineComment sets the inline comment
func (d *DefaultInlineComment) SetInlineComment(comment string) {
	d.InlineComment = comment
}

// FileDirective a statement that saves its own file
type FileDirective interface {
	isFileDirective()
}

// IncludeDirective represents include statement in nginx
type IncludeDirective interface {
	FileDirective
}

// Parameter represents a parameter in a directive
type Parameter struct {
	Value             string
	RelativeLineIndex int // relative line index to the directive
}

// String returns the value of the parameter
func (p *Parameter) String() string {
	return p.Value
}

// SetValue sets the value of the parameter
func (p *Parameter) SetValue(v string) {
	p.Value = v
}

// GetValue returns the value of the parameter
func (p *Parameter) GetValue() string {
	return p.Value
}

// SetRelativeLineIndex sets the relative line index of the parameter
func (p *Parameter) SetRelativeLineIndex(i int) {
	p.RelativeLineIndex = i
}

// GetRelativeLineIndex returns the relative line index of the parameter
func (p *Parameter) GetRelativeLineIndex() int {
	return p.RelativeLineIndex
}
