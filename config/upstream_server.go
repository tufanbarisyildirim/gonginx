package config

import "strings"

//UpstreamServer represents `server` directive in `upstream{}` block
type UpstreamServer struct {
	*Directive
	Address    string
	Flags      []string
	Parameters map[string]string
}

//ToString convert it to a string
func (uss *UpstreamServer) ToString() string {
	return uss.Directive.ToString()
}

//NewUpstreamServer creates an upstream server from a directive
func NewUpstreamServer(directive *Directive) *UpstreamServer {
	uss := &UpstreamServer{
		Directive:  directive,
		Flags:      make([]string, 0),
		Parameters: make(map[string]string, 0),
	}

	for i, parameter := range directive.Parameters {
		if i == 0 { // alright, we asuume that firstone should be a server address
			uss.Address = parameter
			continue
		}
		if strings.Contains(parameter, "=") { //a parameter like weight=5
			s := strings.SplitN(parameter, "=", 2)
			uss.Parameters[s[0]] = s[1]
		} else {
			uss.Flags = append(uss.Flags, parameter)
		}
	}

	return uss
}
