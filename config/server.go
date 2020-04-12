package config

//Server represents server block
type Server struct {
	*Directive
}

//ToString return config as string
func (s *Server) ToString() string {
	return s.Directive.ToString()
}
