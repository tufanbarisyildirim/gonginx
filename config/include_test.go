package config

import (
	"testing"

	"github.com/tufanbarisyildirim/gonginx/token"
	"gotest.tools/v3/assert"
)

func TestConfig_IncludeToString(t *testing.T) {
	include := &Include{
		token: token.Token{
			Type:    token.Keyword,
			Literal: "include",
		},
		IncludePath: "/etc/nginx/conf.d/*.conf",
	}
	assert.Equal(t, "include /etc/nginx/conf.d/*.conf;", include.ToString())
	var i interface{} = include
	_, ok := i.(Statement)
	_, ok2 := i.(IncludeStatement)
	assert.Assert(t, ok)
	assert.Assert(t, ok2)
}
