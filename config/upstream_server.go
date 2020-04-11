package config

//UpstreamServer represents `server` directive in `upstream{}` block
type UpstreamServer struct {
	*Directive
	ServerIP string
	Port     int
	Weight   int
}

//ToString convert it to a string
func (us *UpstreamServer) ToString() string {
	return us.Directive.ToString()
}
