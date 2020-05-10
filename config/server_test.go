package config

import (
	"testing"

	"github.com/tufanbarisyildirim/gonginx/dumper"
)

func TestServer_ToString(t *testing.T) {
	type fields struct {
		Directive *Directive
	}
	tests := []struct {
		name   string
		fields fields
		args   *dumper.Style
		want   string
	}{
		{
			name: "empty server block",
			fields: fields{
				Directive: &Directive{
					Block: &Block{
						Directives: make([]IDirective, 0),
					},
					Name: "server",
				},
			},
			args: dumper.NoIndentStyle,
			want: "server {\n\n}",
		},
		{
			name: "styled server block with some directives",
			fields: fields{
				Directive: &Directive{
					Block: &Block{
						Directives: []IDirective{
							&Directive{
								Name:       "server_name",
								Parameters: []string{"gonginx.dev"},
							},
							&Directive{
								Name:       "root",
								Parameters: []string{"/var/sites/gonginx"},
							},
						},
					},
					Name: "server",
				},
			},
			args: dumper.NewStyle(),
			want: `server {
    server_name gonginx.dev;
    root /var/sites/gonginx;
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewServer(tt.fields.Directive)
			if err != nil {
				t.Error("NewServer(tt.fields.Directive) failed")
			}
			if got := s.ToString(tt.args); got != tt.want {
				t.Errorf("Server.ToString() = \"%v\", want \"%v\"", got, tt.want)
			}
		})
	}
}
