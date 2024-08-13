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
	GetParameters() []string
	GetBlock() IBlock
	GetComment() []string
	SetComment(comment []string)
	SetParent(IBlock)
	GetParent() IBlock
	InlineCommenter
}

type InlineCommenter interface {
	GetInlineComment() string
	SetInlineComment(comment string)
}

type DefaultInline struct {
	InlineComment string
}

func (d *DefaultInline) GetInlineComment() string {
	return d.InlineComment
}

func (d *DefaultInline) SetInlineComment(comment string) {
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
