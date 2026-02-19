package parser

import (
	"encoding/json"
	"testing"

	"gotest.tools/v3/assert"

	"github.com/tufanbarisyildirim/gonginx/parser/token"
)

func TestScanner_Lex(t *testing.T) {
	t.Parallel()
	conf := `
server { # simple reverse-proxy
    listen       80;
    server_name  gonginx.com www.gonginx.com;
    access_log   logs/gonginx.access.log  main;

    # serve static files
    location ~ ^/(images|javascript|js|css|flash|media|static)/  {
	  root    /var/www/virtual/gonginx/;
	  fastcgi_param  SERVER_SOFTWARE    nginx/$nginx_version/$server_name;
      expires 30d;
    }

    # pass requests for dynamic content
    location / {
      proxy_pass      http://127.0.0.1:8080;
      proxy_set_header   X-Real-IP        $remote_addr;
    }
  }
include /etc/nginx/conf.d/*.conf;
directive "with a quoted string\t \r\n \\ with some escaped thing s\" good.";
#also cmment right before eof`
	actual := lex(conf).all()

	var expect = token.Tokens{
		{Type: token.EndOfLine, Literal: "\n", Line: 1, Column: 0},
		{Type: token.Keyword, Literal: "server", Line: 2, Column: 1},
		{Type: token.BlockStart, Literal: "{", Line: 2, Column: 8},
		{Type: token.Comment, Literal: "# simple reverse-proxy", Line: 2, Column: 10},
		{Type: token.EndOfLine, Literal: "\n", Line: 2, Column: 32},
		{Type: token.Keyword, Literal: "listen", Line: 3, Column: 5},
		{Type: token.Keyword, Literal: "80", Line: 3, Column: 18},
		{Type: token.Semicolon, Literal: ";", Line: 3, Column: 20},
		{Type: token.EndOfLine, Literal: "\n", Line: 3, Column: 21},
		{Type: token.Keyword, Literal: "server_name", Line: 4, Column: 5},
		{Type: token.Keyword, Literal: "gonginx.com", Line: 4, Column: 18},
		{Type: token.Keyword, Literal: "www.gonginx.com", Line: 4, Column: 30},
		{Type: token.Semicolon, Literal: ";", Line: 4, Column: 45},
		{Type: token.EndOfLine, Literal: "\n", Line: 4, Column: 46},
		{Type: token.Keyword, Literal: "access_log", Line: 5, Column: 5},
		{Type: token.Keyword, Literal: "logs/gonginx.access.log", Line: 5, Column: 18},
		{Type: token.Keyword, Literal: "main", Line: 5, Column: 43},
		{Type: token.Semicolon, Literal: ";", Line: 5, Column: 47},
		{Type: token.EndOfLine, Literal: "\n", Line: 5, Column: 48},
		{Type: token.EndOfLine, Literal: "\n", Line: 6, Column: 1},
		{Type: token.Comment, Literal: "# serve static files", Line: 7, Column: 5},
		{Type: token.EndOfLine, Literal: "\n", Line: 7, Column: 25},
		{Type: token.Keyword, Literal: "location", Line: 8, Column: 5},
		{Type: token.Keyword, Literal: "~", Line: 8, Column: 14},
		{Type: token.Keyword, Literal: "^/(images|javascript|js|css|flash|media|static)/", Line: 8, Column: 16},
		{Type: token.BlockStart, Literal: "{", Line: 8, Column: 66},
		{Type: token.EndOfLine, Literal: "\n", Line: 8, Column: 67},
		{Type: token.Keyword, Literal: "root", Line: 9, Column: 4},
		{Type: token.Keyword, Literal: "/var/www/virtual/gonginx/", Line: 9, Column: 12},
		{Type: token.Semicolon, Literal: ";", Line: 9, Column: 37},
		{Type: token.EndOfLine, Literal: "\n", Line: 9, Column: 38},
		{Type: token.Keyword, Literal: "fastcgi_param", Line: 10, Column: 4},
		{Type: token.Keyword, Literal: "SERVER_SOFTWARE", Line: 10, Column: 19},
		{Type: token.Keyword, Literal: "nginx/$nginx_version/$server_name", Line: 10, Column: 38},
		{Type: token.Semicolon, Literal: ";", Line: 10, Column: 71},
		{Type: token.EndOfLine, Literal: "\n", Line: 10, Column: 72},
		{Type: token.Keyword, Literal: "expires", Line: 11, Column: 7},
		{Type: token.Keyword, Literal: "30d", Line: 11, Column: 15},
		{Type: token.Semicolon, Literal: ";", Line: 11, Column: 18},
		{Type: token.EndOfLine, Literal: "\n", Line: 11, Column: 19},
		{Type: token.BlockEnd, Literal: "}", Line: 12, Column: 5},
		{Type: token.EndOfLine, Literal: "\n", Line: 12, Column: 6},
		{Type: token.EndOfLine, Literal: "\n", Line: 13, Column: 1},
		{Type: token.Comment, Literal: "# pass requests for dynamic content", Line: 14, Column: 5},
		{Type: token.EndOfLine, Literal: "\n", Line: 14, Column: 40},
		{Type: token.Keyword, Literal: "location", Line: 15, Column: 5},
		{Type: token.Keyword, Literal: "/", Line: 15, Column: 14},
		{Type: token.BlockStart, Literal: "{", Line: 15, Column: 16},
		{Type: token.EndOfLine, Literal: "\n", Line: 15, Column: 17},
		{Type: token.Keyword, Literal: "proxy_pass", Line: 16, Column: 7},
		{Type: token.Keyword, Literal: "http://127.0.0.1:8080", Line: 16, Column: 23},
		{Type: token.Semicolon, Literal: ";", Line: 16, Column: 44},
		{Type: token.EndOfLine, Literal: "\n", Line: 16, Column: 45},
		{Type: token.Keyword, Literal: "proxy_set_header", Line: 17, Column: 7},
		{Type: token.Keyword, Literal: "X-Real-IP", Line: 17, Column: 26},
		{Type: token.Keyword, Literal: "$remote_addr", Line: 17, Column: 43},
		{Type: token.Semicolon, Literal: ";", Line: 17, Column: 55},
		{Type: token.EndOfLine, Literal: "\n", Line: 17, Column: 56},
		{Type: token.BlockEnd, Literal: "}", Line: 18, Column: 5},
		{Type: token.EndOfLine, Literal: "\n", Line: 18, Column: 6},
		{Type: token.BlockEnd, Literal: "}", Line: 19, Column: 3},
		{Type: token.EndOfLine, Literal: "\n", Line: 19, Column: 4},
		{Type: token.Keyword, Literal: "include", Line: 20, Column: 1},
		{Type: token.Keyword, Literal: "/etc/nginx/conf.d/*.conf", Line: 20, Column: 9},
		{Type: token.Semicolon, Literal: ";", Line: 20, Column: 33},
		{Type: token.EndOfLine, Literal: "\n", Line: 20, Column: 34},
		{Type: token.Keyword, Literal: "directive", Line: 21, Column: 1},
		{Type: token.QuotedString, Literal: "\"with a quoted string\\t \\r\\n \\\\ with some escaped thing s\\\" good.\"", Line: 21, Column: 11},
		{Type: token.Semicolon, Literal: ";", Line: 21, Column: 77},
		{Type: token.EndOfLine, Literal: "\n", Line: 21, Column: 78},
		{Type: token.Comment, Literal: "#also cmment right before eof", Line: 22, Column: 1},
	}
	//assert.Equal(t, actual, 1)
	tokenString, err := json.Marshal(actual)
	assert.NilError(t, err)
	expectJSON, err := json.Marshal(expect)
	assert.NilError(t, err)

	assert.NilError(t, actual.Diff(expect))

	//assert.Assert(t, tokens, 1)
	assert.Equal(t, string(tokenString), string(expectJSON))
	assert.Assert(t, actual.EqualTo(expect))
	assert.Equal(t, len(actual), len(expect))
}

func TestScanner_LexUnclosedQuoteError(t *testing.T) {
	t.Parallel()

	l := lex(`
	server { 
	directive "with an unclosed quote \t \r\n \\ with some escaped thing s\" good.;
	`)
	_ = l.all()
	assert.ErrorContains(t, l.Err, "quoted string")
	assert.ErrorContains(t, l.Err, "line")
	assert.ErrorContains(t, l.Err, "column")
}

func TestScanner_LexLuaCode(t *testing.T) {
	conf := `
server {
  location = /foo {
    rewrite_by_lua_block {
      res = ngx.location.capture("/memc",
        { args = { cmd = "incr", key = ngx.var.uri } } # comment contained unexpect '{'
         # comment contained unexpect '}' 
      )
      t = { key="foo", val="bar" }
    }
  }
}`
	actual := lex(conf).all()
	var expect = token.Tokens{
		{Type: token.EndOfLine, Literal: "\n", Line: 1, Column: 0},
		{Type: token.Keyword, Literal: "server", Line: 2, Column: 1},
		{Type: token.BlockStart, Literal: "{", Line: 2, Column: 8},
		{Type: token.EndOfLine, Literal: "\n", Line: 2, Column: 9},
		{Type: token.Keyword, Literal: "location", Line: 3, Column: 3},
		{Type: token.Keyword, Literal: "=", Line: 3, Column: 12},
		{Type: token.Keyword, Literal: "/foo", Line: 3, Column: 14},
		{Type: token.BlockStart, Literal: "{", Line: 3, Column: 19},
		{Type: token.EndOfLine, Literal: "\n", Line: 3, Column: 20},
		{Type: token.Keyword, Literal: "rewrite_by_lua_block", Line: 4, Column: 5},
		{Type: token.BlockStart, Literal: "{", Line: 4, Column: 26},
		{Type: token.LuaCode, Literal: `
      res = ngx.location.capture("/memc",
        { args = { cmd = "incr", key = ngx.var.uri } } # comment contained unexpect '{'
         # comment contained unexpect '}' 
      )
      t = { key="foo", val="bar" }
    `, Line: 4, Column: 27},
		{Type: token.BlockEnd, Literal: "}", Line: 10, Column: 6},
		{Type: token.EndOfLine, Literal: "\n", Line: 10, Column: 7},
		{Type: token.BlockEnd, Literal: "}", Line: 11, Column: 3},
		{Type: token.EndOfLine, Literal: "\n", Line: 11, Column: 4},
		{Type: token.BlockEnd, Literal: "}", Line: 12, Column: 1},
	}
	tokenString, err := json.Marshal(actual)
	assert.NilError(t, err)
	expectJSON, err := json.Marshal(expect)
	assert.NilError(t, err)
	assert.NilError(t, actual.Diff(expect))
	assert.Equal(t, string(tokenString), string(expectJSON))
	assert.Equal(t, len(actual), len(expect))
}
