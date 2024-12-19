package dumper

import (
	"testing"

	"github.com/tufanbarisyildirim/gonginx/config"
)

func TestDirective_ToString(t *testing.T) {
	t.Parallel()
	type fields struct {
		Name       string
		Parameters []config.Parameter
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "server_name direction",
			fields: fields{
				Name: "server_name",
				Parameters: []config.Parameter{
					{Value: "gonginx.dev"},
					{Value: "gonginx.local"},
					{Value: "microspector.com"},
				},
			},
			want: "server_name gonginx.dev gonginx.local microspector.com;",
		},
		{
			name: "proxy_pass direction",
			fields: fields{
				Name: "proxy_pass",
				Parameters: []config.Parameter{
					{Value: "http://127.0.0.1/"},
				},
			},
			want: "proxy_pass http://127.0.0.1/;",
		},
		{
			name: "proxy_set_header direction",
			fields: fields{
				Name: "proxy_set_header",
				Parameters: []config.Parameter{
					{Value: "Host"},
					{Value: "$host"},
				},
			},
			want: "proxy_set_header Host $host;",
		},
		{
			name: "proxy_buffers direction",
			fields: fields{
				Name: "proxy_buffers",
				Parameters: []config.Parameter{
					{Value: "4"},
					{Value: "32k"},
				},
			},
			want: "proxy_buffers 4 32k;",
		},
		{
			name: "charset direction",
			fields: fields{
				Name: "charset",
				Parameters: []config.Parameter{
					{Value: "koi8-r"},
				},
			},
			want: "charset koi8-r;",
		},
		{
			name: "'' close",
			fields: fields{
				Name: "''",
				Parameters: []config.Parameter{
					{Value: "close"},
				},
			},
			want: "'' close;",
		},
	}
	for _, tt := range tests {
		tt2 := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			d := &config.Directive{
				Name:       tt2.fields.Name,
				Parameters: tt2.fields.Parameters,
			}
			if got := DumpDirective(d, NoIndentStyle); got != tt2.want {
				t.Errorf("Directive.ToString() = %v, want %v", got, tt2.want)
			}
		})
	}
}
