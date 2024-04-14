package dumper

import (
	"reflect"
	"testing"

	"github.com/tufanbarisyildirim/gonginx/config"
	"gotest.tools/v3/assert"
)

func TestNewUpstreamServer(t *testing.T) {
	t.Parallel()
	type args struct {
		directive *config.Directive
	}
	tests := []struct {
		name       string
		args       args
		want       *config.UpstreamServer
		wantString string
	}{
		{
			name: "new upstream server",
			args: args{
				directive: &config.Directive{
					Name:       "server",
					Parameters: []string{"127.0.0.1:8080"},
				},
			},
			want: &config.UpstreamServer{
				Address:    "127.0.0.1:8080",
				Flags:      make([]string, 0),
				Parameters: make(map[string]string, 0),
			},
			wantString: "server 127.0.0.1:8080;",
		},
		{
			name: "new upstream server with weight",
			args: args{
				directive: &config.Directive{
					Name:       "server",
					Parameters: []string{"127.0.0.1:8080", "weight=5"},
				},
			},
			want: &config.UpstreamServer{
				Address: "127.0.0.1:8080",
				Flags:   make([]string, 0),
				Parameters: map[string]string{
					"weight": "5",
				},
			},
			wantString: "server 127.0.0.1:8080 weight=5;",
		},
		{
			name: "new upstream server with weight and a flag",
			args: args{
				directive: &config.Directive{
					Name:       "server",
					Parameters: []string{"127.0.0.1:8080", "weight=5", "down"},
				},
			},
			want: &config.UpstreamServer{
				Address: "127.0.0.1:8080",
				Flags:   []string{"down"},
				Parameters: map[string]string{
					"weight": "5",
				},
			},
			wantString: "server 127.0.0.1:8080 weight=5 down;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := config.NewUpstreamServer(tt.args.directive)
			assert.NilError(t, err, "no error expected here")
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUpstreamServer() = %v, want %v", got, tt.want)
			}

			if got.GetBlock() != nil {
				t.Error("Upstream server returns a block")
			}

			gotString := DumpDirective(got, NoIndentStyle)
			if gotString != tt.wantString {
				t.Errorf("NewUpstreamServer().ToString = %v, want %v", gotString, tt.wantString)
			}
		})
	}
}
