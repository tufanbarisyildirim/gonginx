package config

//Server represents server block
type Server struct {
	*Directive
}

//ToString return config as string
func (s *Server) ToString() string {
	return string(s.Directive.ToByteArray())
}

//ToByteArray return config as byte array
func (s *Server) ToByteArray() []byte {
	return s.Directive.ToByteArray()
}
