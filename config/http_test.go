package config

import (
	"strings"
	"testing"

	"github.com/tufanbarisyildirim/gonginx/dumper"
)

func TestHttp_ToString(t *testing.T) {
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
			name: "empty http block",
			fields: fields{
				Directive: &Directive{
					Block: &Block{
						Directives: make([]IDirective, 0),
					},
					Name: "http",
				},
			},
			args: dumper.NoIndentStyle,
			want: "http {\n\n}",
		},
		{
			name: "styled http block with some directives",
			fields: fields{
				Directive: &Directive{
					Block: &Block{
						Directives: []IDirective{
							&Directive{
								Name:       "access_log",
								Parameters: []string{"logs/access.log", "main"},
							},
							&Directive{
								Name:       "default_type",
								Parameters: []string{"application/octet-stream"},
							},
						},
					},
					Name: "http",
				},
			},
			args: dumper.NewStyle(),
			want: `http {
    access_log logs/access.log main;
    default_type application/octet-stream;
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewHttp(tt.fields.Directive)
			if err != nil {
				t.Error("NewHttp(tt.fields.Directive) failed")
			}
			if got := s.ToString(tt.args); got != tt.want {
				t.Errorf("Http.ToString() = \"%v\", want \"%v\"", strings.ReplaceAll(got, " ", "."), strings.ReplaceAll(tt.want, " ", "."))
			}
		})
	}
}
