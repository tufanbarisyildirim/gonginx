package token

import (
	"fmt"
	"reflect"
	"testing"
)

func TestType_String(t *testing.T) {
	tests := []struct {
		name string
		tt   Type
		want string
	}{
		{
			name: "QuotedString",
			tt:   QuotedString,
			want: "QuotedString",
		},
		{
			name: "Eof",
			tt:   EOF,
			want: "Eof",
		},
		{
			name: "Keyword",
			tt:   Keyword,
			want: "Keyword",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tt.String(); got != tt.want {
				t.Errorf("Type.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToken_String(t *testing.T) {
	type fields struct {
		Type    Type
		Literal string
		Line    int
		Column  int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "serialize quoted string",
			fields: fields{
				Type:    QuotedString,
				Literal: "my test string",
				Line:    0,
				Column:  0,
			},
			want: fmt.Sprintf("{Type:%s,Literal:\"%s\",Line:%d,Column:%d}", "QuotedString", "my test string", 0, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok := Token{
				Type:    tt.fields.Type,
				Literal: tt.fields.Literal,
				Line:    tt.fields.Line,
				Column:  tt.fields.Column,
			}
			if got := tok.String(); got != tt.want {
				t.Errorf("Token.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToken_Lit(t *testing.T) {
	type fields struct {
		Type    Type
		Literal string
		Line    int
		Column  int
	}
	type args struct {
		literal string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Token
	}{
		{
			name: "a string literal",
			fields: fields{
				Type:    QuotedString,
				Literal: "a test string",
				Line:    0,
				Column:  0,
			},
			args: args{
				literal: "new test string",
			},
			want: Token{
				Type:    QuotedString,
				Literal: "new test string",
				Line:    0,
				Column:  0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok := Token{
				Type:    tt.fields.Type,
				Literal: tt.fields.Literal,
				Line:    tt.fields.Line,
				Column:  tt.fields.Column,
			}
			if got := tok.Lit(tt.args.literal); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Token.Lit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToken_EqualTo(t *testing.T) {

	tests := []struct {
		name string
		tok1 Token
		tok2 Token
		want bool
	}{
		{
			name: "keyword is keyword",
			tok1: Token{
				Type:    Keyword,
				Literal: "server",
			},
			tok2: Token{
				Type:    Keyword,
				Literal: "server",
			},
			want: true,
		},
		{
			name: "keyword is keyword but needs same directive",
			tok1: Token{
				Type:    Keyword,
				Literal: "server",
			},
			tok2: Token{
				Type:    Keyword,
				Literal: "location",
			},
			want: false,
		}, {
			name: "string is string",
			tok1: Token{
				Type:    QuotedString,
				Literal: "same quoted strings",
			},
			tok2: Token{
				Type:    QuotedString,
				Literal: "same quoted strings",
			},
			want: true,
		},
		{
			name: "Blockstart is Blockstart even if they are in different lines",
			tok1: Token{
				Type:    BlockStart,
				Literal: "{",
				Line:    1,
			},
			tok2: Token{
				Type:    BlockStart,
				Literal: "{",
				Line:    2,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tok1.EqualTo(tt.tok2); got != tt.want {
				t.Errorf("Token.EqualTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokens_EqualTo(t *testing.T) {
	type args struct {
		tokens Tokens
	}
	tests := []struct {
		name string
		ts   Tokens
		args args
		want bool
	}{
		{
			name: "token array matching",
			ts: Tokens{
				{Type: Keyword, Literal: "server", Line: 2, Column: 1},
				{Type: BlockStart, Literal: "{", Line: 2, Column: 8},
				{Type: Comment, Literal: "# simple reverse-proxy", Line: 2, Column: 10},
				{Type: BlockEnd, Literal: "}", Line: 3, Column: 5},
			},
			args: args{
				tokens: Tokens{
					{Type: Keyword, Literal: "server", Line: 2, Column: 1},
					{Type: BlockStart, Literal: "{", Line: 2, Column: 8},
					{Type: Comment, Literal: "# simple reverse-proxy", Line: 2, Column: 10},
					{Type: BlockEnd, Literal: "}", Line: 3, Column: 5},
				},
			},
			want: true,
		},
		{
			name: "token array matching",
			ts: Tokens{
				{Type: Keyword, Literal: "server", Line: 2, Column: 1},
				{Type: BlockStart, Literal: "{", Line: 2, Column: 8},
				{Type: Comment, Literal: "# simple reverse-proxy", Line: 2, Column: 10},
				{Type: BlockEnd, Literal: "}", Line: 3, Column: 5},
			},
			args: args{
				tokens: Tokens{
					{Type: Keyword, Literal: "server", Line: 2, Column: 1},
					{Type: BlockStart, Literal: "{", Line: 2, Column: 8},
					{Type: Comment, Literal: "# simple reverse-proxy", Line: 2, Column: 10},
				},
			},
			want: false,
		},
		{
			name: "token array matching",
			ts: Tokens{
				{Type: Keyword, Literal: "server", Line: 2, Column: 1},
				{Type: BlockStart, Literal: "{", Line: 2, Column: 8},
				{Type: Comment, Literal: "# simple reverse-proxy", Line: 2, Column: 10},
				{Type: BlockEnd, Literal: "}", Line: 3, Column: 5},
			},
			args: args{
				tokens: Tokens{
					{Type: Keyword, Literal: "server", Line: 2, Column: 1},
					{Type: BlockStart, Literal: "{", Line: 2, Column: 8},
					{Type: QuotedString, Literal: "simple reverse-proxy", Line: 2, Column: 10},
					{Type: BlockEnd, Literal: "}", Line: 3, Column: 5},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ts.EqualTo(tt.args.tokens); got != tt.want {
				t.Errorf("Tokens.EqualTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToken_Is(t *testing.T) {
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
		{
			name: "QuotedString Type",
			fields: fields{
				Type:    QuotedString,
				Literal: "hello",
			},
			args: args{
				typ: QuotedString,
			},
			want: true,
		},
		{
			name: "QuotedString Type",
			fields: fields{
				Type:    QuotedString,
				Literal: "hello",
			},
			args: args{
				typ: Keyword,
			},
			want: false,
		},
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
				t.Errorf("Token.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToken_IsParameterEligible(t *testing.T) {

	tests := []struct {
		name  string
		token Token
		want  bool
	}{
		{
			name: "Keyword can be a parameter",
			token: Token{
				Type: Keyword,
			},
			want: true,
		},
		{
			name: "Variable can be a parameter",
			token: Token{
				Type: Variable,
			},
			want: true,
		},
		{
			name: "Quoted string can be a parameter",
			token: Token{
				Type: QuotedString,
			},
			want: true,
		},
		{
			name: "Quoted string can be a parameter",
			token: Token{
				Type: QuotedString,
			},
			want: true,
		},
		{
			name: "Regex string can be a parameter",
			token: Token{
				Type: Regex,
			},
			want: true,
		},
		{
			name: "Blockstart cant can be a parameter",
			token: Token{
				Type: BlockStart,
			},
			want: false,
		},
		{
			name: "Blockend cant can be a parameter",
			token: Token{
				Type: BlockEnd,
			},
			want: false,
		},
		{
			name: "Comment cant can be a parameter",
			token: Token{
				Type: Comment,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.token.IsParameterEligible(); got != tt.want {
				t.Errorf("Token.IsParameterEligible() = %v, want %v", got, tt.want)
			}
		})
	}
}
