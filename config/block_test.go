package config

import (
	"reflect"
	"testing"
)

func NewServerOrNill(directive IDirective) *Server {
	s, _ := NewServer(directive)
	return s
}

func TestBlock_FindDirectives(t *testing.T) {
	t.Parallel()
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
									Parameters: []Parameter{{Value: "gonginx.dev"}},
								},
							},
						},
					},
					&Server{
						Block: &Block{
							Directives: []IDirective{
								&Directive{
									Name:       "server_name",
									Parameters: []Parameter{{Value: "gonginx2.dev"}},
								},
							},
						},
					},
					&HTTP{
						Servers: []*Server{
							{
								Block: &Block{
									Directives: []IDirective{
										&Directive{
											Name:       "server_name",
											Parameters: []Parameter{{Value: "gonginx3.dev"}},
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
								Parameters: []Parameter{{Value: "gonginx.dev"}},
							},
						},
					},
				},
				&Server{
					Block: &Block{
						Directives: []IDirective{
							&Directive{
								Name:       "server_name",
								Parameters: []Parameter{{Value: "gonginx2.dev"}},
							},
						},
					},
				},
				&Server{
					Block: &Block{
						Directives: []IDirective{
							&Directive{
								Name:       "server_name",
								Parameters: []Parameter{{Value: "gonginx3.dev"}},
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
