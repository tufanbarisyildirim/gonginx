package config

import (
	"fmt"
	"sort"
	"strings"
)

// UpstreamServer represents a `server` directive in an `upstream{}` block.
type UpstreamServer struct {
	Address    string
	Flags      []string
	Parameters map[string]string
	Comment    []string
	DefaultInlineComment
	Parent IDirective
	Line   int
}

// SetLine sets the line number.
func (uss *UpstreamServer) SetLine(line int) {
	uss.Line = line
}

// GetLine returns the line number.
func (uss *UpstreamServer) GetLine() int {
	return uss.Line
}

// SetParent sets the parent directive.
func (uss *UpstreamServer) SetParent(parent IDirective) {
	uss.Parent = parent
}

// GetParent returns the parent directive.
func (uss *UpstreamServer) GetParent() IDirective {
	return uss.Parent
}

// SetComment sets the directive comment.
func (uss *UpstreamServer) SetComment(comment []string) {
	uss.Comment = comment
}

// GetComment returns the directive comment.
func (uss *UpstreamServer) GetComment() []string {
	return uss.Comment
}

// GetName implements the Statement interface.
func (uss *UpstreamServer) GetName() string {
	return "server"
}

// GetBlock returns nil because upstream servers do not have blocks.
func (uss *UpstreamServer) GetBlock() IBlock {
	return nil
}

// GetParameters returns parameters for the upstream server.
func (uss *UpstreamServer) GetParameters() []Parameter {
	return uss.GetDirective().Parameters
}

// GetDirective returns the directive representation of the upstream server.
func (uss *UpstreamServer) GetDirective() *Directive {
	//First, generate a new directive from upstream server
	directive := &Directive{
		Name:       "server",
		Parameters: make([]Parameter, 0),
		Block:      nil,
	}

	//address it the first parameter of an upstream directive
	directive.Parameters = append(directive.Parameters, Parameter{Value: uss.Address})

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
		directive.Parameters = append(directive.Parameters, Parameter{Value: fmt.Sprintf("%s=%s", k, uss.Parameters[k])})
	}

	//append flags to the end of the directive.
	for _, flag := range uss.Flags {
		directive.Parameters = append(directive.Parameters, Parameter{Value: flag})
	}

	directive.Comment = uss.GetComment()

	return directive
}

// NewUpstreamServer creates an UpstreamServer from a directive.
func NewUpstreamServer(directive IDirective) (*UpstreamServer, error) {
	uss := &UpstreamServer{
		Flags:      make([]string, 0),
		Parameters: make(map[string]string, 0),
		Comment:    make([]string, 0),
	}

	for i, parameter := range directive.GetParameters() {
		if i == 0 { // alright, we asuume that firstone should be a server address
			uss.Address = parameter.GetValue()
			continue
		}
		if strings.Contains(parameter.GetValue(), "=") { //a parameter like weight=5
			s := strings.SplitN(parameter.GetValue(), "=", 2)
			uss.Parameters[s[0]] = s[1]
		} else {
			uss.Flags = append(uss.Flags, parameter.GetValue())
		}
	}

	uss.Comment = directive.GetComment()
	uss.InlineComment = directive.GetInlineComment()

	return uss, nil
}
