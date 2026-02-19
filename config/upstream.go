package config

import (
	"errors"
)

// Upstream represents an `upstream{}` block.
type Upstream struct {
	UpstreamName    string
	UpstreamServers []*UpstreamServer
	//Directives Other directives in upstream (ip_hash; etc)
	Directives []IDirective
	Comment    []string
	DefaultInlineComment
	Parent IDirective
	Line   int
}

// SetLine sets the line number.
func (us *Upstream) SetLine(line int) {
	us.Line = line
}

// GetLine returns the line number.
func (us *Upstream) GetLine() int {
	return us.Line
}

// SetParent sets the parent directive.
func (us *Upstream) SetParent(parent IDirective) {
	us.Parent = parent
}

// GetParent returns the parent directive.
func (us *Upstream) GetParent() IDirective {
	return us.Parent
}

// SetComment sets the directive comment.
func (us *Upstream) SetComment(comment []string) {
	us.Comment = comment
}

// GetName implements the Statement interface.
func (us *Upstream) GetName() string {
	return "upstream"
}

// GetParameters returns the upstream parameters.
func (us *Upstream) GetParameters() []Parameter {
	return []Parameter{{Value: us.UpstreamName}} //the only parameter for an upstream is its name
}

// GetBlock returns the upstream itself, which implements IBlock.
func (us *Upstream) GetBlock() IBlock {
	return us
}

// GetComment returns the directive comment.
func (us *Upstream) GetComment() []string {
	return us.Comment
}

// GetDirectives returns sub directives of the upstream.
func (us *Upstream) GetDirectives() []IDirective {
	directives := make([]IDirective, 0)
	directives = append(directives, us.Directives...)
	for _, uss := range us.UpstreamServers {
		directives = append(directives, uss)
	}

	return directives
}

// NewUpstream creates a new Upstream from a directive.
func NewUpstream(directive IDirective) (*Upstream, error) {
	parameters := directive.GetParameters()
	if len(parameters) == 0 {
		return nil, errors.New("upstream directive requires a name parameter")
	}

	us := &Upstream{
		UpstreamName: parameters[0].GetValue(), //first parameter of the directive is the upstream name
	}

	if directive.GetBlock() == nil {
		return nil, errors.New("missing upstream block")
	}

	if len(directive.GetBlock().GetDirectives()) > 0 {
		for _, d := range directive.GetBlock().GetDirectives() {
			if d.GetName() == "server" {
				uss, err := NewUpstreamServer(d)
				if err != nil {
					return nil, err
				}
				uss.SetParent(us)
				uss.SetLine(d.GetLine())
				us.UpstreamServers = append(us.UpstreamServers, uss)
			} else {
				us.Directives = append(us.Directives, d)
			}
		}
	}

	us.Comment = directive.GetComment()
	us.InlineComment = directive.GetInlineComment()

	return us, nil
}

// AddServer adds a server to the upstream.
func (us *Upstream) AddServer(server *UpstreamServer) {
	us.UpstreamServers = append(us.UpstreamServers, server)
}

// GetCodeBlock returns the literal code block.
func (us *Upstream) GetCodeBlock() string {
	return ""
}

// FindDirectives finds directives in the block recursively.
func (us *Upstream) FindDirectives(directiveName string) []IDirective {
	directives := make([]IDirective, 0)
	for _, directive := range us.Directives {
		if directive.GetName() == directiveName {
			directives = append(directives, directive)
		}
		if directive.GetBlock() != nil {
			directives = append(directives, directive.GetBlock().FindDirectives(directiveName)...)
		}
	}

	return directives
}
