package dumper

import (
	"testing"

	"github.com/tufanbarisyildirim/gonginx/config"
)

func NewServerOrNill(directive config.IDirective) *config.Server {
	s, _ := config.NewServer(directive)
	return s
}
func TestBlock_ToString(t *testing.T) {
	t.Parallel()
	type fields struct {
		Directives []config.IDirective
	}
	tests := []struct {
		name                        string
		fields                      fields
		want                        string
		wantSorted                  string
		wantSortedSpaceBeforeBlocks string
	}{
		{
			name: "empty block",
			fields: fields{
				Directives: make([]config.IDirective, 0),
			},
			want: "",
		},
		{
			name: "statement list",
			fields: fields{
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
				},
			},
			want:                        "user nginx nginx;\nworker_processes 5;",
			wantSorted:                  "user nginx nginx;\nworker_processes 5;",
			wantSortedSpaceBeforeBlocks: "user nginx nginx;\nworker_processes 5;",
		},
		{
			name: "statement list with wrapped directives",
			fields: fields{
				Directives: []config.IDirective{
					&config.Directive{
						Name: "user",
						Parameters: []config.Parameter{
							{Value: "nginx"},
							{Value: "nginx"}},
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
					NewServerOrNill(&config.Directive{
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
						Name: "server",
					}),
				},
			},
			want:                        "user nginx nginx;\nworker_processes 5;\ninclude /etc/nginx/conf/*.conf;\nserver {\nuser nginx nginx;\nworker_processes 5;\ninclude /etc/nginx/conf/*.conf;\n}",
			wantSorted:                  "include /etc/nginx/conf/*.conf;\nserver {\ninclude /etc/nginx/conf/*.conf;\nuser nginx nginx;\nworker_processes 5;\n}\nuser nginx nginx;\nworker_processes 5;",
			wantSortedSpaceBeforeBlocks: "include /etc/nginx/conf/*.conf;\n\nserver {\ninclude /etc/nginx/conf/*.conf;\nuser nginx nginx;\nworker_processes 5;\n}\nuser nginx nginx;\nworker_processes 5;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &config.Block{
				Directives: tt.fields.Directives,
			}
			if got := DumpBlock(b, NoIndentStyle); got != tt.want {
				t.Errorf("Block.ToString(NoIndentStyle) = \"%v\", want \"%v\"", got, tt.want)
			}
			if got := DumpBlock(b, NoIndentSortedStyle); got != tt.wantSorted {
				t.Errorf("Block.ToString(NoIndentSortedStyle) = \"%v\", want \"%v\"", got, tt.wantSorted)
			}
			if got := DumpBlock(b, NoIndentSortedSpaceStyle); got != tt.wantSortedSpaceBeforeBlocks {
				t.Errorf("Block.ToString(NoIndentSortedSpaceStyle) = \"%v\", want \"%v\"", got, tt.wantSortedSpaceBeforeBlocks)
			}
		})
	}
}
