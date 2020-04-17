package config

import "github.com/tufanbarisyildirim/gonginx/dumper"

//Server represents server block
type Server struct {
	*Directive
}

//ToString return config as string
func (s *Server) ToString(style *dumper.Style) string {
	return s.Directive.ToString(style)
}
