package dumper

import (
	"testing"

	"github.com/tufanbarisyildirim/gonginx/config"
)

func TestConfig_ToString(t *testing.T) {
	t.Parallel()
	type fields struct {
		Block    *config.Block
		FilePath string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "block",
			fields: fields{
				Block: &config.Block{
					Directives: []config.IDirective{
						&config.Directive{
							Name: "user",
							Parameters: []config.Parameter{
								{Value: "nginx"},
								{Value: "nginx"},
							},
						},
						&config.Directive{
							Name:       "worker_processes",
							Parameters: []config.Parameter{{Value: "5"}},
						},
						&config.Include{
							Directive: &config.Directive{
								Name:       "include",
								Parameters: []config.Parameter{{Value: "/etc/nginx/conf/*.conf"}},
							},
							IncludePath: "/etc/nginx/conf/*.conf",
						},
					},
				},
			},
			want: "user nginx nginx;\nworker_processes 5;\ninclude /etc/nginx/conf/*.conf;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &config.Config{
				Block:    tt.fields.Block,
				FilePath: tt.fields.FilePath,
			}
			//TODO(tufan): create another dumper for a config and include statement (file thingis)
			if got := DumpConfig(c, NoIndentStyle); got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
