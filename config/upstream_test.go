package config

import (
	"testing"

	"github.com/tufanbarisyildirim/gonginx/dumper"
)

func TestUpstream_ToString(t *testing.T) {
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
							NewUpstreamServer(&Directive{
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
							NewUpstreamServer(&Directive{
								Name:       "server",
								Parameters: []string{"127.0.0.1:9005"},
							}),
							NewUpstreamServer(&Directive{
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
							NewUpstreamServer(&Directive{
								Name:       "server",
								Parameters: []string{"127.0.0.1:9005", "weight=5"},
							}),
							NewUpstreamServer(&Directive{
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
			us, err := NewUpstream(*tt.fields.Directive)
			if err != nil {
				t.Error("Failed to create NewUpstream(*tt.fields.Directive)")
			}
			if got := us.ToString(dumper.NoIndentStyle); got != tt.want {
				t.Errorf("Upstream.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
