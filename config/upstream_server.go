package config

//UpstreamServer represents `server` directive in `upstream{}` block
type UpstreamServer struct {
	*Directive
	ServerIP string
	Port     int
	Weight   int
}

//ToString convert it to a string
func (uss *UpstreamServer) ToString() string {
	return uss.Directive.ToString()
}
