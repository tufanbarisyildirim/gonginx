package config

import (
	"errors"
)

// Server represents server block
type Server struct {
	Block   IBlock
	Comment []string
	DefaultInlineComment
	Parent IDirective
	Line   int
}

// SetLine Set line number
func (s *Server) SetLine(line int) {
	s.Line = line
}

// GetLine Get the line number
func (s *Server) GetLine() int {
	return s.Line
}

// SetParent change the parent block
func (s *Server) SetParent(parent IDirective) {
	s.Parent = parent
}

// GetParent the parent block
func (s *Server) GetParent() IDirective {
	return s.Parent
}

// SetComment set comment of server directive
func (s *Server) SetComment(comment []string) {
	s.Comment = comment
}

// GetComment get comment of server directive
func (s *Server) GetComment() []string {
	return s.Comment
}

// NewServer create a server block from a directive with block
func NewServer(directive IDirective) (*Server, error) {
	if block := directive.GetBlock(); block != nil {
		return &Server{
			Block:                block,
			Comment:              directive.GetComment(),
			DefaultInlineComment: DefaultInlineComment{InlineComment: directive.GetInlineComment()},
		}, nil
	}
	return nil, errors.New("server directive must have a block")
}

// GetName get directive name to construct the statment string
func (s *Server) GetName() string { //the directive name.
	return "server"
}

// GetParameters get directive parameters if any
func (s *Server) GetParameters() []Parameter {
	return []Parameter{}
}

// GetBlock get block if any
func (s *Server) GetBlock() IBlock {
	return s.Block
}

// FindDirectives find directives
func (s *Server) FindDirectives(directiveName string) []IDirective {
	directives := make([]IDirective, 0)
	for _, directive := range s.GetDirectives() {
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

// GetDirectives get all directives in Server
func (s *Server) GetDirectives() []IDirective {
	block := s.GetBlock()
	if block == nil {
		return []IDirective{}
	}
	return block.GetDirectives()
}
