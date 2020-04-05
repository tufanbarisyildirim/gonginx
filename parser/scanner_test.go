package parser

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestScanner_Set(t *testing.T) {
	tokens := Parse(` server 
{ 
	hello this is block; # with a comment
	oh yes;
} `).All()
	assert.Equal(t, len(tokens), 12)
}
