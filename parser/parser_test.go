package parser

import (
	"testing"

	"github.com/tufanbarisyildirim/gonginx/config"
	"github.com/tufanbarisyildirim/gonginx/parser/token"
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

func TestParser_Include(t *testing.T) {
	conf := `
	include /etc/ngin/conf.d/mime.types;
	`
	p := NewStringParser(conf)
	c := p.Parse()
	_, ok := c.Statements[0].(config.IncludeStatement) //we expect the first statement to be an include
	assert.Assert(t, ok)
}

func TestParser_UnendedInclude(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	NewParserFromLexer(
		lex(`
	server { 
	include /but/no/semicolon {}
	`)).Parse()
}

func TestParser_ParseFromFile(t *testing.T) {
	_, err := NewParser("../full-example/nginx.conf")
	assert.NilError(t, err)
	_, err2 := NewParser("../full-example/nginx.conf-not-found")
	assert.ErrorContains(t, err2, "no such file or directory")
}

func TestParser_MultiParamDirecive(t *testing.T) {

	NewParserFromLexer(
		lex(`
	server { 
	a_directive has multi params /and/ends;
	`)).Parse()
}

func TestParser_UnendedMultiParams(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	NewParserFromLexer(
		lex(`
	server { 
	a_driective with mutli params /but/no/semicolon/to/panic }
	`)).Parse()
}
