package config

import (
	"fmt"
)

// LuaBlock represents *_by_lua_block
type LuaBlock struct {
	Directives []IDirective
	Name       string
	Comment    []string
	DefaultInlineComment
	LuaCode    string
	Parent     IDirective
	Line       int
	Parameters []Parameter
}

// NewLuaBlock creates a lua block
func NewLuaBlock(directive IDirective) (*LuaBlock, error) {
	if block := directive.GetBlock(); block != nil {
		lb := &LuaBlock{
			Directives: []IDirective{},
			Name:       directive.GetName(),
			LuaCode:    block.GetCodeBlock(),
			Parameters: directive.GetParameters(),
		}

		lb.Directives = append(lb.Directives, block.GetDirectives()...)
		lb.Comment = directive.GetComment()
		lb.InlineComment = directive.GetInlineComment()

		return lb, nil
	}
	return nil, fmt.Errorf("%s must have a block", directive.GetName())
}

// SetLine Set line number
func (lb *LuaBlock) SetLine(line int) {
	lb.Line = line
}

// GetLine Get the line number
func (lb *LuaBlock) GetLine() int {
	return lb.Line
}

// SetParent change the parent block
func (lb *LuaBlock) SetParent(parent IDirective) {
	lb.Parent = parent
}

// GetParent the parent block
func (lb *LuaBlock) GetParent() IDirective {
	return lb.Parent
}

// GetName get directive name to construct the statment string
func (lb *LuaBlock) GetName() string { //the directive name.
	return lb.Name
}

// GetParameters get directive parameters if any
func (lb *LuaBlock) GetParameters() []Parameter {
	return lb.Parameters
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

// GetCodeBlock returns the literal code block
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
