package gonginx

import (
	"reflect"
	"testing"

)

func TestBlock_ToString(t *testing.T) {
	type fields struct {
		Directives []IDirective
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
				Directives: make([]IDirective, 0),
			},
			want: "",
		},
		{
			name: "statement list",
			fields: fields{
				Directives: []IDirective{
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
			want:                        "user nginx nginx;\nworker_processes 5;",
			wantSorted:                  "user nginx nginx;\nworker_processes 5;",
			wantSortedSpaceBeforeBlocks: "user nginx nginx;\nworker_processes 5;",
		},
		{
			name: "statement list with wrapped directives",
			fields: fields{
				Directives: []IDirective{
					&Directive{
						Name:       "user",
						Parameters: []string{"nginx", "nginx"},
					},
					&Directive{
						Name:       "worker_processes",
						Parameters: []string{"5"},
					},
					&Include{
						Directive: &Directive{
							Name:       "include",
							Parameters: []string{"/etc/nginx/conf/*.conf"},
						},
						IncludePath: "/etc/nginx/conf/*.conf",
					},
					NewServerOrNill(&Directive{
						Block: &Block{
							Directives: []IDirective{
								&Directive{
									Name:       "user",
									Parameters: []string{"nginx", "nginx"},
								},
								&Directive{
									Name:       "worker_processes",
									Parameters: []string{"5"},
								},
								&Include{
									Directive: &Directive{
										Name:       "include",
										Parameters: []string{"/etc/nginx/conf/*.conf"},
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
			b := &Block{
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

func NewServerOrNill(directive IDirective) *Server {
	s, _ := NewServer(directive)
	return s
}

func TestBlock_FindDirectives(t *testing.T) {
	type args struct {
		directiveName string
	}
	tests := []struct {
		name  string
		block *Block
		args  args
		want  []IDirective
	}{
		{
			name: "find all servers",
			block: &Block{
				Directives: []IDirective{
					&Server{
						Block: &Block{
							Directives: []IDirective{
								&Directive{
									Name:       "server_name",
									Parameters: []string{"gonginx.dev"},
								},
							},
						},
					},
					&Server{
						Block: &Block{
							Directives: []IDirective{
								&Directive{
									Name:       "server_name",
									Parameters: []string{"gonginx2.dev"},
								},
							},
						},
					},
					&Http{
						&Block{
							Directives: []IDirective{
								&Server{
									Block: &Block{
										Directives: []IDirective{
											&Directive{
												Name:       "server_name",
												Parameters: []string{"gonginx3.dev"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: []IDirective{
				&Server{
					Block: &Block{
						Directives: []IDirective{
							&Directive{
								Name:       "server_name",
								Parameters: []string{"gonginx.dev"},
							},
						},
					},
				},
				&Server{
					Block: &Block{
						Directives: []IDirective{
							&Directive{
								Name:       "server_name",
								Parameters: []string{"gonginx2.dev"},
							},
						},
					},
				},
				&Server{
					Block: &Block{
						Directives: []IDirective{
							&Directive{
								Name:       "server_name",
								Parameters: []string{"gonginx3.dev"},
							},
						},
					},
				},
			},
			args: args{
				directiveName: "server",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.block.FindDirectives(tt.args.directiveName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Block.FindDirectives() = %v want %v", got, tt.want)
			}
		})
	}
}
