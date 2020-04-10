package parser

import (
	"testing"

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
}
