package config

import (
	"errors"
)

// HTTP represents an http block.
type HTTP struct {
	Servers    []*Server
	Directives []IDirective
	Comment    []string
	DefaultInlineComment
	Parent IDirective
	Line   int
}

// SetLine sets the line number.
func (h *HTTP) SetLine(line int) {
	h.Line = line
}

// GetLine returns the line number.
func (h *HTTP) GetLine() int {
	return h.Line
}

// SetParent sets the parent directive.
func (h *HTTP) SetParent(parent IDirective) {
	h.Parent = parent
}

// GetParent returns the parent directive.
func (h *HTTP) GetParent() IDirective {
	return h.Parent
}

// GetComment returns the comment of the HTTP directive.
func (h *HTTP) GetComment() []string {
	return h.Comment
}

// SetComment sets the comment of the HTTP directive.
func (h *HTTP) SetComment(comment []string) {
	h.Comment = comment
}

// NewHTTP creates an HTTP block from a directive that has a block.
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

// GetName returns the directive name to construct the statement string.
func (h *HTTP) GetName() string { //the directive name.
	return "http"
}

// GetParameters returns directive parameters, if any.
func (h *HTTP) GetParameters() []Parameter {
	return []Parameter{}
}

// GetDirectives returns all directives in the http block.
func (h *HTTP) GetDirectives() []IDirective {
	directives := make([]IDirective, 0)
	directives = append(directives, h.Directives...)
	for _, directive := range h.Servers {
		directives = append(directives, directive)
	}
	return directives
}

// FindDirectives finds directives in the http block.
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

// GetBlock returns the block itself.
func (h *HTTP) GetBlock() IBlock {
	return h
}

// GetCodeBlock returns the literal code block.
func (h *HTTP) GetCodeBlock() string {
	return ""
}
