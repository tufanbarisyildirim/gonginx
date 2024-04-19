package config

// Block a block statement
type Block struct {
	Directives  []IDirective
	IsLuaBlock  bool
	LiteralCode string
	Parent      IBlock
}

func (b *Block) SetParent(parent IBlock) {
	b.Parent = parent
}

func (b *Block) GetParent() IBlock {
	return b.Parent
}

// GetDirectives get all directives in this block
func (b *Block) GetDirectives() []IDirective {
	return b.Directives
}

// GetCodeBlock returns the literal code block
func (b *Block) GetCodeBlock() string {
	return b.LiteralCode
}

// FindDirectives find directives in block recursively
func (b *Block) FindDirectives(directiveName string) []IDirective {
	directives := make([]IDirective, 0)
	for _, directive := range b.GetDirectives() {
		if directive.GetName() == directiveName {
			directives = append(directives, directive)
		}
		if include, ok := directive.(*Include); ok {
			for _, c := range include.Configs {
				directives = append(directives, c.FindDirectives(directiveName)...)
			}
		}
		if directive.GetBlock() != nil {
			directives = append(directives, directive.GetBlock().FindDirectives(directiveName)...)
		}
	}

	return directives
}
