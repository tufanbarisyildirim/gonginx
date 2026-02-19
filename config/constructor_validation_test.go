package config

import (
	"errors"
	"testing"

	"gotest.tools/v3/assert"
)

func TestNewInclude_Validation(t *testing.T) {
	t.Parallel()

	_, err := NewInclude(&Directive{Name: "include"})
	assert.ErrorContains(t, err, "requires exactly 1 parameter")

	_, err = NewInclude(&Directive{
		Name: "include",
		Parameters: []Parameter{
			{Value: "a.conf"},
			{Value: "b.conf"},
		},
	})
	assert.ErrorContains(t, err, "requires exactly 1 parameter")

	_, err = NewInclude(&Directive{
		Name:       "include",
		Parameters: []Parameter{{Value: "a.conf"}},
		Block:      &Block{},
	})
	assert.ErrorContains(t, err, "cannot have a block")
}

func TestNewUpstream_Validation(t *testing.T) {
	t.Parallel()

	_, err := NewUpstream(&Directive{
		Name:  "upstream",
		Block: &Block{},
	})
	assert.ErrorContains(t, err, "requires a name parameter")
}

func TestNewLocation_TypeErrorMessage(t *testing.T) {
	t.Parallel()

	_, err := NewLocation(&Server{})
	assert.ErrorContains(t, err, "location directive type error")
}

func TestConfig_FindUpstreams_SkipsUnexpectedTypes(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Block: &Block{
			Directives: []IDirective{
				&Directive{
					Name:       "upstream",
					Parameters: []Parameter{{Value: "backend"}},
					Block:      &Block{},
				},
			},
		},
	}

	upstreams := cfg.FindUpstreams()
	assert.Equal(t, len(upstreams), 0)
}

func TestConfig_FindUpstreamsStrict_ReturnsTypedError(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Block: &Block{
			Directives: []IDirective{
				&Directive{
					Name:       "upstream",
					Parameters: []Parameter{{Value: "backend"}},
					Block:      &Block{},
				},
			},
		},
	}

	_, err := cfg.FindUpstreamsStrict()
	assert.Assert(t, err != nil)

	var typeErr *UnexpectedUpstreamTypeError
	assert.Assert(t, errors.As(err, &typeErr))
	assert.Equal(t, typeErr.Index, 0)
}

func TestConfig_FindUpstreamsStrict_Success(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Block: &Block{
			Directives: []IDirective{
				&Upstream{
					UpstreamName: "backend",
				},
			},
		},
	}

	upstreams, err := cfg.FindUpstreamsStrict()
	assert.NilError(t, err)
	assert.Equal(t, len(upstreams), 1)
	assert.Equal(t, upstreams[0].UpstreamName, "backend")
}
