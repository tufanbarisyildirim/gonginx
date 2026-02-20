package config

import "testing"

func TestHTTP_GetServerByServerName(t *testing.T) {
	t.Parallel()

	serverA := &Server{
		Block: &Block{
			Directives: []IDirective{
				&Directive{
					Name:       "server_name",
					Parameters: []Parameter{{Value: "example.com"}, {Value: "www.example.com"}},
				},
			},
		},
	}
	serverB := &Server{
		Block: &Block{
			Directives: []IDirective{
				&Directive{
					Name:       "server_name",
					Parameters: []Parameter{{Value: "api.example.com"}},
				},
			},
		},
	}
	http := &HTTP{
		Servers: []*Server{serverA, serverB},
	}

	got := http.GetServerByServerName("www.example.com")
	if got != serverA {
		t.Fatalf("expected serverA, got %v", got)
	}

	got = http.GetServerByServerName("api.example.com")
	if got != serverB {
		t.Fatalf("expected serverB, got %v", got)
	}

	got = http.GetServerByServerName("missing.example.com")
	if got != nil {
		t.Fatalf("expected nil for missing server name, got %v", got)
	}
}
