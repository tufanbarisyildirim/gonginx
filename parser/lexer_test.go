package parser

import (
	"encoding/json"
	"testing"

	"github.com/tufanbarisyildirim/gonginx/parser/token"
	"gotest.tools/v3/assert"
)

func TestScanner_Lex(t *testing.T) {
	actual := lex(`
server { # simple reverse-proxy
    listen       80;
    server_name  gonginx.com www.gonginx.com;
    access_log   logs/gonginx.access.log  main;

    # serve static files
    location ~ ^/(images|javascript|js|css|flash|media|static)/  {
      root    /var/www/virtual/gonginx/;
      expires 30d;
    }

    # pass requests for dynamic content
    location / {
      proxy_pass      http://127.0.0.1:8080;
      proxy_set_header   X-Real-IP        $remote_addr;
    }
  }
include /etc/nginx/conf.d/*.conf;
`).all()

	var expect = token.Tokens{
		{Type: token.Keyword, Literal: "server", Line: 2, Column: 1},
		{Type: token.BlockStart, Literal: "{", Line: 2, Column: 8},
		{Type: token.Comment, Literal: "# simple reverse-proxy", Line: 2, Column: 10},
		{Type: token.Keyword, Literal: "listen", Line: 3, Column: 5},
		{Type: token.Keyword, Literal: "80", Line: 3, Column: 18},
		{Type: token.Semicolon, Literal: ";", Line: 3, Column: 20},
		{Type: token.Keyword, Literal: "server_name", Line: 4, Column: 5},
		{Type: token.Keyword, Literal: "gonginx.com", Line: 4, Column: 18},
		{Type: token.Keyword, Literal: "www.gonginx.com", Line: 4, Column: 30},
		{Type: token.Semicolon, Literal: ";", Line: 4, Column: 45},
		{Type: token.Keyword, Literal: "access_log", Line: 5, Column: 5},
		{Type: token.Keyword, Literal: "logs/gonginx.access.log", Line: 5, Column: 18},
		{Type: token.Keyword, Literal: "main", Line: 5, Column: 43},
		{Type: token.Semicolon, Literal: ";", Line: 5, Column: 47},
		{Type: token.Comment, Literal: "# serve static files", Line: 7, Column: 5},
		{Type: token.Keyword, Literal: "location", Line: 8, Column: 5},
		{Type: token.Keyword, Literal: "~", Line: 8, Column: 14},
		{Type: token.Keyword, Literal: "^/(images|javascript|js|css|flash|media|static)/", Line: 8, Column: 16},
		{Type: token.BlockStart, Literal: "{", Line: 8, Column: 66},
		{Type: token.Keyword, Literal: "root", Line: 9, Column: 7},
		{Type: token.Keyword, Literal: "/var/www/virtual/gonginx/", Line: 9, Column: 15},
		{Type: token.Semicolon, Literal: ";", Line: 9, Column: 40},
		{Type: token.Keyword, Literal: "expires", Line: 10, Column: 7},
		{Type: token.Keyword, Literal: "30d", Line: 10, Column: 15},
		{Type: token.Semicolon, Literal: ";", Line: 10, Column: 18},
		{Type: token.BlockEnd, Literal: "}", Line: 11, Column: 5},
		{Type: token.Comment, Literal: "# pass requests for dynamic content", Line: 13, Column: 5},
		{Type: token.Keyword, Literal: "location", Line: 14, Column: 5},
		{Type: token.Keyword, Literal: "/", Line: 14, Column: 14},
		{Type: token.BlockStart, Literal: "{", Line: 14, Column: 16},
		{Type: token.Keyword, Literal: "proxy_pass", Line: 15, Column: 7},
		{Type: token.Keyword, Literal: "http://127.0.0.1:8080", Line: 15, Column: 23},
		{Type: token.Semicolon, Literal: ";", Line: 15, Column: 44},
		{Type: token.Keyword, Literal: "proxy_set_header", Line: 16, Column: 7},
		{Type: token.Keyword, Literal: "X-Real-IP", Line: 16, Column: 26},
		{Type: token.Variable, Literal: "$remote_addr", Line: 16, Column: 43},
		{Type: token.Semicolon, Literal: ";", Line: 16, Column: 55},
		{Type: token.BlockEnd, Literal: "}", Line: 17, Column: 5},
		{Type: token.BlockEnd, Literal: "}", Line: 18, Column: 3},
		{Type: token.Keyword, Literal: "include", Line: 19, Column: 1},
		{Type: token.Keyword, Literal: "/etc/nginx/conf.d/*.conf", Line: 19, Column: 9},
		{Type: token.Semicolon, Literal: ";", Line: 19, Column: 33},
	}
	//assert.Equal(t, actual, 1)
	tokenString, err := json.Marshal(actual)
	assert.NilError(t, err)
	expectJSON, err := json.Marshal(expect)
	assert.NilError(t, err)

	//assert.Assert(t, tokens, 1)
	assert.Assert(t, actual.EqualTo(expect))
	assert.Equal(t, string(tokenString), string(expectJSON))
	assert.Equal(t, len(actual), len(expect))
}
