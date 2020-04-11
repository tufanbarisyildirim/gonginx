package config

//Server represents server block
type Server struct {
	Block
}

//ToString return config as string
func (s *Server) ToString() string {
	return string(s.Block.ToByteArray())
}

//ToByteArray return config as byte array
func (s *Server) ToByteArray() []byte {
	return s.Block.ToByteArray()
}
