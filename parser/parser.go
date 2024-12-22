package parser

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tufanbarisyildirim/gonginx/config"
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

func defaultOptions() options {
	return options{
		parseInclude:               false,
		skipIncludeParsingErr:      false,
		skipComments:               false,
		customDirectives:           map[string]string{},
		skipValidSubDirectiveBlock: map[string]struct{}{},
		skipValidDirectivesErr:     false,
	}
}

// Parser is an nginx config parser
type Parser struct {
	opts              options
	configRoot        string // TODO: confirmation needed (whether this is the parent of nginx.conf)
	lexer             *lexer
	currentToken      token.Token
	followingToken    token.Token
	parsedIncludes    map[*config.Include]*config.Config
	statementParsers  map[string]func() (config.IDirective, error)
	blockWrappers     map[string]func(*config.Directive) (config.IDirective, error)
	directiveWrappers map[string]func(*config.Directive) (config.IDirective, error)
	includeWrappers   map[string]func(*config.Directive) (config.IDirective, error)

	commentBuffer []string
	file          *os.File
}

// WithSameOptions copy options from another parser
func WithSameOptions(p *Parser) Option {
	return func(curr *Parser) {
		curr.opts = p.opts
	}
}

func withParsedIncludes(parsedIncludes map[*config.Include]*config.Config) Option {
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
		p.opts = defaultOptions()
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
		opts:           defaultOptions(),
		parsedIncludes: make(map[*config.Include]*config.Config),
		configRoot:     configRoot,
	}

	for _, o := range opts {
		o(parser)
	}

	parser.nextToken()
	parser.nextToken()

	parser.blockWrappers = config.BlockWrappers
	parser.directiveWrappers = config.DirectiveWrappers
	parser.includeWrappers = config.IncludeWrappers
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
func (p *Parser) Parse() (*config.Config, error) {
	parsedBlock, err := p.parseBlock(false, false)
	if err != nil {
		return nil, err
	}
	c := &config.Config{
		FilePath: p.lexer.file, //TODO: set filepath here,
		Block:    parsedBlock,
	}
	err = p.Close()
	return c, err
}

// ParseBlock parse a block statement
func (p *Parser) parseBlock(inBlock bool, isSkipValidDirective bool) (*config.Block, error) {

	context := &config.Block{
		Directives: make([]config.IDirective, 0),
	}
	var s config.IDirective
	var err error
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
			s, err = p.parseStatement(isSkipValidDirective)
			if err != nil {
				return nil, err
			}
			if s.GetBlock() == nil {
				s.SetParent(s)
			} else {
				// each directive should have a parent directive, not a block
				// find each directive in the block and set the parent directive
				b := s.GetBlock()
				for _, dir := range b.GetDirectives() {
					dir.SetParent(s)
				}
			}
			line = p.currentToken.Line
			s.SetLine(line)
			context.Directives = append(context.Directives, s)
		case p.curTokenIs(token.Comment):
			if p.opts.skipComments {
				break
			}
			// inline comment
			if line == p.currentToken.Line {
				if s == nil && len(context.Directives) > 0 {
					s = context.Directives[len(context.Directives)-1]
				}
				s.SetInlineComment(p.currentToken.Literal)
				s.SetLine(line)
				p.commentBuffer = nil
			} else {
				p.commentBuffer = append(p.commentBuffer, p.currentToken.Literal)
			}

		}
		p.nextToken()
	}

	return context, nil
}

func (p *Parser) parseStatement(isSkipValidDirective bool) (config.IDirective, error) {
	d := &config.Directive{
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
	// keep track of the line index of the directive
	directiveLineIndex := p.currentToken.Line
	for p.nextToken(); p.currentToken.IsParameterEligible(); p.nextToken() {
		d.Parameters = append(d.Parameters, config.Parameter{
			Value:             p.currentToken.Literal,
			RelativeLineIndex: p.currentToken.Line - directiveLineIndex}) // save the relative line index of the parameter
		if p.currentToken.Is(token.BlockEnd) {
			return d, nil
		}
	}

	//if we find a semicolon it is a directive, we will check directive converters
	if p.curTokenIs(token.Semicolon) {
		if iw, ok := p.includeWrappers[d.Name]; ok {
			include, err := iw(d)
			if err != nil {
				return nil, err
			}
			return p.ParseInclude(include.(*config.Include))
		} else if dw, ok := p.directiveWrappers[d.Name]; ok {
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

// ParseInclude just parse include confs
func (p *Parser) ParseInclude(include *config.Include) (config.IDirective, error) {
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
			if conf, ok := p.parsedIncludes[include]; ok {
				// same file includes itself? don't blow up the parser
				if conf == nil {
					continue
				}
			} else {
				p.parsedIncludes[include] = nil
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
			//TODO: link parent config or include direcitve?
			p.parsedIncludes[include] = config
			include.Configs = append(include.Configs, config)
		}
	}
	return include, nil
}

// Close closes the file handler and releases the resources
func (p *Parser) Close() (err error) {
	if p.file != nil {
		err = p.file.Close()
	}
	return err
}
