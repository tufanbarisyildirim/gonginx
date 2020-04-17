package dumper

import (
	"reflect"
	"testing"
)

func TestStyle_Iterate(t *testing.T) {
	tests := []struct {
		name  string
		style *Style
		want  *Style
	}{
		{
			name:  "iteration test",
			style: NewStyle(),
			want: &Style{
				SortDirectives: false,
				StartIndent:    4,
				Indent:         4,
			},
		},
		{
			name:  "always empty no interation constant",
			style: NoIndentStyle,
			want: &Style{
				SortDirectives: false,
				StartIndent:    0,
				Indent:         0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.style.Iterate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Style.Iterate() = %v, want %v", got, tt.want)
			}
		})
	}
}
