package config

import (
	"testing"

	"github.com/tufanbarisyildirim/gonginx/dumper"
)

func TestLocation_ToString(t *testing.T) {
	type fields struct {
		Directive *Directive
		Modifier  string
		Match     string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "empty bloc and one match location empty block",
			fields: fields{
				Directive: &Directive{
					Name:       "location",
					Parameters: []string{"/admin"},
					Block: &Block{
						Statements: make([]Statement, 0),
					},
				},
				Modifier: "",
				Match:    "/admin",
			},
			want: "location /admin {\n\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Location{
				Directive: tt.fields.Directive,
				Modifier:  tt.fields.Modifier,
				Match:     tt.fields.Match,
			}
			if got := l.ToString(dumper.NoIndentStyle); got != tt.want {
				t.Errorf("Location.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
