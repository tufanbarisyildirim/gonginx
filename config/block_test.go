package config

import "testing"

func TestBlock_ToString(t *testing.T) {
	type fields struct {
		Statements []Statement
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "empty block",
			fields: fields{
				Statements: make([]Statement, 0),
			},
			want: "",
		},
		{
			name: "statement list",
			fields: fields{
				Statements: []Statement{
					&Directive{
						Name:       "user",
						Parameters: []string{"nginx", "nginx"},
					},
					&Directive{
						Name:       "worker_processes",
						Parameters: []string{"5"},
					},
				},
			},
			want: "user nginx nginx;\nworker_processes 5;",
		},
		{
			name: "statement list with wrapped directives",
			fields: fields{
				Statements: []Statement{
					&Directive{
						Name:       "user",
						Parameters: []string{"nginx", "nginx"},
					},
					&Directive{
						Name:       "worker_processes",
						Parameters: []string{"5"},
					},
					&Include{
						IncludePath: "/etc/nginx/conf/*.conf",
					},
					&Server{
						Directive: &Directive{
							Block: &Block{
								Statements: []Statement{
									&Directive{
										Name:       "user",
										Parameters: []string{"nginx", "nginx"},
									},
									&Directive{
										Name:       "worker_processes",
										Parameters: []string{"5"},
									},
									&Include{
										IncludePath: "/etc/nginx/conf/*.conf",
									},
								},
							},
							Name: "server",
						},
					},
				},
			},
			want: "user nginx nginx;\nworker_processes 5;\ninclude /etc/nginx/conf/*.conf;\nserver {\nuser nginx nginx;\nworker_processes 5;\ninclude /etc/nginx/conf/*.conf;\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Block{
				Statements: tt.fields.Statements,
			}
			if got := b.ToString(); got != tt.want {
				t.Errorf("Block.ToString() = \"%v\", want \"%v\"", got, tt.want)
			}
		})
	}
}
