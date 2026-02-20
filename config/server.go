package config

import (
	"errors"
)

// Server represents a server block.
type Server struct {
	Block   IBlock
	Comment []string
	DefaultInlineComment
	Parent IDirective
	Line   int
}

// SetLine sets the line number.
func (s *Server) SetLine(line int) {
	s.Line = line
}

// GetLine returns the line number.
func (s *Server) GetLine() int {
	return s.Line
}

// SetParent sets the parent directive.
func (s *Server) SetParent(parent IDirective) {
	s.Parent = parent
}

// GetParent returns the parent directive.
func (s *Server) GetParent() IDirective {
	return s.Parent
}

// SetComment sets the comment of the server directive.
func (s *Server) SetComment(comment []string) {
	s.Comment = comment
}

// GetComment returns the comment of the server directive.
func (s *Server) GetComment() []string {
	return s.Comment
}

// NewServer creates a server block from a directive with a block.
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

// GetName returns the directive name to construct the statement string.
func (s *Server) GetName() string { //the directive name.
	return "server"
}

// GetParameters returns directive parameters if any.
func (s *Server) GetParameters() []Parameter {
	return []Parameter{}
}

// GetBlock returns the block if any.
func (s *Server) GetBlock() IBlock {
	return s.Block
}

// FindDirectives finds directives within the server block.
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

// GetDirectives returns all directives in the server.
func (s *Server) GetDirectives() []IDirective {
	block := s.GetBlock()
	if block == nil {
		return []IDirective{}
	}
	return block.GetDirectives()
}

// AddLocation appends a location block to the server.
func (s *Server) AddLocation(location *Location) {
	if location == nil {
		return
	}

	if s.Block == nil {
		s.Block = &Block{Directives: []IDirective{}}
	}

	block, ok := s.Block.(*Block)
	if !ok {
		block = &Block{
			Directives: append([]IDirective(nil), s.Block.GetDirectives()...),
		}
		s.Block = block
	}

	location.SetParent(s)
	if locationBlock := location.GetBlock(); locationBlock != nil {
		locationBlock.SetParent(location)
	}

	block.Directives = append(block.Directives, location)
}
