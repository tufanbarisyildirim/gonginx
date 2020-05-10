package config

import (
	"errors"

	"github.com/tufanbarisyildirim/gonginx/dumper"
)

//Upstream represents `upstream{}` block
type Upstream struct {
	UpstreamName    string
	UpstreamServers []*UpstreamServer
	//Directives Other directives in upstream (ip_hash; etc)
	Directives []IDirective
}

//GetName Statement interface
func (us *Upstream) GetName() string {
	return "upstream"
}

//GetParameters upsrema parameters
func (us *Upstream) GetParameters() []string {
	return []string{us.UpstreamName} //the only parameter for an upstream is its name
}

//NewUpstream creaste new upstream from a directive
func NewUpstream(directive Directive) (*Upstream, error) {
	us := &Upstream{
		UpstreamName: directive.Parameters[0], //first parameter of the directive is the upstream name
	}

	if directive.Block == nil {
		return nil, errors.New("missing upstream block")
	}

	if len(directive.Block.Directives) > 0 {
		for _, d := range directive.Block.Directives {
			if d.GetName() == "server" {
				us.UpstreamServers = append(us.UpstreamServers, NewUpstreamServer(d))
			}
		}
	}

	return us, nil
}

//ToString convert it to a string
func (us *Upstream) ToString(style *dumper.Style) string {
	directive := Directive{
		Name:       us.GetName(),
		Parameters: us.GetParameters(),
		Block: &Block{
			Directives: []IDirective{},
		},
	}
	//first add other directives
	if us.Directives != nil {
		for _, d := range us.Directives {
			directive.Block.Directives = append(directive.Block.Directives, d)
		}
	}
	//then upstream
	for _, uss := range us.UpstreamServers {
		directive.Block.Directives = append(directive.Block.Directives, uss)
	}

	return directive.ToString(style)
}

//AddServer add a server to upstream
func (us *Upstream) AddServer(server *UpstreamServer) {
	us.UpstreamServers = append(us.UpstreamServers, server)
}
