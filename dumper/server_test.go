package dumper

import (
	"testing"

	"github.com/tufanbarisyildirim/gonginx/config"
)

func TestServer_ToString(t *testing.T) {
	t.Parallel()
	type fields struct {
		Directive *config.Directive
	}
	tests := []struct {
		name   string
		fields fields
		args   *Style
		want   string
	}{
		{
			name: "empty server block",
			fields: fields{
				Directive: &config.Directive{
					Block: &config.Block{
						Directives: make([]config.IDirective, 0),
					},
					Name: "server",
				},
			},
			args: NoIndentStyle,
			want: "server {\n\n}",
		},
		{
			name: "styled server block with some directives",
			fields: fields{
				Directive: &config.Directive{
					Block: &config.Block{
						Directives: []config.IDirective{
							&config.Directive{
								Name:       "server_name",
								Parameters: []config.Parameter{{Value: "gonginx.dev"}},
							},
							&config.Directive{
								Name:       "root",
								Parameters: []config.Parameter{{Value: "/var/sites/gonginx"}},
							},
						},
					},
					Name: "server",
				},
			},
			args: NewStyle(),
			want: `server {
    server_name gonginx.dev;
    root /var/sites/gonginx;
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := config.NewServer(tt.fields.Directive)
			if err != nil {
				t.Error("NewServer(tt.fields.Directive) failed")
			}

			if got := DumpDirective(s, tt.args); got != tt.want {
				t.Errorf("Server.ToString() = \"%v\", want \"%v\"", got, tt.want)
			}
		})
	}
}
