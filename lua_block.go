package gonginx

import (
	"fmt"
)

// LuaBlock represents *_by_lua_block
type LuaBlock struct {
	Directives []IDirective
	Name       string
	Comment    []string
	LuaCode    string
}

// NewLuaBlock creates a lua block
func NewLuaBlock(directive IDirective) (*LuaBlock, error) {
	if block := directive.GetBlock(); block != nil {
		lb := &LuaBlock{
			Directives: []IDirective{},
			Name:       directive.GetName(),
			LuaCode:    block.GetCodeBlock(),
		}

		lb.Directives = append(lb.Directives, block.GetDirectives()...)
		lb.Comment = directive.GetComment()

		return lb, nil
	}
	return nil, fmt.Errorf("%s must have a block", directive.GetName())
}

// GetName get directive name to construct the statment string
func (lb *LuaBlock) GetName() string { //the directive name.
	return lb.Name
}

// GetParameters get directive parameters if any
func (lb *LuaBlock) GetParameters() []string {
	return []string{}
}

// GetDirectives get all directives in lua block
// this should return 1 code block and that should be it
func (lb *LuaBlock) GetDirectives() []IDirective {
	directives := make([]IDirective, 0)
	directives = append(directives, lb.Directives...)
	return directives
}

// FindDirectives find directives
func (lb *LuaBlock) FindDirectives(directiveName string) []IDirective {
	directives := make([]IDirective, 0)
	for _, directive := range lb.GetDirectives() {
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

func (lb *LuaBlock) GetCodeBlock() string {
	return lb.LuaCode
}

// GetBlock get block if any
func (lb *LuaBlock) GetBlock() IBlock {
	return lb
}

// GetComment get directive comment
func (lb *LuaBlock) GetComment() []string {
	return lb.Comment
}

// SetComment sets comment tied to this directive
func (lb *LuaBlock) SetComment(comment []string) {
	lb.Comment = comment
}
