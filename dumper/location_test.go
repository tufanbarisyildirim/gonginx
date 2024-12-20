package dumper

import (
	"testing"

	"github.com/tufanbarisyildirim/gonginx/config"
)

func TestLocation_ToString(t *testing.T) {
	t.Parallel()
	type fields struct {
		Directive *config.Directive
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
				Directive: &config.Directive{
					Name:       "location",
					Parameters: []config.Parameter{{Value: "/admin"}},
					Block: &config.Block{
						Directives: make([]config.IDirective, 0),
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
			l := &config.Location{
				Directive: tt.fields.Directive,
				Modifier:  tt.fields.Modifier,
				Match:     tt.fields.Match,
			}
			if got := DumpDirective(l, NoIndentStyle); got != tt.want {
				t.Errorf("Location.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
