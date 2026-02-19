package config

import "fmt"

// Config represents a complete nginx configuration file.
type Config struct {
	*Block
	FilePath string
}

// Global wrappers provide extension points for custom directive handling.
var (
	BlockWrappers     = map[string]func(*Directive) (IDirective, error){}
	DirectiveWrappers = map[string]func(*Directive) (IDirective, error){}
	IncludeWrappers   = map[string]func(*Directive) (IDirective, error){}
)

//TODO(tufan): move that part inti dumper package
//SaveToFile save config to a file
//TODO: add custom file / folder path support
//func (c *Config) SaveToFile(style *dumper.Style) error {
//	//wrilte file
//	dirPath := filepath.Dir(c.FilePath)
//	if _, err := os.Stat(dirPath); err != nil {
//		err := os.MkdirAll(dirPath, os.ModePerm)
//		if err != nil {
//			return err //TODO: do we reallt need to find a way to test dir creating error?
//		}
//	}
//
//	//write main file
//	err := ioutil.WriteFile(c.FilePath, c.ToByteArray(style), 0644)
//	if err != nil {
//		return err //TODO: do we need to find a way to test writing error?
//	}
//
//	//write sub files (inlude /file/path)
//	for _, directive := range c.Block.Directives {
//		if fs, ok := (interface{}(directive)).(FileDirective); ok {
//			err := fs.SaveToFile(style)
//			if err != nil {
//				return err
//			}
//		}
//	}
//
//	return nil
//}

// FindDirectives find directives from whole config block
func (c *Config) FindDirectives(directiveName string) []IDirective {
	return c.Block.FindDirectives(directiveName)
}

// FindUpstreams find directives from whole config block
func (c *Config) FindUpstreams() []*Upstream {
	var upstreams []*Upstream
	directives := c.Block.FindDirectives("upstream")
	for _, directive := range directives {
		upstream, ok := directive.(*Upstream)
		if !ok {
			continue
		}

		upstreams = append(upstreams, upstream)
	}
	return upstreams
}

// UnexpectedUpstreamTypeError indicates a directive named "upstream" that is
// not represented by *config.Upstream in the AST.
type UnexpectedUpstreamTypeError struct {
	Index int
	Got   any
}

func (e *UnexpectedUpstreamTypeError) Error() string {
	return fmt.Sprintf("unexpected upstream directive type at index %d: %T", e.Index, e.Got)
}

// FindUpstreamsStrict returns all upstream blocks or an error when a matched
// "upstream" directive is not typed as *Upstream.
func (c *Config) FindUpstreamsStrict() ([]*Upstream, error) {
	var upstreams []*Upstream
	directives := c.Block.FindDirectives("upstream")
	for i, directive := range directives {
		upstream, ok := directive.(*Upstream)
		if !ok {
			return nil, &UnexpectedUpstreamTypeError{
				Index: i,
				Got:   directive,
			}
		}
		upstreams = append(upstreams, upstream)
	}
	return upstreams, nil
}

func init() {
	BlockWrappers["http"] = func(directive *Directive) (IDirective, error) {
		return NewHTTP(directive)
	}
	BlockWrappers["location"] = func(directive *Directive) (IDirective, error) {
		return NewLocation(directive)
	}
	BlockWrappers["_by_lua_block"] = func(directive *Directive) (IDirective, error) {
		return NewLuaBlock(directive)
	}
	BlockWrappers["server"] = func(directive *Directive) (IDirective, error) {
		return NewServer(directive)
	}
	BlockWrappers["upstream"] = func(directive *Directive) (IDirective, error) {
		return NewUpstream(directive)
	}

	DirectiveWrappers["server"] = func(directive *Directive) (IDirective, error) {
		return NewUpstreamServer(directive)
	}

	IncludeWrappers["include"] = func(directive *Directive) (IDirective, error) {
		return NewInclude(directive)
	}
}
