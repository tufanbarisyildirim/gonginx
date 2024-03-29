package gonginx

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestConfig_IncludeToString(t *testing.T) {
	t.Parallel()
	include := &Include{
		Directive: &Directive{
			Name:       "include",
			Parameters: []string{"/etc/nginx/conf.d/*.conf"},
		},
		IncludePath: "/etc/nginx/conf.d/*.conf",
	}
	assert.Equal(t, "include /etc/nginx/conf.d/*.conf;", DumpDirective(include, NoIndentStyle))
	var i interface{} = include
	_, ok := i.(IDirective)
	//_, ok2 := i.(IncludeDirective)// TODO(tufan):reactivate here after getting include and file things done
	assert.Assert(t, ok)
	//assert.Assert(t, ok2)
}

//TODO(tufan): reactivate here after getting SaveToFile() done
//func TestInclude_SaveToFile(t *testing.T) {
//	type fields struct {
//		IncludePath string
//		Config      *Config
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		wantErr bool
//	}{
//		{
//			name: "test saving file",
//			fields: fields{
//				IncludePath: "../full-example/unittest/*.conf",
//				Config: &Config{
//					FilePath: "../full-example/unittest/included.conf",
//					Block: &Block{
//						Directives: []IDirective{},
//					},
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			i := &Include{
//				IncludePath: tt.fields.IncludePath,
//				Configs: []*Config{
//					tt.fields.Config,
//				},
//			}
//			if err := i.SaveToFile(NoIndentStyle); (err != nil) != tt.wantErr {
//				t.Errorf("Include.SaveToFile() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
