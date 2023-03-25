package gonginx

import (
	"fmt"
	"sort"
	"strings"
)

// UpstreamServer represents `server` directive in `upstream{}` block
type UpstreamServer struct {
	Address    string
	Flags      []string
	Parameters map[string]string
	Comment    []string
}

// SetComment set comment of the directive
func (uss *UpstreamServer) SetComment(comment []string) {
	uss.Comment = comment
}

// GetComment get comment of the directive
func (uss *UpstreamServer) GetComment() []string {
	return uss.Comment
}

// GetName return directive name, Statement  interface
func (uss *UpstreamServer) GetName() string {
	return "server"
}

// GetBlock block of an upstream, basically nil
func (uss *UpstreamServer) GetBlock() IBlock {
	return nil
}

// GetParameters block of an upstream, basically nil
func (uss *UpstreamServer) GetParameters() []string {
	return uss.GetDirective().Parameters
}

// GetDirective get directive of the upstreamserver
func (uss *UpstreamServer) GetDirective() *Directive {
	//First, generate a new directive from upstream server
	directive := &Directive{
		Name:       "server",
		Parameters: make([]string, 0),
		Block:      nil,
	}

	//address it the first parameter of an upstream directive
	directive.Parameters = append(directive.Parameters, uss.Address)

	//Iterations are random in golang maps https://blog.golang.org/maps#TOC_7.
	//we sort keys in different slice and print them sorted.
	//we always expect key=values parameters to be sorted by key
	paramNames := make([]string, 0)
	for k := range uss.Parameters {
		paramNames = append(paramNames, k)
	}
	sort.Strings(paramNames)

	//append named parameters first
	for _, k := range paramNames {
		directive.Parameters = append(directive.Parameters, fmt.Sprintf("%s=%s", k, uss.Parameters[k]))
	}

	//append flags to the end of the directive.
	directive.Parameters = append(directive.Parameters, uss.Flags...)

	directive.Comment = uss.GetComment()

	return directive
}

// NewUpstreamServer creates an upstream server from a directive
func NewUpstreamServer(directive IDirective) *UpstreamServer {
	uss := &UpstreamServer{
		Flags:      make([]string, 0),
		Parameters: make(map[string]string, 0),
		Comment:    make([]string, 0),
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

	uss.Comment = directive.GetComment()

	return uss
}
