package dumper

import (
	"reflect"
	"testing"

	"path/filepath"

	"github.com/tufanbarisyildirim/gonginx/config"
	"gotest.tools/v3/assert"
)

func TestStyle_Iterate(t *testing.T) {
	t.Parallel()
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

func TestDumpBlock_SortDoesNotMutateSource(t *testing.T) {
	t.Parallel()

	b := &config.Block{
		Directives: []config.IDirective{
			&config.Directive{
				Name:       "worker_processes",
				Parameters: []config.Parameter{{Value: "1"}},
			},
			&config.Directive{
				Name: "user",
				Parameters: []config.Parameter{
					{Value: "nginx"},
					{Value: "nginx"},
				},
			},
		},
	}

	firstUnsorted := DumpBlock(b, NoIndentStyle)
	assert.Equal(t, firstUnsorted, "worker_processes 1;\nuser nginx nginx;")

	sorted := DumpBlock(b, NoIndentSortedStyle)
	assert.Equal(t, sorted, "user nginx nginx;\nworker_processes 1;")

	secondUnsorted := DumpBlock(b, NoIndentStyle)
	assert.Equal(t, secondUnsorted, firstUnsorted)
}

func TestWriteConfig_ErrorOnIncludeTypeMismatch(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	cfg := &config.Config{
		FilePath: filepath.Join(tmp, "out.conf"),
		Block: &config.Block{
			Directives: []config.IDirective{
				&config.Directive{
					Name:       "include",
					Parameters: []config.Parameter{{Value: "a.conf"}},
				},
			},
		},
	}

	err := WriteConfig(cfg, NoIndentStyle, true)
	assert.ErrorContains(t, err, "include directive type mismatch")
}
