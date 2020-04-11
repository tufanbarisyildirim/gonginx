package parser

import (
	"testing"

	"github.com/tufanbarisyildirim/gonginx/config"
	"github.com/tufanbarisyildirim/gonginx/token"
	"gotest.tools/v3/assert"
)

func TestParser_CurrFollow(t *testing.T) {
	conf := `
	server { # simple reverse-proxy
	}
	`
	p := NewStringParser(conf)
	//assert.Assert(t, tokens, 1)
	assert.Assert(t, p.curTokenIs(token.Keyword))
	assert.Assert(t, p.followingTokenIs(token.BlockStart))
	assert.Assert(t, p.followingTokenIs(token.BlockStart))
}

func TestParser_Include(t *testing.T) {
	conf := `
	include /etc/ngin/conf.d/mime.types;
	`
	p := NewStringParser(conf)
	c := p.Parse()
	_, ok := c.Statements[0].(config.IncludeStatement) //we expect the first statement to be an include
	assert.Assert(t, ok)
}
