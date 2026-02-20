package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Server represents a server block.
type Server struct {
	Block   IBlock
	Comment []string
	DefaultInlineComment
	Parent IDirective
	Line   int
}

// SetLine sets the line number.
func (s *Server) SetLine(line int) {
	s.Line = line
}

// GetLine returns the line number.
func (s *Server) GetLine() int {
	return s.Line
}

// SetParent sets the parent directive.
func (s *Server) SetParent(parent IDirective) {
	s.Parent = parent
}

// GetParent returns the parent directive.
func (s *Server) GetParent() IDirective {
	return s.Parent
}

// SetComment sets the comment of the server directive.
func (s *Server) SetComment(comment []string) {
	s.Comment = comment
}

// GetComment returns the comment of the server directive.
func (s *Server) GetComment() []string {
	return s.Comment
}

// NewServer creates a server block from a directive with a block.
func NewServer(directive IDirective) (*Server, error) {
	if block := directive.GetBlock(); block != nil {
		return &Server{
			Block:                block,
			Comment:              directive.GetComment(),
			DefaultInlineComment: DefaultInlineComment{InlineComment: directive.GetInlineComment()},
		}, nil
	}
	return nil, errors.New("server directive must have a block")
}

// GetName returns the directive name to construct the statement string.
func (s *Server) GetName() string { //the directive name.
	return "server"
}

// GetParameters returns directive parameters if any.
func (s *Server) GetParameters() []Parameter {
	return []Parameter{}
}

// GetBlock returns the block if any.
func (s *Server) GetBlock() IBlock {
	return s.Block
}

// FindDirectives finds directives within the server block.
func (s *Server) FindDirectives(directiveName string) []IDirective {
	directives := make([]IDirective, 0)
	for _, directive := range s.GetDirectives() {
		if directive.GetName() == directiveName {
			directives = append(directives, directive)
		}
		if include, ok := directive.(*Include); ok {
			for _, c := range include.Configs {
				directives = append(directives, c.FindDirectives(directiveName)...)
			}
		}
		if directive.GetBlock() != nil {
			directives = append(directives, directive.GetBlock().FindDirectives(directiveName)...)
		}
	}

	return directives
}

// GetDirectives returns all directives in the server.
func (s *Server) GetDirectives() []IDirective {
	block := s.GetBlock()
	if block == nil {
		return []IDirective{}
	}
	return block.GetDirectives()
}

// GetLocations returns direct child location blocks in the server.
func (s *Server) GetLocations() []*Location {
	locations := make([]*Location, 0)
	for _, directive := range s.GetDirectives() {
		location, ok := directive.(*Location)
		if !ok {
			continue
		}
		locations = append(locations, location)
	}
	return locations
}

// GetListenPorts returns listen ports from direct child listen directives.
// Non-port listen endpoints (for example unix sockets) are ignored.
func (s *Server) GetListenPorts() []int {
	ports := make([]int, 0)
	for _, listen := range s.getListenDirectives() {
		if len(listen.Parameters) == 0 {
			continue
		}
		port, ok := parseListenPortValue(listen.Parameters[0].GetValue())
		if !ok {
			continue
		}
		ports = append(ports, port)
	}
	return ports
}

// SetListenPort updates the port of the listen directive at the given index.
func (s *Server) SetListenPort(index int, port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("invalid listen port: %d", port)
	}

	listens := s.getListenDirectives()
	if index < 0 || index >= len(listens) {
		return fmt.Errorf("listen index out of range: %d", index)
	}

	listen := listens[index]
	if len(listen.Parameters) == 0 {
		listen.Parameters = append(listen.Parameters, Parameter{Value: strconv.Itoa(port)})
		return nil
	}

	updated, err := formatListenEndpointWithPort(listen.Parameters[0].GetValue(), port)
	if err != nil {
		return err
	}
	listen.Parameters[0].SetValue(updated)
	return nil
}

func (s *Server) getListenDirectives() []*Directive {
	listens := make([]*Directive, 0)
	for _, directive := range s.GetDirectives() {
		listen, ok := directive.(*Directive)
		if !ok {
			continue
		}
		if listen.GetName() != "listen" {
			continue
		}
		listens = append(listens, listen)
	}
	return listens
}

func parseListenPortValue(value string) (int, bool) {
	value = strings.TrimSpace(value)
	if value == "" || strings.HasPrefix(value, "unix:") {
		return 0, false
	}
	if isDecimal(value) {
		port, err := strconv.Atoi(value)
		if err != nil {
			return 0, false
		}
		return port, true
	}
	if strings.HasPrefix(value, "[") {
		sep := strings.LastIndex(value, "]:")
		if sep > 0 && isDecimal(value[sep+2:]) {
			port, err := strconv.Atoi(value[sep+2:])
			if err != nil {
				return 0, false
			}
			return port, true
		}
	}
	sep := strings.LastIndex(value, ":")
	if sep > 0 && isDecimal(value[sep+1:]) {
		port, err := strconv.Atoi(value[sep+1:])
		if err != nil {
			return 0, false
		}
		return port, true
	}
	return 0, false
}

func formatListenEndpointWithPort(value string, port int) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" || isDecimal(value) {
		return strconv.Itoa(port), nil
	}
	if strings.HasPrefix(value, "unix:") {
		return "", fmt.Errorf("cannot set port on unix socket listen %q", value)
	}
	if strings.HasPrefix(value, "[") {
		sep := strings.LastIndex(value, "]:")
		if sep > 0 && isDecimal(value[sep+2:]) {
			return value[:sep+2] + strconv.Itoa(port), nil
		}
	}
	sep := strings.LastIndex(value, ":")
	if sep > 0 && isDecimal(value[sep+1:]) {
		return value[:sep+1] + strconv.Itoa(port), nil
	}
	return "", fmt.Errorf("unsupported listen value %q", value)
}

func isDecimal(value string) bool {
	if value == "" {
		return false
	}
	for i := 0; i < len(value); i++ {
		if value[i] < '0' || value[i] > '9' {
			return false
		}
	}
	return true
}

// AddLocation appends a location block to the server.
func (s *Server) AddLocation(location *Location) {
	if location == nil {
		return
	}

	if s.Block == nil {
		s.Block = &Block{Directives: []IDirective{}}
	}

	block, ok := s.Block.(*Block)
	if !ok {
		block = &Block{
			Directives: append([]IDirective(nil), s.Block.GetDirectives()...),
		}
		s.Block = block
	}

	location.SetParent(s)
	if locationBlock := location.GetBlock(); locationBlock != nil {
		locationBlock.SetParent(location)
	}

	block.Directives = append(block.Directives, location)
}
