package parser

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/tufanbarisyildirim/gonginx"
	"github.com/tufanbarisyildirim/gonginx/parser/token"
)

// Option parsing option
type Option func(*Parser)

type options struct {
	parseInclude          bool
	skipIncludeParsingErr bool
}

// Parser is an nginx config parser
type Parser struct {
	opts              options
	configRoot        string // TODO: confirmation needed (whether this is the parent of nginx.conf)
	lexer             *lexer
	currentToken      token.Token
	followingToken    token.Token
	parsedIncludes    map[string]*gonginx.Config
	statementParsers  map[string]func() gonginx.IDirective
	blockWrappers     map[string]func(*gonginx.Directive) gonginx.IDirective
	directiveWrappers map[string]func(*gonginx.Directive) gonginx.IDirective
	commentBuffer     []string
}

// WithSameOptions copy options from another parser
func WithSameOptions(p *Parser) Option {
	return func(curr *Parser) {
		curr.opts = p.opts
	}
}

func withParsedIncludes(parsedIncludes map[string]*gonginx.Config) Option {
	return func(p *Parser) {
		p.parsedIncludes = parsedIncludes
	}
}

func withConfigRoot(configRoot string) Option {
	return func(p *Parser) {
		p.configRoot = configRoot
	}
}

// WithSkipIncludeParsingErr ignores include parsing errors
func WithSkipIncludeParsingErr() Option {
	return func(p *Parser) {
		p.opts.skipIncludeParsingErr = true
	}
}

// WithDefaultOptions default options
func WithDefaultOptions() Option {
	return func(p *Parser) {
		p.opts = options{}
	}
}

// WithIncludeParsing enable parsing included files
func WithIncludeParsing() Option {
	return func(p *Parser) {
		p.opts.parseInclude = true
	}
}

// NewStringParser parses nginx conf from string
func NewStringParser(str string, opts ...Option) *Parser {
	return NewParserFromLexer(lex(str), opts...)
}

// NewParser create new parser
func NewParser(filePath string, opts ...Option) (*Parser, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	l := newLexer(bufio.NewReader(f))
	l.file = filePath
	p := NewParserFromLexer(l, opts...)
	return p, nil
}

// NewParserFromLexer initilizes a new Parser
func NewParserFromLexer(lexer *lexer, opts ...Option) *Parser {
	configRoot, _ := filepath.Split(lexer.file)
	parser := &Parser{
		lexer:          lexer,
		opts:           options{},
		parsedIncludes: make(map[string]*gonginx.Config),
		configRoot:     configRoot,
	}

	for _, o := range opts {
		o(parser)
	}

	parser.nextToken()
	parser.nextToken()

	parser.blockWrappers = map[string]func(*gonginx.Directive) gonginx.IDirective{
		"http": func(directive *gonginx.Directive) gonginx.IDirective {
			return parser.wrapHTTP(directive)
		},
		"server": func(directive *gonginx.Directive) gonginx.IDirective {
			return parser.wrapServer(directive)
		},
		"location": func(directive *gonginx.Directive) gonginx.IDirective {
			return parser.wrapLocation(directive)
		},
		"upstream": func(directive *gonginx.Directive) gonginx.IDirective {
			return parser.wrapUpstream(directive)
		},
	}

	parser.directiveWrappers = map[string]func(*gonginx.Directive) gonginx.IDirective{
		"server": func(directive *gonginx.Directive) gonginx.IDirective {
			return parser.parseUpstreamServer(directive)
		},
		"include": func(directive *gonginx.Directive) gonginx.IDirective {
			return parser.parseInclude(directive)
		},
	}

	return parser
}

func (p *Parser) nextToken() {
	p.currentToken = p.followingToken
	p.followingToken = p.lexer.scan()
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.currentToken.Type == t
}

func (p *Parser) followingTokenIs(t token.Type) bool {
	return p.followingToken.Type == t
}

// Parse the gonginx.
func (p *Parser) Parse() *gonginx.Config {
	return &gonginx.Config{
		FilePath: p.lexer.file, //TODO: set filepath here,
		Block:    p.parseBlock(),
	}
}

// ParseBlock parse a block statement
func (p *Parser) parseBlock() *gonginx.Block {

	context := &gonginx.Block{
		Directives: make([]gonginx.IDirective, 0),
	}
	var s gonginx.IDirective
	var line int
parsingloop:
	for {
		switch {
		case p.curTokenIs(token.EOF) || p.curTokenIs(token.BlockEnd):
			break parsingloop
		case p.curTokenIs(token.Keyword) || p.curTokenIs(token.QuotedString):
			s = p.parseStatement()
			context.Directives = append(context.Directives, s)
			line = p.currentToken.Line
		case p.curTokenIs(token.Comment):
			p.commentBuffer = append(p.commentBuffer, p.currentToken.Literal)
			// inline comment
			if line == p.currentToken.Line {
				s.SetComment(p.commentBuffer)
				p.commentBuffer = nil
			}
		}
		p.nextToken()
	}

	return context
}

func (p *Parser) parseStatement() gonginx.IDirective {
	d := &gonginx.Directive{
		Name: p.currentToken.Literal,
	}

	//if we have a special parser for the directive, we use it.
	if sp, ok := p.statementParsers[d.Name]; ok {
		return sp()
	}

	if len(p.commentBuffer) > 0 {
		d.Comment = p.commentBuffer
		p.commentBuffer = nil
	}

	//parse parameters until the end.
	for p.nextToken(); p.currentToken.IsParameterEligible(); p.nextToken() {
		d.Parameters = append(d.Parameters, p.currentToken.Literal)
	}

	//if we find a semicolon it is a directive, we will check directive converters
	if p.curTokenIs(token.Semicolon) {
		if dw, ok := p.directiveWrappers[d.Name]; ok {
			return dw(d)
		}
		return d
	}
	for {
		if p.curTokenIs(token.Comment) {
			p.commentBuffer = append(p.commentBuffer, p.currentToken.Literal)
			p.nextToken()
		} else {
			break
		}
	}

	//ok, it does not end with a semicolon but a block starts, we will convert that block if we have a converter
	if p.curTokenIs(token.BlockStart) {
		d.Block = p.parseBlock()
		if bw, ok := p.blockWrappers[d.Name]; ok {
			return bw(d)
		}
		return d
	}

	return &gonginx.Directive{}
	//panic(fmt.Errorf("unexpected token %s (%s) on line %d, column %d", p.currentToken.Type.String(), p.currentToken.Literal, p.currentToken.Line, p.currentToken.Column))
}

// TODO: move this into gonginx.Include
func (p *Parser) parseInclude(directive *gonginx.Directive) *gonginx.Include {
	include := &gonginx.Include{
		Directive:   directive,
		IncludePath: directive.Parameters[0],
	}

	if len(directive.Parameters) > 1 {
		panic("include directive can not have multiple parameters")
	}

	if directive.Block != nil {
		panic("include can not have a block, or missing semicolon at the end of include statement")
	}

	if p.opts.parseInclude {
		includePath := include.IncludePath
		if !filepath.IsAbs(includePath) {
			includePath = filepath.Join(p.configRoot, include.IncludePath)
		}
		includePaths, err := filepath.Glob(includePath)
		if err != nil && !p.opts.skipIncludeParsingErr {
			panic(err)
		}
		for _, includePath := range includePaths {
			if conf, ok := p.parsedIncludes[includePath]; ok {
				// same file includes itself? don't blow up the parser
				if conf == nil {
					continue
				}
			} else {
				p.parsedIncludes[includePath] = nil
			}

			parser, err := NewParser(includePath,
				WithSameOptions(p),
				withParsedIncludes(p.parsedIncludes),
				withConfigRoot(p.configRoot),
			)

			if err != nil && !p.opts.skipIncludeParsingErr {
				panic(err)
			}

			config := parser.Parse()
			p.parsedIncludes[includePath] = config
			include.Configs = append(include.Configs, config)
		}
	}

	return include
}

// TODO: move this into gonginx.Location
func (p *Parser) wrapLocation(directive *gonginx.Directive) *gonginx.Location {
	location := &gonginx.Location{
		Modifier:  "",
		Match:     "",
		Directive: directive,
	}

	if len(directive.Parameters) == 0 {
		panic("no enough parameter for location")
	}

	if len(directive.Parameters) == 1 {
		location.Match = directive.Parameters[0]
		return location
	} else if len(directive.Parameters) == 2 {
		location.Modifier = directive.Parameters[0]
		location.Match = directive.Parameters[1]
		return location
	}

	panic("too many arguments for location directive")
}

func (p *Parser) wrapServer(directive *gonginx.Directive) *gonginx.Server {
	s, _ := gonginx.NewServer(directive)
	return s
}

func (p *Parser) wrapUpstream(directive *gonginx.Directive) *gonginx.Upstream {
	s, _ := gonginx.NewUpstream(directive)
	return s
}

func (p *Parser) wrapHTTP(directive *gonginx.Directive) *gonginx.HTTP {
	h, _ := gonginx.NewHTTP(directive)
	return h
}

func (p *Parser) parseUpstreamServer(directive *gonginx.Directive) *gonginx.UpstreamServer {
	return gonginx.NewUpstreamServer(directive)
}
