package gonginx

import (
	"strings"
	"testing"
)

func TestUpstream_ToString(t *testing.T) {
	t.Parallel()
	type fields struct {
		Directive *Directive
		Name      string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "empty upstream block with name",
			fields: fields{
				Directive: &Directive{
					Name:       "upstream",
					Parameters: []string{"gonginx_upstream"},
					Block: &Block{
						Directives: make([]IDirective, 0),
					},
				},
			},
			want: "upstream gonginx_upstream {\n\n}",
		},
		{
			name: "empty upstream block with name and upstream server",
			fields: fields{
				Directive: &Directive{
					Name:       "upstream",
					Parameters: []string{"gonginx_upstream"},
					Block: &Block{
						Directives: []IDirective{
							NewUpstreamServerIgnoreErr(&Directive{
								Name:       "server",
								Parameters: []string{"127.0.0.1:9005"},
							}),
						},
					},
				},
			},
			want: "upstream gonginx_upstream {\nserver 127.0.0.1:9005;\n}",
		},
		{
			name: "empty upstream block with name and multi upstream server",
			fields: fields{
				Directive: &Directive{
					Name:       "upstream",
					Parameters: []string{"gonginx_upstream"},
					Block: &Block{
						Directives: []IDirective{
							NewUpstreamServerIgnoreErr(&Directive{
								Name:       "server",
								Parameters: []string{"127.0.0.1:9005"},
							}),
							NewUpstreamServerIgnoreErr(&Directive{
								Name:       "server",
								Parameters: []string{"127.0.0.2:9005"},
							}),
						},
					},
				},
			},
			want: "upstream gonginx_upstream {\nserver 127.0.0.1:9005;\nserver 127.0.0.2:9005;\n}",
		},
		{
			name: "empty upstream block with name and multi upstream server and some flags, params",
			fields: fields{
				Directive: &Directive{
					Name:       "upstream",
					Parameters: []string{"gonginx_upstream"},
					Block: &Block{
						Directives: []IDirective{
							NewUpstreamServerIgnoreErr(&Directive{
								Name:       "server",
								Parameters: []string{"127.0.0.1:9005", "weight=5"},
							}),
							NewUpstreamServerIgnoreErr(&Directive{
								Name:       "server",
								Parameters: []string{"127.0.0.2:9005", "weight=4", "down"},
							}),
						},
					},
				},
			},
			want: "upstream gonginx_upstream {\nserver 127.0.0.1:9005 weight=5;\nserver 127.0.0.2:9005 weight=4 down;\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us, err := NewUpstream(tt.fields.Directive)
			if err != nil {
				t.Error("Failed to create NewUpstream(*tt.fields.Directive)")
			}
			if got := DumpDirective(us, NoIndentStyle); got != tt.want {
				t.Errorf("Upstream.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func NewUpstreamServerIgnoreErr(directive IDirective) *UpstreamServer {
	server, _ := NewUpstreamServer(directive)
	return server
}

func TestUpstream_AddServer(t *testing.T) {
	t.Parallel()
	type fields struct {
		UpstreamName    string
		UpstreamServers []*UpstreamServer
		Directives      []IDirective
	}
	type args struct {
		server *UpstreamServer
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		toString string
	}{
		{
			name: "add simple server",
			fields: fields{
				UpstreamName: "my_backend",
				UpstreamServers: []*UpstreamServer{
					{
						Address: "127.0.0.1:8080",
						Flags:   []string{"backup"},
						Parameters: map[string]string{
							"weight": "1",
						},
					},
				},
			},
			args: args{
				server: &UpstreamServer{
					Address: "backend2.gonginx.org:8090",
					Flags:   []string{"resolve"},
					Parameters: map[string]string{
						"fail_timeout": "5s",
						"slow_start":   "30s",
					},
				},
			},
			toString: `upstream my_backend {
server 127.0.0.1:8080 weight=1 backup;
server backend2.gonginx.org:8090 fail_timeout=5s slow_start=30s resolve;
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &Upstream{
				UpstreamName:    tt.fields.UpstreamName,
				UpstreamServers: tt.fields.UpstreamServers,
				Directives:      tt.fields.Directives,
			}
			us.AddServer(tt.args.server)
			if got := DumpDirective(us, NoIndentStyle); got != tt.toString {
				t.Errorf("us.ToString() = `%v`, want `%v`", strings.ReplaceAll(got, " ", "."), strings.ReplaceAll(tt.toString, " ", "."))
			}
		})
	}
}
