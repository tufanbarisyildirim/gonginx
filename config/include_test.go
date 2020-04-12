package config

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestConfig_IncludeToString(t *testing.T) {

	include := &Include{
		IncludePath: "/etc/nginx/conf.d/*.conf",
	}
	assert.Equal(t, "include /etc/nginx/conf.d/*.conf;", include.ToString())
	var i interface{} = include
	_, ok := i.(Statement)
	_, ok2 := i.(IncludeStatement)
	assert.Assert(t, ok)
	assert.Assert(t, ok2)
}

func TestInclude_SaveToFile(t *testing.T) {
	type fields struct {
		IncludePath string
		Config      *Config
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test saving file",
			fields: fields{
				IncludePath: "../full-example/makedir/*.conf",
				Config: &Config{
					FilePath: "../full-example/makedir/included.conf",
					Block: &Block{
						Statements: []Statement{
							
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				IncludePath: tt.fields.IncludePath,
				Config:      tt.fields.Config,
			}
			if err := i.SaveToFile(); (err != nil) != tt.wantErr {
				t.Errorf("Include.SaveToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
