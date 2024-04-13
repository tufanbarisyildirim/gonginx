package parser

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tufanbarisyildirim/gonginx"
	"github.com/tufanbarisyildirim/gonginx/parser/token"
)

// Option parsing option
type Option func(*Parser)

type options struct {
	parseInclude               bool
	skipIncludeParsingErr      bool
	skipComments               bool
	customDirectives           map[string]string
	skipValidSubDirectiveBlock map[string]struct{}
	skipValidDirectivesErr     bool
}

// Parser is an nginx config parser
type Parser struct {
	opts              options
	configRoot        string // TODO: confirmation needed (whether this is the parent of nginx.conf)
	lexer             *lexer
	currentToken      token.Token
	followingToken    token.Token
	parsedIncludes    map[string]*gonginx.Config
	statementParsers  map[string]func() (gonginx.IDirective, error)
	blockWrappers     map[string]func(*gonginx.Directive) (gonginx.IDirective, error)
	directiveWrappers map[string]func(*gonginx.Directive) (gonginx.IDirective, error)
	commentBuffer     []string
	file              *os.File
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

// WithSkipComments default options
func WithSkipComments() Option {
	return func(p *Parser) {
		p.opts.skipComments = true
	}
}

// WithIncludeParsing enable parsing included files
func WithIncludeParsing() Option {
	return func(p *Parser) {
		p.opts.parseInclude = true
	}
}

// WithCustomDirectives add your custom directives as valid directives
func WithCustomDirectives(directives ...string) Option {
	return func(p *Parser) {
		for _, directive := range directives {
			p.opts.customDirectives[directive] = directive
		}
	}
}

// WithSkipValidBlocks add your custom block as valid
func WithSkipValidBlocks(directives ...string) Option {
	return func(p *Parser) {
		for _, directive := range directives {
			p.opts.skipValidSubDirectiveBlock[directive] = struct{}{}
		}
	}
}

// WithSkipValidDirectivesErr ignores unknown directive errors
func WithSkipValidDirectivesErr() Option {
	return func(p *Parser) {
		p.opts.skipValidDirectivesErr = true
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
	p.file = f
	return p, nil
}

// NewParserFromLexer initilizes a new Parser
func NewParserFromLexer(lexer *lexer, opts ...Option) *Parser {
	configRoot, _ := filepath.Split(lexer.file)
	parser := &Parser{
		lexer:          lexer,
		opts:           options{customDirectives: make(map[string]string), skipValidSubDirectiveBlock: map[string]struct{}{}},
		parsedIncludes: make(map[string]*gonginx.Config),
		configRoot:     configRoot,
	}

	for _, o := range opts {
		o(parser)
	}

	parser.nextToken()
	parser.nextToken()

	parser.blockWrappers = map[string]func(*gonginx.Directive) (gonginx.IDirective, error){
		"http": func(directive *gonginx.Directive) (gonginx.IDirective, error) {
			return parser.wrapHTTP(directive)
		},
		"server": func(directive *gonginx.Directive) (gonginx.IDirective, error) {
			return parser.wrapServer(directive)
		},
		"location": func(directive *gonginx.Directive) (gonginx.IDirective, error) {
			return parser.wrapLocation(directive)
		},
		"upstream": func(directive *gonginx.Directive) (gonginx.IDirective, error) {
			return parser.wrapUpstream(directive)
		},
		"_by_lua_block": func(directive *gonginx.Directive) (gonginx.IDirective, error) {
			return parser.wrapLuaBlock(directive)
		},
	}

	parser.directiveWrappers = map[string]func(*gonginx.Directive) (gonginx.IDirective, error){
		"server": func(directive *gonginx.Directive) (gonginx.IDirective, error) {
			return parser.parseUpstreamServer(directive)
		},
		"include": func(directive *gonginx.Directive) (gonginx.IDirective, error) {
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
func (p *Parser) Parse() (*gonginx.Config, error) {
	parsedBlock, err := p.parseBlock(false, false)
	if err != nil {
		return nil, err
	}
	c := &gonginx.Config{
		FilePath: p.lexer.file, //TODO: set filepath here,
		Block:    parsedBlock,
	}
	err = p.Close()
	return c, err
}

// ParseBlock parse a block statement
func (p *Parser) parseBlock(inBlock bool, isSkipValidDirective bool) (*gonginx.Block, error) {

	context := &gonginx.Block{
		Directives: make([]gonginx.IDirective, 0),
	}
	var s gonginx.IDirective
	var line int
parsingLoop:
	for {
		switch {
		case p.curTokenIs(token.EOF):
			if inBlock {
				return nil, errors.New("unexpected eof in block")
			}
			break parsingLoop
		case p.curTokenIs(token.LuaCode):
			context.IsLuaBlock = true
			context.LiteralCode = p.currentToken.Literal
		case p.curTokenIs(token.BlockEnd):
			break parsingLoop
		case p.curTokenIs(token.Keyword) || p.curTokenIs(token.QuotedString):
			s, err := p.parseStatement(isSkipValidDirective)
			if err != nil {
				return nil, err
			}
			s.SetParent(context)
			context.Directives = append(context.Directives, s)
			line = p.currentToken.Line
		case p.curTokenIs(token.Comment):
			if p.opts.skipComments {
				break
			}
			p.commentBuffer = append(p.commentBuffer, p.currentToken.Literal)
			// inline comment
			if line == p.currentToken.Line {
				if s == nil && len(context.Directives) > 0 {
					s = context.Directives[len(context.Directives)-1]
				}
				s.SetComment(p.commentBuffer)
				p.commentBuffer = nil
			}

		}
		p.nextToken()
	}

	return context, nil
}

func (p *Parser) parseStatement(isSkipValidDirective bool) (gonginx.IDirective, error) {
	d := &gonginx.Directive{
		Name: p.currentToken.Literal,
	}

	if !p.opts.skipValidDirectivesErr && !isSkipValidDirective {
		_, ok := ValidDirectives[d.Name]
		_, ok2 := p.opts.customDirectives[d.Name]

		if !ok && !ok2 {
			return nil, fmt.Errorf("unknown directive '%s' on line %d, column %d", d.Name, p.currentToken.Line, p.currentToken.Column)
		}
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
		return d, nil
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
		_, blockSkip1 := SkipValidBlocks[d.Name]
		_, blockSkip2 := p.opts.skipValidSubDirectiveBlock[d.Name]
		isSkipBlockSubDirective := blockSkip1 || blockSkip2 || isSkipValidDirective
		b, err := p.parseBlock(true, isSkipBlockSubDirective)
		if err != nil {
			return nil, err
		}
		d.Block = b

		if strings.HasSuffix(d.Name, "_by_lua_block") {
			return p.blockWrappers["_by_lua_block"](d)
		}

		if bw, ok := p.blockWrappers[d.Name]; ok {
			return bw(d)
		}
		return d, nil
	}

	return nil, fmt.Errorf("unexpected token %s (%s) on line %d, column %d", p.currentToken.Type.String(), p.currentToken.Literal, p.currentToken.Line, p.currentToken.Column)
}

// TODO: move this into gonginx.Include
func (p *Parser) parseInclude(directive *gonginx.Directive) (*gonginx.Include, error) {
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
			return nil, err
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

			config, err := parser.Parse()
			if err != nil {
				return nil, err
			}
			p.parsedIncludes[includePath] = config
			include.Configs = append(include.Configs, config)
		}
	}

	return include, nil
}

// TODO: move this into gonginx.Location
func (p *Parser) wrapLocation(directive *gonginx.Directive) (*gonginx.Location, error) {
	location := &gonginx.Location{
		Modifier:  "",
		Match:     "",
		Directive: directive,
	}

	if len(directive.Parameters) == 0 {
		return nil, errors.New("no enough parameter for location")
	}

	if len(directive.Parameters) == 1 {
		location.Match = directive.Parameters[0]
		return location, nil
	} else if len(directive.Parameters) == 2 {
		location.Modifier = directive.Parameters[0]
		location.Match = directive.Parameters[1]
		return location, nil
	}
	return nil, errors.New("too many arguments for location directive")
}

func (p *Parser) wrapServer(directive *gonginx.Directive) (*gonginx.Server, error) {
	return gonginx.NewServer(directive)
}

func (p *Parser) wrapUpstream(directive *gonginx.Directive) (*gonginx.Upstream, error) {
	return gonginx.NewUpstream(directive)
}

func (p *Parser) wrapLuaBlock(directive *gonginx.Directive) (*gonginx.LuaBlock, error) {
	return gonginx.NewLuaBlock(directive)
}

func (p *Parser) wrapHTTP(directive *gonginx.Directive) (*gonginx.HTTP, error) {
	return gonginx.NewHTTP(directive)
}

func (p *Parser) parseUpstreamServer(directive *gonginx.Directive) (*gonginx.UpstreamServer, error) {
	return gonginx.NewUpstreamServer(directive)
}

// Close closes the file handler and releases the resources
func (p *Parser) Close() (err error) {
	if p.file != nil {
		err = p.file.Close()
	}
	return err
}
