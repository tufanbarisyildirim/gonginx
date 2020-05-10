package config

import (
	"fmt"
	"strings"

	"github.com/tufanbarisyildirim/gonginx/dumper"
)

//UpstreamServer represents `server` directive in `upstream{}` block
type UpstreamServer struct {
	Address    string
	Flags      []string
	Parameters map[string]string
}

//GetName return directive name, config.Statement  interface
func (uss *UpstreamServer) GetName() string {
	return "server"
}

//GetBlock block of an upstream, basically nil
func (uss *UpstreamServer) GetBlock() *Block {
	return nil
}

//GetParameters block of an upstream, basically nil
func (uss *UpstreamServer) GetParameters() []string {
	return uss.GetDirective().Parameters
}

//GetDirective get directive of the upstreamserver
func (uss *UpstreamServer) GetDirective() *Directive {
	//First, generate a new directive from upstream server
	directive := &Directive{
		Name:       "server",
		Parameters: make([]string, 1+len(uss.Flags)+len(uss.Parameters)),
		Block:      nil,
	}

	directive.Parameters[0] = uss.Address

	pIndex := 1
	for pName, pVal := range uss.Parameters {
		directive.Parameters[pIndex] = fmt.Sprintf("%s=%s", pName, pVal)
		pIndex++
	}

	for _, flag := range uss.Flags {
		directive.Parameters[pIndex] = flag
		pIndex++
	}

	return directive
}

//ToString convert it to a string
func (uss *UpstreamServer) ToString(style *dumper.Style) string {
	return uss.GetDirective().ToString(style)
}

//NewUpstreamServer creates an upstream server from a directive
func NewUpstreamServer(directive IDirective) *UpstreamServer {
	uss := &UpstreamServer{
		Flags:      make([]string, 0),
		Parameters: make(map[string]string, 0),
	}

	for i, parameter := range directive.GetParameters() {
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
