package parser

import "testing"

func FuzzParserParseNoPanic(f *testing.F) {
	seeds := []string{
		``,
		`user nginx;`,
		`server { listen 80; }`,
		`http { upstream backend { server 127.0.0.1:8080; } }`,
		`location / { content_by_lua_block { local s = "#hash" } }`,
		`include a.conf;`,
		`map $http_upgrade $connection_upgrade { default upgrade; '' close; }`,
		`server { location / { proxy_pass http://backend/; }`,
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, conf string) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("parser panicked for input %q: %v", conf, r)
			}
		}()

		p := NewStringParser(conf, WithSkipValidDirectivesErr())
		_, _ = p.Parse()
	})
}
