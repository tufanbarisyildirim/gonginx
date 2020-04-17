package config

import (
	"testing"

	"github.com/tufanbarisyildirim/gonginx/dumper"
)

func TestDirective_ToString(t *testing.T) {
	type fields struct {
		Name       string
		Parameters []string
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
				Parameters: []string{
					"gonginx.dev",
					"gonginx.local",
					"microspector.com",
				},
			},
			want: "server_name gonginx.dev gonginx.local microspector.com;",
		},
		{
			name: "proxy_pass direction",
			fields: fields{
				Name: "proxy_pass",
				Parameters: []string{
					"http://127.0.0.1/",
				},
			},
			want: "proxy_pass http://127.0.0.1/;",
		},
		{
			name: "proxy_set_header direction",
			fields: fields{
				Name: "proxy_set_header",
				Parameters: []string{
					"Host",
					"$host",
				},
			},
			want: "proxy_set_header Host $host;",
		},
		{
			name: "proxy_buffers direction",
			fields: fields{
				Name: "proxy_buffers",
				Parameters: []string{
					"4",
					"32k",
				},
			},
			want: "proxy_buffers 4 32k;",
		},
		{
			name: "charset direction",
			fields: fields{
				Name: "charset",
				Parameters: []string{
					"koi8-r",
				},
			},
			want: "charset koi8-r;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Directive{
				Name:       tt.fields.Name,
				Parameters: tt.fields.Parameters,
			}
			if got := d.ToString(dumper.NoIndentStyle); got != tt.want {
				t.Errorf("Directive.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
