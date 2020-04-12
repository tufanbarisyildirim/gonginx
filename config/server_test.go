package config

import "testing"

func TestServer_ToString(t *testing.T) {
	type fields struct {
		Directive *Directive
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "empty server block",
			fields: fields{
				Directive: &Directive{
					Block: &Block{
						Statements: make([]Statement, 0),
					},
					Name: "server",
				},
			},
			want: "server {\n\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Directive: tt.fields.Directive,
			}
			if got := s.ToString(); got != tt.want {
				t.Errorf("Server.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
