package gonginx

import (
	"testing"
)

func TestServer_ToString(t *testing.T) {
	type fields struct {
		Directive *Directive
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
				Directive: &Directive{
					Block: &Block{
						Directives: make([]IDirective, 0),
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
			args: NewStyle(),
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

			if got := DumpDirective(s, tt.args); got != tt.want {
				t.Errorf("Server.ToString() = \"%v\", want \"%v\"", got, tt.want)
			}
		})
	}
}
