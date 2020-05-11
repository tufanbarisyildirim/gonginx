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

//GetBlock upstream does not have block
func (us *Upstream) GetBlock() *Block {
	return us.ToDirective().Block
}

//NewUpstream creaste new upstream from a directive
func NewUpstream(directive IDirective) (*Upstream, error) {
	parameters := directive.GetParameters()
	us := &Upstream{
		UpstreamName: parameters[0], //first parameter of the directive is the upstream name
	}

	if directive.GetBlock() == nil {
		return nil, errors.New("missing upstream block")
	}

	if len(directive.GetBlock().Directives) > 0 {
		for _, d := range directive.GetBlock().Directives {
			if d.GetName() == "server" {
				us.UpstreamServers = append(us.UpstreamServers, NewUpstreamServer(d))
			}
		}
	}

	return us, nil
}

//ToDirective get upstream as a directive
func (us *Upstream) ToDirective() *Directive {
	directive := &Directive{
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

	return directive
}

//ToString convert it to a string
func (us *Upstream) ToString(style *dumper.Style) string {
	return us.ToDirective().ToString(style)
}

//AddServer add a server to upstream
func (us *Upstream) AddServer(server *UpstreamServer) {
	us.UpstreamServers = append(us.UpstreamServers, server)
}
