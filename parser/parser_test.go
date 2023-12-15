package parser

import (
	"testing"

	"github.com/tufanbarisyildirim/gonginx"
	"github.com/tufanbarisyildirim/gonginx/parser/token"
	"gotest.tools/v3/assert"
)

func TestParser_CurrFollow(t *testing.T) {
	t.Parallel()
	conf := `
	server { # simple reverse-proxy
	}
	`
	p := NewStringParser(conf)
	//assert.Assert(t, tokens, 1)
	assert.Assert(t, p.curTokenIs(token.Keyword))
	assert.Assert(t, p.followingTokenIs(token.BlockStart))
}

//TODO(tufan): reactivate here after getting include specific things done
//func TestParser_Include(t *testing.T) {
//	conf := `
//	include /etc/ngin/conf.d/mime.types;
//	`
//	p := NewStringParser(conf)
//	c := p.Parse()
//	_, ok := c.Directives[0].(gonginx.IncludeDirective) //we expect the first statement to be an include
//	assert.Assert(t, ok)
//}

func TestParser_UnendedInclude(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	NewParserFromLexer(
		lex(`
	server { 
	include /but/no/semicolon before block;
	`)).Parse()
}

func TestParser_LocationNoParam(t *testing.T) {
	t.Parallel()

	_, err := NewParserFromLexer(
		lex(`
	server { 
	location  {} #location with no param
	`)).Parse()

	assert.Error(t, err, "no enough parameter for location")
}

func TestParser_LocationTooManyParam(t *testing.T) {
	t.Parallel()
	_, err := NewParserFromLexer(
		lex(`
	server { 
	location one two three four {} #location with too many arguments
	`)).Parse()
	assert.Error(t, err, "too many arguments for location directive")
}

func TestParser_ParseValidLocations(t *testing.T) {
	t.Parallel()
	NewParserFromLexer(
		lex(`
	server { 
		location  ~ /(.*)php/{

		} #location with no param

		location  /admin {

			} #location with no param

	`)).Parse()
}

func TestParser_ParseUpstream(t *testing.T) {
	t.Parallel()
	_, err := NewParserFromLexer(
		lex(`
		upstream my_upstream{
			server 127.0.0.1:8080;
			server 127.0.0.1:8080 weight=5 failure=3;
		}
	server { 
		location  ~ /(.*)php/{

		} #location with no param

		location  /admin {

			} #location with no param

	}`)).Parse()
	assert.NilError(t, err, "no error expected here")
}

func TestParser_ParseFromFile(t *testing.T) {
	t.Parallel()
	p, err := NewParser("../full-example/nginx.conf")
	assert.NilError(t, err)
	assert.Assert(t, p.file != nil, "file must be non-nil")
	_, err2 := NewParser("../full-example/nginx.conf-not-found")
	assert.ErrorContains(t, err2, "no such file or directory")
}

func TestParser_MultiParamDirecive(t *testing.T) {
	t.Parallel()
	NewParserFromLexer(
		lex(`
http{
		server { 
			a_directive has multi params /and/ends;
			location ~ /and/ends{
				
			}
		}
}
	`)).Parse()
}

func TestParser_Location(t *testing.T) {
	t.Parallel()
	c, err := NewParserFromLexer(
		lex(`
		location ~ /and/ends{
			
		} 
	`)).Parse()
	assert.NilError(t, err, "no error expected here")

	_, ok := c.Directives[0].(*gonginx.Location)
	assert.Assert(t, ok, "expecting a location as first statement")
}

func TestParser_VariableAsParameter(t *testing.T) {
	t.Parallel()
	c, err := NewParserFromLexer(
		lex(`
			map $host $clientname {
				default -;
			}
	`)).Parse()

	assert.NilError(t, err, "no error expected here")

	d, ok := c.Directives[0].(*gonginx.Directive)
	assert.Assert(t, ok, "expecting a directive(http) as first statement")
	assert.Equal(t, d.Name, "map", "first directive needs to be ")
	assert.Equal(t, len(d.Parameters), 2, "map must have 2 parameters here")
	assert.Equal(t, d.Parameters[0], "$host", "invalid first parameter")
	assert.Equal(t, d.Parameters[1], "$clientname", "invalid second parameter")
}

func TestParser_UnendedMultiParams(t *testing.T) {
	t.Parallel()
	_, err := NewParserFromLexer(
		lex(`
	server { 
	a_driective with mutli params /but/no/semicolon/to/panic }
	`)).Parse()
	assert.Error(t, err, "unexpected token BlockEnd (}) on line 3, column 59")
}

func TestParser_SkipComment(t *testing.T) {
	t.Parallel()
	NewParserFromLexer(lex(`
if ($a ~* "")#comment
#comment
{#comment
return 400;
}
`)).Parse()
}

func TestParser_Include(t *testing.T) {
	t.Parallel()
	p, err := NewParser("../testdata/include-glob/nginx.conf", WithIncludeParsing())
	if err != nil {
		t.Fatal(err)
	}

	c, err := p.Parse()
	assert.NilError(t, err, "no error expected here")
	s := gonginx.DumpConfig(c, gonginx.IndentedStyle)

	assert.Equal(t, `user www www;
worker_processes 5;
include events.conf;
include http.conf;`, s)
}

func Benchmark_ParseFullExample(t *testing.B) {
	fullconf := `user www www;
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
	for n := 0; n < t.N; n++ {
		NewParserFromLexer(
			lex(fullconf)).Parse()
	}
}

func TestParser_Issue17(t *testing.T) {
	t.Parallel()
	p, err := NewParser("../testdata/issues/17.conf", WithSkipComments())
	if err != nil {
		t.Fatal(err)
	}

	c, err := p.Parse()
	assert.NilError(t, err, "no error expected here")
	s := gonginx.DumpConfig(c, gonginx.IndentedStyle)
	assert.Equal(t, `location / {
    set $serve_URL $fullurl${uri}index.html;
    try_files $serve_URL $uri $uri/ /index.php$is_args$args;
}
location ~* ^/xmlrpc.php$ {
    return 403;
}
location ~ \.php$ {
    include snippets/fastcgi-php.conf;
    fastcgi_param PHP_VALUE "open_basedir=$document_root:/tmp;\nerror_log=/public_html/logs/demo1-php_errors.log;";
    fastcgi_pass unix:/run/php/php7.4-fpm.sock;
}
location /wp-content/uploads/ {
    location ~ .(aspx|php|jsp|cgi)$ {
        return 410;
    }
}
location ~* \.(css|gif|ico|jpeg|jpg|js|png|woff|woff2|ttf|ttc|otf|eot)$ {
    expires 30d;
    log_not_found off;
}
location ~ /\.ht {
    deny all;
}
location ~ ^/(status)$ {
    allow 127.0.0.1;
    fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
    fastcgi_index index.php;
    include fastcgi_params;
    fastcgi_pass unix:/run/php/php7.4-fpm.sock;
}
location /.well-known/acme-challenge {
    alias /public_html/certbot_temp/.well-known/acme-challenge;
}`, s)
}

func TestParser_Issue22(t *testing.T) {
	t.Parallel()
	p, err := NewParser("../testdata/issues/22.conf")
	if err != nil {
		t.Fatal(err)
	}

	c, err := p.Parse()
	assert.NilError(t, err, "no error expected here")
	st := &gonginx.Style{
		SortDirectives: false,
		StartIndent:    0,
		Indent:         4,
		Debug:          false,
	}
	s := gonginx.DumpConfig(c, st)
	assert.Equal(t, `server {
    location = /foo {
        rewrite_by_lua_block {
            
            res = ngx.location.capture("/memc",
            { args = { cmd = "incr", key = ngx.var.uri } } # comment contained unexpect '{'
            # comment contained unexpect '}'
            )
            t = { key="foo", val="bar" }
            
        }
    }
}`, s)
}

func TestParser_Issue31(t *testing.T) {
	t.Parallel()
	p, err := NewParser("../testdata/issues/31.conf")
	assert.NilError(t, err, "no error expected here")

	_, err = p.Parse()
	if err == nil {
		t.Fatal("error expected here")
	}
}
