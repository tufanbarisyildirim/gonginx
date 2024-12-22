package config

import (
	"errors"
)

// HTTP represents http block
type HTTP struct {
	Servers    []*Server
	Directives []IDirective
	Comment    []string
	DefaultInlineComment
	Parent IDirective
	Line   int
}

// SetLine Set line number
func (h *HTTP) SetLine(line int) {
	h.Line = line
}

// GetLine Get the line number
func (h *HTTP) GetLine() int {
	return h.Line
}

// SetParent change the parent block
func (h *HTTP) SetParent(parent IDirective) {
	h.Parent = parent
}

// GetParent the parent block
func (h *HTTP) GetParent() IDirective {
	return h.Parent
}

// GetComment comment of the HTTP directive
func (h *HTTP) GetComment() []string {
	return h.Comment
}

// SetComment set the comment of the HTTP directive
func (h *HTTP) SetComment(comment []string) {
	h.Comment = comment
}

// NewHTTP create an http block from a directive which has a block
func NewHTTP(directive IDirective) (*HTTP, error) {
	if block := directive.GetBlock(); block != nil {
		http := &HTTP{
			Servers:    []*Server{},
			Directives: []IDirective{},
		}
		for _, directive := range block.GetDirectives() {
			if server, ok := directive.(*Server); ok {
				server.Parent = http
				http.Servers = append(http.Servers, server)
				continue
			}
			http.Directives = append(http.Directives, directive)
		}
		http.Comment = directive.GetComment()
		http.InlineComment = directive.GetInlineComment()

		return http, nil
	}
	return nil, errors.New("http directive must have a block")
}

// GetName get directive name to construct the statment string
func (h *HTTP) GetName() string { //the directive name.
	return "http"
}

// GetParameters get directive parameters if any
func (h *HTTP) GetParameters() []Parameter {
	return []Parameter{}
}

// GetDirectives get all directives in http
func (h *HTTP) GetDirectives() []IDirective {
	directives := make([]IDirective, 0)
	directives = append(directives, h.Directives...)
	for _, directive := range h.Servers {
		directives = append(directives, directive)
	}
	return directives
}

// FindDirectives find directives
func (h *HTTP) FindDirectives(directiveName string) []IDirective {
	directives := make([]IDirective, 0)
	for _, directive := range h.GetDirectives() {
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

// GetBlock get block if any
func (h *HTTP) GetBlock() IBlock {
	return h
}

// GetCodeBlock returns the literal code block
func (h *HTTP) GetCodeBlock() string {
	return ""
}
