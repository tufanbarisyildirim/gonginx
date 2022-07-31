package gonginx

import (
	"testing"
)

func TestConfig_ToString(t *testing.T) {
	t.Parallel()
	type fields struct {
		Block    *Block
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
			},
			want: "user nginx nginx;\nworker_processes 5;\ninclude /etc/nginx/conf/*.conf;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
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

//TODO(tufan): reactive after getting SaveToFile() done
//func TestConfig_SaveToFile(t *testing.T) {
//	type fields struct {
//		Block    *Block
//		FilePath string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		wantErr bool
//	}{
//		{
//			name: "block",
//			fields: fields{
//				FilePath: "../full-example/unit-test.conf",
//				Block: &Block{
//					Directives: []IDirective{
//						&Directive{
//							Name:       "user",
//							Parameters: []string{"nginx", "nginx"},
//						},
//						&Directive{
//							Name:       "worker_processes",
//							Parameters: []string{"5"},
//						},
//					},
//				},
//			},
//			wantErr: false,
//		},
//		{
//			name: "block",
//			fields: fields{
//				FilePath: "../full-example/unit-test.conf",
//				Block: &Block{
//					Directives: []IDirective{
//						&Directive{
//							Name:       "user",
//							Parameters: []string{"nginx", "nginx"},
//						},
//						&Directive{
//							Name:       "worker_processes",
//							Parameters: []string{"5"},
//						},
//						&Include{
//							Directive: &Directive{
//								Name:       "include",
//								Parameters: []string{"/etc/nginx/conf/*.conf"},
//							},
//							IncludePath: "/etc/nginx/conf/*.conf",
//						},
//					},
//				},
//			},
//			wantErr: true,
//		},
//		{
//			name: "block",
//			fields: fields{
//				FilePath: "../full-example/unittest/file.conf",
//				Block: &Block{
//					Directives: []IDirective{
//						&Directive{
//							Name:       "user",
//							Parameters: []string{"nginx", "nginx"},
//						},
//						&Directive{
//							Name:       "worker_processes",
//							Parameters: []string{"5"},
//						},
//						&Include{
//							Directive: &Directive{
//								Name:       "include",
//								Parameters: []string{"/etc/nginx/conf/*.conf"},
//							},
//							IncludePath: "/etc/nginx/conf/*.conf",
//						},
//					},
//				},
//			},
//			wantErr: true,
//		},
//	}
//	os.RemoveAll("../full-example/makedir")
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &Config{
//				Block:    tt.fields.Block,
//				FilePath: tt.fields.FilePath,
//			}
//			if err := c.SaveToFile(NoIndentStyle); (err != nil) != tt.wantErr {
//				t.Errorf("SaveToFile() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
