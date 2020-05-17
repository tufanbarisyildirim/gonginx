package gonginx

import (
	"errors"
)

//Http represents http block
type Http struct {
	Servers    []*Server
	Directives []IDirective
}

//NewHttp create an http block from a directive which has a block
func NewHttp(directive IDirective) (*Http, error) {
	if block := directive.GetBlock(); block != nil {
		http := &Http{
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

		return http, nil
	}
	return nil, errors.New("http directive must have a block")
}

//GetName get directive name to construct the statment string
func (h *Http) GetName() string { //the directive name.
	return "http"
}

//GetParameters get directive parameters if any
func (h *Http) GetParameters() []string {
	return []string{}
}

//GetDirectives get all directives in http
func (h *Http) GetDirectives() []IDirective {
	directives := make([]IDirective, 0)
	for _, directive := range h.Directives {
		directives = append(directives, directive)
	}
	for _, directive := range h.Servers {
		directives = append(directives, directive)
	}
	return directives
}

//FindDirectives find directives
func (h *Http) FindDirectives(directiveName string) []IDirective {
	directives := make([]IDirective, 0)
	for _, directive := range h.GetDirectives() {
		if directive.GetName() == directiveName {
			directives = append(directives, directive)
		}
		if directive.GetBlock() != nil {
			directives = append(directives, directive.GetBlock().FindDirectives(directiveName)...)
		}
	}

	return directives
}

//GetBlock get block if any
func (h *Http) GetBlock() IBlock {
	return h
}
