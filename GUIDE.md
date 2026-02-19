# Gonginx Guide

This guide explains how to parse, modify and regenerate Nginx configuration
files with Gonginx. The [examples](./examples) directory contains runnable
programs showing each feature.

## Table of Contents
- [Quick Start](#quick-start)
- [Tips for Coding Agents](#tips-for-coding-agents)
- [Examples](#examples)
- [Library Reference](#library-reference)

## Quick Start
You can find all examples in [examples](./examples).
### Parse nginx config file
Parse Nginx config file, Get server listen port
```go
func parseConfigAndGetPorts(filePath string) ([]string, error) {
	p, err := parser.NewParser(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create parser: %w", err)
	}
	conf, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	servers := conf.FindDirectives("server")
	ports := make([]string, 0)
	for _, server := range servers {
		listens := server.GetBlock().FindDirectives("listen")
		if len(listens) > 0 {
			listenPorts := listens[0].GetParameters()
			ports = append(ports, listenPorts...)
		}
	}
	return ports, nil
}
func main() {
	ports, err := parseConfigAndGetPorts("../../testdata/full_conf/nginx.conf")
	if err != nil {
		panic(err)
	}
	fmt.Println(ports)
}
```

### Dump nginx config string to file with indent
```go
func dumpConfigToFile(fullConf string, filePath string) error {
	p := parser.NewStringParser(fullConf)
	conf, err := p.Parse()
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	dumpString := dumper.DumpConfig(conf, dumper.IndentedStyle)
	if err := os.WriteFile(filePath, []byte(dumpString), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

func dumpAndWriteConfigFile(fullConf string, filePath string) error {
	p := parser.NewStringParser(fullConf)
	conf, err := p.Parse()
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	// set config file path
	conf.FilePath = filePath
	err = dumper.WriteConfig(conf, dumper.IndentedStyle, false)
	if err != nil {
		panic(err)
	}
	return nil
}

func main() {
	fullConf := `user www www;
worker_processes 5;
error_log logs/error.log;
pid logs/nginx.pid;
worker_rlimit_nofile 8192;
events { worker_connections 4096; } http {
include mime.types;
include proxy.conf;
include fastcgi.conf;
index index.html index.htm index.php;
default_type application/octet-stream;
log_format main '$remote_addr - $remote_user [$time_local]  $status '  
'"$request" $body_bytes_sent "$http_referer" '
' "$http_user_agent" "$http_x_forwarded_for"';
access_log logs/access.log main;
sendfile on;
tcp_nopush on;
server_names_hash_bucket_size 128;
server {
listen 80;
server_name domain1.com www.domain1.com;
access_log logs/domain1.access.log main;
root html;
location ~ \.php$ {
fastcgi_pass 127.0.0.1:1025; } } server {
listen 80;
server_name domain2.com www.domain2.com;
access_log logs/domain2.access.log main;
location ~ ^/(images|javascript|js|css|flash|media|static)/ {
root /var/www/virtual/big.server.com/htdocs;
expires 30d;
} location / { proxy_pass http://127.0.0.1:8080; } }
upstream big_server_com {
server 127.0.0.3:8000 weight=5;
server 127.0.0.3:8001 weight=5;
server 192.168.0.1:8000;
server 192.168.0.1:8001;
} server { listen 80;
server_name big.server.com;
access_log logs/big.server.access.log main;
location / { proxy_pass http://big_server_com; } } }`

	// dump config with indented style
	dumpConfigToFile(fullConf, "nginx-temp.conf")

        // dump config to file with indented style
        dumpAndWriteConfigFile(fullConf, "./nginx-temp2.conf")
}
```

## Tips for Coding Agents
The Gonginx API exposes small, composable functions. Prefer operating on
`config.Config` objects and helper methods like `FindDirectives` or `AddServer`
instead of editing raw text.

## Behavior Notes

### Error Model
- Parsing malformed input returns errors.
- Lexer/parser malformed-input paths should not panic.

### Include Parsing
- Enable include parsing with `parser.WithIncludeParsing()`.
- Includes are deduplicated by canonical file path.
- Cyclic include branches are skipped by default to prevent recursion loops.
- Use `parser.WithIncludeCycleErr()` to return an explicit error on cycle detection.

### Dump Sorting
- Sorted dump styles only affect output rendering order.
- Sorted dumps do not mutate in-memory directive order.

### Parent Pointers
- Root-level leaf directives have `nil` parent.
- Nested directives reference their enclosing wrapper/directive as parent.
- This is a hard behavior switch with no compatibility flag.

### Upstream Lookup Modes
- `FindUpstreams()` is permissive and skips unexpected types.
- `FindUpstreamsStrict()` returns a typed error for unexpected upstream directive types.

## Examples

The `examples` directory contains small programs you can run directly. Useful
entries include:

- `adding-server` – append a server to an upstream block.
- `update-directive` – modify a directive in place.
- `dump-nginx-config` – format and write a full configuration.
- `update-server-listen-port` – change the listen port of an existing server.

### Add server in upstream
```go
func main() {
    p := parser.NewStringParser(`http{
	upstream my_backend{
		server 127.0.0.1:443;
		server 127.0.0.2:443 backup;
	}
	}`)

	conf, err := p.Parse()
	if err != nil {
		panic(err)
	}
	upstreams := conf.FindUpstreams()

	upstreams[0].AddServer(&config.UpstreamServer{
		Address: "127.0.0.1:443",
		Parameters: map[string]string{
			"weight": "5",
		},
		Flags: []string{"down"},
	})

	fmt.Println(dumper.DumpBlock(conf.Block, dumper.IndentedStyle))
}
```

### Update directive

```go
func main() {
	p := parser.NewStringParser(`
user nginx;
worker_processes auto;
error_log /var/log/nginx/error.log;
pid /run/nginx.pid;

# Load dynamic modules. See /usr/share/doc/nginx/README.dynamic.
include /usr/share/nginx/modules/*.conf;

events {
    worker_connections 1024;
}

http {
    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';

    access_log  /var/log/nginx/access.log  main;

    sendfile            on;
    tcp_nopush          on;
    tcp_nodelay         on;
    keepalive_timeout   65;
    types_hash_max_size 2048;

    include             /etc/nginx/mime.types;
    default_type        application/octet-stream;

    server {
        listen       80 default_server;
        listen       [::]:80 default_server;
        server_name  _;
        root         /usr/share/nginx/html;

        # Load configuration files for the default server block.
        include /etc/nginx/default.d/*.conf;

        location / {
             proxy_pass http://www.google.com/;
        }

        error_page 404 /404.html;
            location = /40x.html {
        }

        error_page 500 502 503 504 /50x.html;
            location = /50x.html {
        }
    }

}`)

	c, err := p.Parse()
	if err != nil {
		panic(err)
	}
	directives := c.FindDirectives("proxy_pass")
	for _, directive := range directives {
		fmt.Println("found a proxy_pass :  ", directive.GetName(), directive.GetParameters())
		if directive.GetParameters()[0] == "http://www.google.com/" {
			directive.GetParameters()[0] = "http://www.duckduckgo.com/"
		}
	}

	fmt.Println(dumper.DumpBlock(c.Block, dumper.IndentedStyle))

}
```

### Update server listen port
```go
func updateServerListenPort(filePath string, oldPort string, newPort string) (string, error) {
	p, err := parser.NewParser(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create parser: %w", err)
	}
	conf, err := p.Parse()
	if err != nil {
		return "", fmt.Errorf("failed to parse config: %w", err)
	}

	servers := conf.FindDirectives("server")
	for _, server := range servers {
		listens := server.GetBlock().FindDirectives("listen")
		for _, listen := range listens {
			if listen.GetParameters()[0] == oldPort {
				listenDirective := listen.(*config.Directive)
				listenDirective.Parameters[0] = newPort
			}
		}
	}
	changedConf := dumper.DumpConfig(conf, dumper.IndentedStyle)
	return changedConf, nil
}
func main() {

	filePath := "../../testdata/full_conf/nginx.conf"
	oldPort := "80"
	newPort := "8080"
	if changedConf, err := updateServerListenPort(filePath, oldPort, newPort); err != nil {
		log.Fatalf("Error updating server listen port: %v", err)
	} else {
		fmt.Println(changedConf)
	}
}
```

### Add custom directive in any block
```go
func addCustomDirective(fullConf string, blockName string, directiveName string, directiveValue string) (string, error) {
	p := parser.NewStringParser(fullConf)
	conf, err := p.Parse()
	if err != nil {
		return "", fmt.Errorf("failed to parse config: %w", err)
	}

	blocks := conf.FindDirectives(blockName)
	if len(blocks) == 0 {
		return "", fmt.Errorf("no such block: %s", blockName)
	}

	block := blocks[0].GetBlock()
	newDirective := &config.Directive{
		Name:       directiveName,
		Parameters: []string{directiveValue},
	}
	realBlock := block.(*config.Block)
	realBlock.Directives = append(realBlock.Directives, newDirective)

	return dumper.DumpConfig(conf, dumper.IndentedStyle), nil
}

func main() {
	fullConf := `http{
	upstream my_backend{
		server 127.0.0.1:443;
		server 127.0.0.2:443 backup;
	}
	server {
		listen 8080;
		location / {
			root /var/www/html;
			index index.html;
		}
	}
	
	server {
		listen 9090;
		location / {
			root /var/www/html;
			index index.html;
		}
	}
	}`

	blockName := "server"
	directiveName := "access_log"
	directiveValue := "/var/log/nginx/access.log"
	newFullConf, err := addCustomDirective(fullConf, blockName, directiveName, directiveValue)
	if err != nil {
		log.Fatalf("Error adding custom directive: %v", err)
	}
	fmt.Println("New Full Config:", newFullConf)
}

```

## Library Reference
### Parser
Parser is the main package that analyzes and turns nginx structured files into objects. It basically has three pieces: `lexer` breaks the file into tokens and `parser` converts tokens into configuration objects defined in the `config` package.

#### ```NewParser(filePath string, opts ...Option) (*Parser, error)```
+ filePath is the path to the nginx config file.
+ opts is a list of options that can be passed to the parser.
#### ```NewStringParser(content string, opts ...Option) (*Parser, error)```
+ content is the content of the nginx config file.
+ opts is a list of options that can be passed to the parser.
#### ```NewParserFromLexer(lexer *lexer, opts ...Option) *Parser```
+ lexer is the lexer that is used to parse the config file.
+ opts is a list of options that can be passed to the parser.

#### Options
+ **WithSkipIncludeParsingErr()**: If this option is set, the parser will not return an error if it encounters an include directive.
+ **WithDefaultOptions()**: WithDefaultOptions default options
+ **WithSkipComments()**: If this option is set, the parser will not parse comments.
+ **WithIncludeParsing()**: If this option is set, the parser will parse includes.
+ **WithIncludeCycleErr()**: If this option is set, parser returns an error when include cycle is detected.
+ **WithCustomDirectives(directives ...string)**: If this option is set, the parser will parse custom directives without validation.
+ **WithSkipValidBlocks(blocks ...string)**: If this option is set, the parser will not validate directives that are within blocks(recursive)
+ **WithSkipValidDirectivesErr()**: If this option is set, the parser will not return an error if it encounters an invalid directive.

#### Create a new parser with options
```go
p, err := parser.NewParser("nginx.conf", WithSkipComments(), WithCustomDirectives("hello_world"), WithSkipValidBlocks("my_block"))
```

#### ```func (p *Parser) Parse() (*config.Config, error)```
Parse parses the config file(or from config strings) and returns a config object. **It's the only way to get the config object**.

----
### Config
The `config` package models contexts and directives in Go and forms the AST.

#### ```func (c *Config) FindDirectives(directiveName string) []IDirective```
FindDirectives finds all directives with the given name.
#### ```func (c *Config) FindUpstreams() []*Upstream```
FindUpstreams finds all upstreams.

#### IDirective
```go
type IDirective interface {
	GetName() string //the directive name.
	GetParameters() []string
	GetBlock() IBlock
	GetComment() []string
	SetComment(comment []string)
	SetParent(IDirective)
	GetParent() IDirective
}
```

## Migration Notes

If you are upgrading code that depended on older parser/dumper behavior:

- Parent pointer assumptions:
  - root-level leaf directives now have `nil` parent rather than self-reference.
- Include parsing now skips cyclic branches and deduplicates parsed include files by canonical path.
- If you need strict cycle handling, enable `WithIncludeCycleErr()` and treat cycle detection as a parse error.
- Sorted dump operations do not reorder your in-memory AST anymore.
- Lua dump preserves string literal semantics and falls back to original code if formatting fails.
+ GetName() string: the directive name.
+ GetParameters() []string: the directive parameters.
+ GetBlock() IBlock: the directive block.
+ GetComment() []string: the directive comment.
+ SetComment(comment []string): the directive comment.
+ GetParent() IDirective: the directive that contains or encloses the current directive.
#### IBlock
```go
type IBlock interface {
	GetDirectives() []IDirective
	FindDirectives(directiveName string) []IDirective
	GetCodeBlock() string
	SetParent(IDirective)
	GetParent() IDirective
}
```
+ GetDirectives() []IDirective: the block directives.
+ FindDirectives(directiveName string) []IDirective: the block directives.
+ GetCodeBlock() string: the block code.
+ GetParent() IDirective: the directive that contains or encloses the current directive.

#### Directive (impl IDirective)
```go
type Directive struct {
	Block      IBlock
	Name       string
	Parameters []string //TODO: Save parameters with their type
	Comment    []string
	Parent     IBlock
}
```
#### Block (impl IBlock)
```go
type Block struct {
	Directives  []IDirective
	IsLuaBlock  bool
	LiteralCode string
	Parent      IBlock
}
```
+ ```func (b *Block) FindDirectives(directiveName string) []IDirective```

#### Upstream (impl IDirective)
```go
type Upstream struct {
	UpstreamName    string
	UpstreamServers []*UpstreamServer
	//Directives Other directives in upstream (ip_hash; etc)
	Directives []IDirective
	Comment    []string
	Parent     IBlock
}
```
+ ```func (us *Upstream) AddServer(server *UpstreamServer)```

#### UpstreamServer (impl IDirective)
```go
type UpstreamServer struct {
	Address    string
	Flags      []string
	Parameters map[string]string
	Comment    []string
	Parent     IBlock
}
```

#### HTTP (impl IDirective)
```go
type HTTP struct {
	Servers    []*Server
	Directives []IDirective
	Comment    []string
	Parent     IBlock
}
```
+ ```func (h *HTTP) FindDirectives(directiveName string) []IDirective```

#### Server (impl IDirective)
```go
type Server struct {
	Block   IBlock
	Comment []string
	Parent  IBlock
}
```
---
### Dumper
Dumper is the package that holds styling configuration only. 

#### ```func DumpConfig(c *config.Config, style *Style) string```
DumpConfig dump whole config.

#### ```func DumpBlock(b config.IBlock, style *Style) string```
DumpBlock convert a directive to a string.

#### ```func DumpDirective(d config.IDirective, style *Style) string```
DumpDirective convert a directive to a string

#### ```func DumpInclude(i *config.Include, style *Style) map[string]string```
DumpInclude dump(stringify) the included AST

#### ```func WriteConfig(c *config.Config, style *Style, writeInclude bool) error```
WriteConfig writes config.

#### Style
dumping style, you can use it to customize the output style.
```go
type Style struct {
	SortDirectives    bool
	SpaceBeforeBlocks bool
	StartIndent       int
	Indent            int
	Debug             bool
}
```
#### Styles by default
+ NoIndentStyle
```go
NoIndentStyle = &Style{
    SortDirectives: false,
    StartIndent:    0,
    Indent:         0,
    Debug:          false,
}
```
+ IndentStyle
```go
IndentedStyle = &Style{
    SortDirectives: false,
    StartIndent:    0,
    Indent:         4,
    Debug:          false,
}
```
+ NoIndentSortedStyle
```go
NoIndentSortedStyle = &Style{
    SortDirectives: true,
    StartIndent:    0,
    Indent:         0,
    Debug:          false,
}
```
+ NoIndentSortedSpaceStyle
```go
NoIndentSortedSpaceStyle = &Style{
    SortDirectives:    true,
    SpaceBeforeBlocks: true,
    StartIndent:       0,
    Indent:            0,
    Debug:             false,
}
```
