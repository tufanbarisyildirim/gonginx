package gonginx

import (
	"errors"
)

// HTTP represents http block
type HTTP struct {
	Servers    []*Server
	Directives []IDirective
	Comment    []string
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
				http.Servers = append(http.Servers, server)
				continue
			}
			http.Directives = append(http.Directives, directive)
		}
		http.Comment = directive.GetComment()

		return http, nil
	}
	return nil, errors.New("http directive must have a block")
}

// GetName get directive name to construct the statment string
func (h *HTTP) GetName() string { //the directive name.
	return "http"
}

// GetParameters get directive parameters if any
func (h *HTTP) GetParameters() []string {
	return []string{}
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

func (h *HTTP) GetCodeBlock() string {
	return ""
}
