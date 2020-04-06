package parser

import (
	"encoding/json"
	"github.com/tufanbarisyildirim/gonginx/token"
	"gotest.tools/v3/assert"
	"testing"
)

func TestScanner_Set(t *testing.T) {
	tokens := Parse(`
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
    }
  }
`).All()

	var actual = token.Tokens{
		{Type: token.Keyword, Literal: "server", Line: 2, Column: 1},
		{Type: token.OpenBrace, Literal: "{", Line: 2, Column: 8},
		{Type: token.Comment, Literal: "# simple reverse-proxy", Line: 2, Column: 10},
		{Type: token.Keyword, Literal: "listen", Line: 5, Column: 5},
		{Type: token.Keyword, Literal: "80", Line: 5, Column: 18},
		{Type: token.Semicolon, Literal: ";", Line: 5, Column: 20},
		{Type: token.Keyword, Literal: "server_name", Line: 7, Column: 5},
		{Type: token.Keyword, Literal: "gonginx.com", Line: 7, Column: 18},
		{Type: token.Keyword, Literal: "www.gonginx.com", Line: 7, Column: 30},
		{Type: token.Semicolon, Literal: ";", Line: 7, Column: 45},
		{Type: token.Keyword, Literal: "access_log", Line: 9, Column: 5},
		{Type: token.Keyword, Literal: "logs/gonginx.access.log", Line: 9, Column: 18},
		{Type: token.Keyword, Literal: "main", Line: 9, Column: 43},
		{Type: token.Semicolon, Literal: ";", Line: 9, Column: 47},
		{Type: token.Comment, Literal: "# serve static files", Line: 12, Column: 5},
		{Type: token.Keyword, Literal: "location", Line: 15, Column: 5},
		{Type: token.Keyword, Literal: "~", Line: 15, Column: 14},
		{Type: token.Keyword, Literal: "^/(images|javascript|js|css|flash|media|static)/", Line: 15, Column: 16},
		{Type: token.OpenBrace, Literal: "{", Line: 15, Column: 66},
		{Type: token.Keyword, Literal: "root", Line: 17, Column: 7},
		{Type: token.Keyword, Literal: "/var/www/virtual/gonginx/", Line: 17, Column: 15},
		{Type: token.Semicolon, Literal: ";", Line: 17, Column: 40},
		{Type: token.Keyword, Literal: "expires", Line: 19, Column: 7},
		{Type: token.Keyword, Literal: "30d", Line: 19, Column: 15},
		{Type: token.Semicolon, Literal: ";", Line: 19, Column: 18},
		{Type: token.CloseBrace, Literal: "}", Line: 21, Column: 5},
		{Type: token.Comment, Literal: "# pass requests for dynamic content", Line: 24, Column: 5},
		{Type: token.Keyword, Literal: "location", Line: 27, Column: 5},
		{Type: token.Keyword, Literal: "/", Line: 27, Column: 14},
		{Type: token.OpenBrace, Literal: "{", Line: 27, Column: 16},
		{Type: token.Keyword, Literal: "proxy_pass", Line: 29, Column: 7},
		{Type: token.Keyword, Literal: "http://127.0.0.1:8080", Line: 29, Column: 23},
		{Type: token.Semicolon, Literal: ";", Line: 29, Column: 44},
		{Type: token.CloseBrace, Literal: "}", Line: 31, Column: 5},
		{Type: token.CloseBrace, Literal: "}", Line: 33, Column: 3},
	}
	tokenString, err := json.Marshal(tokens)
	assert.NilError(t, err)
	expect, err := json.Marshal(actual)
	assert.NilError(t, err)

	assert.Assert(t, tokens.EqualTo(actual))
	assert.Equal(t, string(tokenString), string(expect))
	assert.Equal(t, len(tokens), len(actual))
}
