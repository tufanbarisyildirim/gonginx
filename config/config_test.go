package config

import (
	"os"
	"testing"
)

func TestConfig_ToString(t *testing.T) {
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
			if got := c.ToString(); got != tt.want || string(c.ToByteArray()) != tt.want {
				t.Errorf("Config.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SaveToFile(t *testing.T) {
	type fields struct {
		Block    *Block
		FilePath string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "block",
			fields: fields{
				FilePath: "../full-example/unit-test.conf",
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
					},
				},
			},
			wantErr: false,
		},
		{
			name: "block",
			fields: fields{
				FilePath: "../full-example/unit-test.conf",
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
			},
			wantErr: true,
		},
		{
			name: "block",
			fields: fields{
				FilePath: "../full-example/unittest/file.conf",
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
			},
			wantErr: true,
		},
	}
	os.RemoveAll("../full-example/makedir")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Block:    tt.fields.Block,
				FilePath: tt.fields.FilePath,
			}
			if err := c.SaveToFile(); (err != nil) != tt.wantErr {
				t.Errorf("Config.SaveToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
