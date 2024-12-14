package config

// IBlock represents any directive block
type IBlock interface {
	GetDirectives() []IDirective
	FindDirectives(directiveName string) []IDirective
	GetCodeBlock() string
	SetParent(IBlock)
	GetParent() IBlock
}

// IDirective represents any directive
type IDirective interface {
	GetName() string //the directive name.
	GetParameters() []Parameter
	GetBlock() IBlock
	GetComment() []string
	SetComment(comment []string)
	SetParent(IBlock)
	GetParent() IBlock
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

func (p *Parameter) String() string {
	return p.Value
}

func (p *Parameter) SetValue(v string) {
	p.Value = v
}

func (p *Parameter) GetValue() string {
	return p.Value
}

func (p *Parameter) SetRelativeLineIndex(i int) {
	p.RelativeLineIndex = i
}

func (p *Parameter) GetRelativeLineIndex() int {
	return p.RelativeLineIndex
}
