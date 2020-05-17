package gonginx

//Block a block statement
type Block struct {
	Directives []IDirective
}

//GetDirectives get all directives in this block
func (b *Block) GetDirectives() []IDirective {
	return b.Directives
}

//FindDirectives find directives in block recursively
func (b *Block) FindDirectives(directiveName string) []IDirective {
	directives := make([]IDirective, 0)
	for _, directive := range b.GetDirectives() {
		if directive.GetName() == directiveName {
			directives = append(directives, directive)
		}
		if directive.GetBlock() != nil {
			directives = append(directives, directive.GetBlock().FindDirectives(directiveName)...)
		}
	}

	return directives
}
