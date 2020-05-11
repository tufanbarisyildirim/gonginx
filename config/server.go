package config

import (
	"errors"

	"github.com/tufanbarisyildirim/gonginx/dumper"
)

//Server represents server block
type Server struct {
	Block
}

//NewServer create a server block from a directive with block
func NewServer(directive IDirective) (*Server, error) {
	if block := directive.GetBlock(); block != nil {
		return &Server{
			Block: *block,
		}, nil
	}

	return nil, errors.New("server directive must have a block")
}

//GetName get directive name to construct the statment string
func (s *Server) GetName() string { //the directive name.
	return "server"
}

//GetParameters get directive parameters if any
func (s *Server) GetParameters() []string {
	return []string{}
}

//GetBlock get block if any
func (s *Server) GetBlock() *Block {
	return &s.Block
}

//ToString return config as string
func (s *Server) ToString(style *dumper.Style) string {
	directive := Directive{
		Block:      s.GetBlock(),
		Name:       s.GetName(),
		Parameters: s.GetParameters(),
	}
	return directive.ToString(style)
}
