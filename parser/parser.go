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
	includeCycleErr            bool
	skipComments               bool
	customDirectives           map[string]string
	skipValidSubDirectiveBlock map[string]struct{}
	skipValidDirectivesErr     bool
}

func defaultOptions() options {
	return options{
		parseInclude:               false,
		skipIncludeParsingErr:      false,
		includeCycleErr:            false,
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
	parsedIncludes    map[string]*config.Config
	includeStack      map[string]struct{}
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

func withParsedIncludes(parsedIncludes map[string]*config.Config) Option {
	return func(p *Parser) {
		p.parsedIncludes = parsedIncludes
	}
}

func withIncludeStack(includeStack map[string]struct{}) Option {
	return func(p *Parser) {
		p.includeStack = includeStack
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

// WithIncludeCycleErr returns an error when include cycle is detected.
// Default behavior skips cyclic include branches.
func WithIncludeCycleErr() Option {
	return func(p *Parser) {
		p.opts.includeCycleErr = true
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
		parsedIncludes: make(map[string]*config.Config),
		includeStack:   make(map[string]struct{}),
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
func (p *Parser) Parse() (_ *config.Config, err error) {
	if p.file != nil {
		defer func() {
			closeErr := p.Close()
			if closeErr == nil {
				return
			}

			if err != nil {
				err = errors.Join(err, closeErr)
				return
			}
			err = closeErr
		}()
	}

	parsedBlock, err := p.parseBlock(false, false)
	if err != nil {
		if p.lexer.Err != nil {
			return nil, errors.Join(p.lexer.Err, err)
		}
		return nil, err
	}
	if p.lexer.Err != nil {
		return nil, p.lexer.Err
	}

	c := &config.Config{
		FilePath: p.lexer.file, //TODO: set filepath here,
		Block:    parsedBlock,
	}
	return c, nil
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
				// Root-level leaf directives have no parent.
				// Nested leaf directives get parent assignment when their containing block wrapper is processed.
				s.SetParent(nil)
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
			// outline comment
			p.commentBuffer = append(p.commentBuffer, p.currentToken.Literal)
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

	// set outline comment
	if len(p.commentBuffer) > 0 {
		d.Comment = p.commentBuffer
		p.commentBuffer = make([]string, 0)
	}

	directiveLineIndex := p.currentToken.Line // keep track of the line index of the directive
	// Parse parameters until reaching the semicolon that ends the directive.
	for {
		p.nextToken()
		if p.currentToken.IsParameterEligible() {
			d.Parameters = append(d.Parameters, config.Parameter{
				Value:             p.currentToken.Literal,
				RelativeLineIndex: p.currentToken.Line - directiveLineIndex}) // save the relative line index of the parameter
			if p.currentToken.Is(token.BlockEnd) {
				return d, nil
			}
		} else if p.curTokenIs(token.Semicolon) {
			// inline comment in following token
			if !p.opts.skipComments {
				if p.followingTokenIs(token.Comment) && p.followingToken.Line == p.currentToken.Line {
					// if following token is a comment, then it is an inline comment, fetch next token
					p.nextToken()
					d.SetInlineComment(config.InlineComment{
						Value:             p.currentToken.Literal,
						RelativeLineIndex: p.currentToken.Line - directiveLineIndex,
					})
				}
			}
			if iw, ok := p.includeWrappers[d.Name]; ok {
				include, err := iw(d)
				if err != nil {
					return nil, err
				}

				inc, ok := include.(*config.Include)
				if !ok {
					return nil, fmt.Errorf("invalid include wrapper result type %T", include)
				}
				return p.ParseInclude(inc)
			} else if dw, ok := p.directiveWrappers[d.Name]; ok {
				return dw(d)
			}
			return d, nil
		} else if p.curTokenIs(token.Comment) {
			// param comment
			d.SetInlineComment(config.InlineComment{
				Value:             p.currentToken.Literal,
				RelativeLineIndex: p.currentToken.Line - directiveLineIndex,
			})
		} else if p.curTokenIs(token.BlockStart) {
			_, blockSkip1 := SkipValidBlocks[d.Name]
			_, blockSkip2 := p.opts.skipValidSubDirectiveBlock[d.Name]
			isSkipBlockSubDirective := blockSkip1 || blockSkip2 || isSkipValidDirective

			// Special handling for *_by_lua_block directives
			if strings.HasSuffix(d.Name, "_by_lua_block") {
				// For Lua blocks, we need to capture the content without parsing it as nginx directives
				b := &config.Block{
					IsLuaBlock:  true,
					Directives:  []config.IDirective{},
					LiteralCode: "",
				}

				// Skip past the opening brace
				p.nextToken()

				// Collect all content until the matching closing brace
				// We need to count braces to handle nested blocks within Lua code
				braceCount := 1
				var luaCode strings.Builder

				for braceCount > 0 && !p.curTokenIs(token.EOF) {
					if p.curTokenIs(token.BlockStart) {
						braceCount++
					} else if p.curTokenIs(token.BlockEnd) {
						braceCount--
						if braceCount == 0 {
							// This is the closing brace of the Lua block
							break
						}
					}

					// Append token to Lua code if it's not the closing brace
					if !(p.curTokenIs(token.BlockEnd) && braceCount == 0) {
						luaCode.WriteString(p.currentToken.Literal)
						// Add space between tokens for readability
						if p.followingToken.Type != token.BlockEnd &&
							p.followingToken.Type != token.Semicolon &&
							p.followingToken.Type != token.EndOfLine {
							luaCode.WriteString(" ")
						}
					}

					p.nextToken()
				}

				b.LiteralCode = strings.TrimSpace(luaCode.String())
				d.Block = b

				// Use the appropriate wrapper based on the directive name
				if strings.HasSuffix(d.Name, "_by_lua_block") {
					return p.blockWrappers["_by_lua_block"](d)
				}
				return d, nil
			}

			b, err := p.parseBlock(true, isSkipBlockSubDirective)
			if err != nil {
				return nil, err
			}
			d.Block = b

			if bw, ok := p.blockWrappers[d.Name]; ok {
				return bw(d)
			}
			return d, nil
		} else if p.currentToken.Is(token.EndOfLine) {
			continue
		} else {
			return nil, fmt.Errorf("unexpected token %s (%s) on line %d, column %d", p.currentToken.Type.String(), p.currentToken.Literal, p.currentToken.Line, p.currentToken.Column)
		}
	}
}

// ParseInclude just parse include confs
func (p *Parser) ParseInclude(include *config.Include) (config.IDirective, error) {
	if p.opts.parseInclude {
		includePath := include.IncludePath
		if !filepath.IsAbs(includePath) {
			includePath = filepath.Join(p.configRoot, include.IncludePath)
		}
		hasWildcard := hasGlobMeta(includePath)
		includePaths, err := filepath.Glob(includePath)
		if err != nil && !p.opts.skipIncludeParsingErr {
			return nil, err
		}

		for _, matchedPath := range includePaths {
			// Keep parity with nginx include globbing: wildcard includes ignore hidden files.
			if hasWildcard && pathHasHiddenSegment(matchedPath) {
				continue
			}

			canonicalPath, err := filepath.Abs(filepath.Clean(matchedPath))
			if err != nil {
				if p.opts.skipIncludeParsingErr {
					continue
				}
				return nil, err
			}

			if _, inStack := p.includeStack[canonicalPath]; inStack {
				if p.opts.includeCycleErr && !p.opts.skipIncludeParsingErr {
					return nil, fmt.Errorf("include cycle detected for %s", canonicalPath)
				}
				// cyclic include graph, skip this branch and continue.
				continue
			}

			if cached, ok := p.parsedIncludes[canonicalPath]; ok {
				if cached != nil {
					include.Configs = append(include.Configs, cached)
				}
				continue
			}

			p.includeStack[canonicalPath] = struct{}{}

			parser, err := NewParser(canonicalPath,
				WithSameOptions(p),
				withParsedIncludes(p.parsedIncludes),
				withIncludeStack(p.includeStack),
				withConfigRoot(p.configRoot),
			)
			if err != nil {
				delete(p.includeStack, canonicalPath)
				if p.opts.skipIncludeParsingErr {
					continue
				}
				return nil, err
			}

			config, err := parser.Parse()
			delete(p.includeStack, canonicalPath)
			if err != nil {
				if p.opts.skipIncludeParsingErr {
					continue
				}
				return nil, err
			}

			//TODO: link parent config or include direcitve?
			p.parsedIncludes[canonicalPath] = config
			include.Configs = append(include.Configs, config)
		}
	}
	return include, nil
}

func hasGlobMeta(path string) bool {
	return strings.ContainsAny(path, "*?[")
}

func pathHasHiddenSegment(path string) bool {
	cleaned := filepath.Clean(path)
	for _, segment := range strings.Split(cleaned, string(filepath.Separator)) {
		if len(segment) == 0 || segment == "." || segment == ".." {
			continue
		}
		if strings.HasPrefix(segment, ".") {
			return true
		}
	}
	return false
}

// Close closes the file handler and releases the resources
func (p *Parser) Close() (err error) {
	if p.file != nil {
		err = p.file.Close()
	}
	return err
}
