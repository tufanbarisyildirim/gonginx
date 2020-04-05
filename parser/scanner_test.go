package parser

import (
	"encoding/json"
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

	var actual = []Token{
		{Type: keyword, Literal: "server", Line: 2, LineEnd: 0, Column: 1, ColumnEnd: 0},
		{Type: openBrace, Literal: "{", Line: 2, LineEnd: 0, Column: 8, ColumnEnd: 0},
		{Type: comment, Literal: "# simple reverse-proxy", Line: 2, LineEnd: 0, Column: 10, ColumnEnd: 0},
		{Type: keyword, Literal: "listen", Line: 5, LineEnd: 0, Column: 5, ColumnEnd: 0},
		{Type: keyword, Literal: "80", Line: 5, LineEnd: 0, Column: 18, ColumnEnd: 0},
		{Type: semiColon, Literal: ";", Line: 5, LineEnd: 0, Column: 20, ColumnEnd: 0},
		{Type: keyword, Literal: "server_name", Line: 7, LineEnd: 0, Column: 5, ColumnEnd: 0},
		{Type: keyword, Literal: "gonginx.com", Line: 7, LineEnd: 0, Column: 18, ColumnEnd: 0},
		{Type: keyword, Literal: "www.gonginx.com", Line: 7, LineEnd: 0, Column: 30, ColumnEnd: 0},
		{Type: semiColon, Literal: ";", Line: 7, LineEnd: 0, Column: 45, ColumnEnd: 0},
		{Type: keyword, Literal: "access_log", Line: 9, LineEnd: 0, Column: 5, ColumnEnd: 0},
		{Type: keyword, Literal: "logs/gonginx.access.log", Line: 9, LineEnd: 0, Column: 18, ColumnEnd: 0},
		{Type: keyword, Literal: "main", Line: 9, LineEnd: 0, Column: 43, ColumnEnd: 0},
		{Type: semiColon, Literal: ";", Line: 9, LineEnd: 0, Column: 47, ColumnEnd: 0},
		{Type: comment, Literal: "# serve static files", Line: 12, LineEnd: 0, Column: 5, ColumnEnd: 0},
		{Type: keyword, Literal: "location", Line: 15, LineEnd: 0, Column: 5, ColumnEnd: 0},
		{Type: keyword, Literal: "~", Line: 15, LineEnd: 0, Column: 14, ColumnEnd: 0},
		{Type: keyword, Literal: "^/(images|javascript|js|css|flash|media|static)/", Line: 15, LineEnd: 0, Column: 16, ColumnEnd: 0},
		{Type: openBrace, Literal: "{", Line: 15, LineEnd: 0, Column: 66, ColumnEnd: 0},
		{Type: keyword, Literal: "root", Line: 17, LineEnd: 0, Column: 7, ColumnEnd: 0},
		{Type: keyword, Literal: "/var/www/virtual/gonginx/", Line: 17, LineEnd: 0, Column: 15, ColumnEnd: 0},
		{Type: semiColon, Literal: ";", Line: 17, LineEnd: 0, Column: 40, ColumnEnd: 0},
		{Type: keyword, Literal: "expires", Line: 19, LineEnd: 0, Column: 7, ColumnEnd: 0},
		{Type: keyword, Literal: "30d", Line: 19, LineEnd: 0, Column: 15, ColumnEnd: 0},
		{Type: semiColon, Literal: ";", Line: 19, LineEnd: 0, Column: 18, ColumnEnd: 0},
		{Type: closeBrace, Literal: "}", Line: 21, LineEnd: 0, Column: 5, ColumnEnd: 0},
		{Type: comment, Literal: "# pass requests for dynamic content", Line: 24, LineEnd: 0, Column: 5, ColumnEnd: 0},
		{Type: keyword, Literal: "location", Line: 27, LineEnd: 0, Column: 5, ColumnEnd: 0},
		{Type: keyword, Literal: "/", Line: 27, LineEnd: 0, Column: 14, ColumnEnd: 0},
		{Type: openBrace, Literal: "{", Line: 27, LineEnd: 0, Column: 16, ColumnEnd: 0},
		{Type: keyword, Literal: "proxy_pass", Line: 29, LineEnd: 0, Column: 7, ColumnEnd: 0},
		{Type: keyword, Literal: "http://127.0.0.1:8080", Line: 29, LineEnd: 0, Column: 23, ColumnEnd: 0},
		{Type: semiColon, Literal: ";", Line: 29, LineEnd: 0, Column: 44, ColumnEnd: 0},
		{Type: closeBrace, Literal: "}", Line: 31, LineEnd: 0, Column: 5, ColumnEnd: 0},
		{Type: closeBrace, Literal: "}", Line: 33, LineEnd: 0, Column: 3, ColumnEnd: 0},
	}
	tokenString, err := json.Marshal(tokens)
	assert.NilError(t, err)
	expect, err := json.Marshal(actual)
	assert.NilError(t, err)

	assert.Equal(t, string(tokenString), string(expect))
	assert.Equal(t, len(tokens), len(actual))
}
