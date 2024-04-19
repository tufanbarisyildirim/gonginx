package dumper

import (
	"testing"

	"github.com/tufanbarisyildirim/gonginx/config"
	"gotest.tools/v3/assert"
)

func TestConfig_IncludeToString(t *testing.T) {
	t.Parallel()
	include := &config.Include{
		Directive: &config.Directive{
			Name:       "include",
			Parameters: []string{"/etc/nginx/conf.d/*.conf"},
		},
		IncludePath: "/etc/nginx/conf.d/*.conf",
	}
	assert.Equal(t, "include /etc/nginx/conf.d/*.conf;", DumpDirective(include, NoIndentStyle))
	var i interface{} = include
	_, ok := i.(config.IDirective)
	//_, ok2 := i.(IncludeDirective)// TODO(tufan):reactivate here after getting include and file things done
	assert.Assert(t, ok)
	//assert.Assert(t, ok2)
}
