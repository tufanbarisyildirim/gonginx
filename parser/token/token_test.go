package token

import (
	"testing"
)

func TestToken_EqualTo(t *testing.T) {
	tests := []struct {
		name   string
		Token1 Token
		Token2 Token
		want   bool
	}{
		{
			name: "server is server",
			Token1: Token{
				Type:    Keyword,
				Literal: "server",
			},
			Token2: Token{
				Type:    Keyword,
				Literal: "server",
			},
			want: true,
		},
		{
			name: "loc is not server",
			Token1: Token{
				Type:    Keyword,
				Literal: "server",
			},
			Token2: Token{
				Type:    Keyword,
				Literal: "location",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.Token1.EqualTo(tt.Token2); got != tt.want {
				t.Errorf("Token.EqualTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToken_is(t *testing.T) {
	type fields struct {
		Type    Type
		Literal string
		Line    int
		Column  int
	}
	type args struct {
		typ Type
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok := Token{
				Type:    tt.fields.Type,
				Literal: tt.fields.Literal,
				Line:    tt.fields.Line,
				Column:  tt.fields.Column,
			}
			if got := tok.Is(tt.args.typ); got != tt.want {
				t.Errorf("Token.is() = %v, want %v", got, tt.want)
			}
		})
	}
}
