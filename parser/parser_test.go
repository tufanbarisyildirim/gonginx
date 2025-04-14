package parser

import (
	"os"
	"testing"

	"github.com/tufanbarisyildirim/gonginx/config"
	"github.com/tufanbarisyildirim/gonginx/dumper"
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
	assert.Assert(t, p.curTokenIs(token.EndOfLine))
	assert.Assert(t, p.followingTokenIs(token.Keyword))
	p.nextToken()
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

	_, err := NewParserFromLexer(
		lex(`
	server { 
	include /but/no/semicolon before block
	}`)).Parse()

	assert.Error(t, err, "unexpected token BlockEnd (}) on line 4, column 2")

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
	_, err := NewParserFromLexer(
		lex(`
	server { 
		location  ~ /(.*)php/{

		} #location with no param

		location  /admin {

			} #location with no param

	}`)).Parse()

	assert.NilError(t, err, "unexpected error")
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
	_, err := NewParserFromLexer(
		lex(`
http{
		server { 
			upstream has multi params /and/ends;
			location ~ /and/ends{
				
			}
		}
}
	`)).Parse()
	assert.NilError(t, err, "unexpected error")
}

func TestParser_Location(t *testing.T) {
	t.Parallel()
	c, err := NewParserFromLexer(
		lex(`
		location ~ /and/ends{
			
		} 
	`)).Parse()
	assert.NilError(t, err, "no error expected here")

	_, ok := c.Directives[0].(*config.Location)
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

	d, ok := c.Directives[0].(*config.Directive)
	assert.Assert(t, ok, "expecting a directive(http) as first statement")
	assert.Equal(t, d.Name, "map", "first directive needs to be ")
	assert.Equal(t, len(d.Parameters), 2, "map must have 2 parameters here")
	assert.Equal(t, d.Parameters[0].GetValue(), "$host", "invalid first parameter")
	assert.Equal(t, d.Parameters[1].GetValue(), "$clientname", "invalid second parameter")
}

func TestParser_UnendedMultiParams(t *testing.T) {
	t.Parallel()
	_, err := NewParserFromLexer(
		lex(`
	server { 
	default with mutli params /but/no/semicolon/to/panic }
	`)).Parse()
	assert.Error(t, err, "unexpected token BlockEnd (}) on line 3, column 55")
}

func TestParser_UnknownDirective(t *testing.T) {
	t.Parallel()
	_, err := NewParserFromLexer(
		lex(`
	server { 
	a_driective param { }
	`)).Parse()
	assert.Error(t, err, "unknown directive 'a_driective' on line 3, column 2")
}

func TestParser_SkipComment(t *testing.T) {
	t.Parallel()
	_, err := NewParserFromLexer(lex(`
if ($a ~* "")#comment
#comment
{#comment
return 400;
}
`)).Parse()

	assert.NilError(t, err, "unexpected error")
}

func TestParser_Include(t *testing.T) {
	t.Parallel()
	p, err := NewParser("../testdata/include-glob/nginx.conf", WithIncludeParsing())
	if err != nil {
		t.Fatal(err)
	}

	c, err := p.Parse()
	assert.NilError(t, err, "no error expected here")
	s := dumper.DumpConfig(c, dumper.IndentedStyle)

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
		_, _ = NewParserFromLexer(
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
	s := dumper.DumpConfig(c, dumper.IndentedStyle)
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

func TestParser_Issue20(t *testing.T) {
	t.Parallel()
	p, err := NewParser("../testdata/issues/20.conf")
	assert.NilError(t, err, "no error expected here")

	c, err := p.Parse()
	assert.NilError(t, err, "no error expected here")

	s := dumper.DumpBlock(c, dumper.IndentedStyle)

	p = NewStringParser(s)
	c, err = p.Parse()
	assert.NilError(t, err, "no error expected here")

	s = dumper.DumpConfig(c, dumper.IndentedStyle)
	assert.Equal(t, `server {
    listen 80;
    listen [::]:80;
    server_name _;
    location / {
        content_by_lua_block {
            # comment
            local foo = "bar" # comment
        }
    }
    location = /random {
        set_by_lua_block $file_name {
            # comment contained unexpect '{'
            local t = ngx.var.uri
            local query = string.find(t, "?", 1)

            if query ~= nil then
             t = string.sub(t, 1, query - 1)
            end

            return t
        }
        set_by_lua_block $random {
            # comment contained unexpect '{'
            return math.random(1, 100)
        }
        return 403 "Random number: $random";
    }
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
	st := &dumper.Style{
		SortDirectives: false,
		StartIndent:    0,
		Indent:         4,
		Debug:          false,
	}
	s := dumper.DumpConfig(c, st)
	assert.Equal(t, `server {
    location = /foo {
        rewrite_by_lua_block {
            res = ngx.location.capture("/memc", {args = {cmd = "incr", key = ngx.var.uri}}) # comment contained unexpect '{'
            # comment contained unexpect '}'
            t = {key = "foo", val = "bar"}
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

func TestParser_Issue32(t *testing.T) {
	t.Parallel()
	p, err := NewParser("../testdata/issues/37.conf",
		WithCustomDirectives("my_custom_directive", "my_custom_directive2"),
		WithCustomDirectives("my_custom_directive3"),
	)
	assert.NilError(t, err, "no error expected here")
	_, err = p.Parse()
	assert.NilError(t, err, "no error expected here")
}

func TestParser_Issue50(t *testing.T) {
	t.Parallel()
	p, err := NewParser("../testdata/issues/50.conf")
	assert.NilError(t, err, "no error expected here")
	data, err := os.ReadFile("../testdata/issues/50.conf")
	assert.NilError(t, err, "no error expected here")
	c, err := p.Parse()
	assert.NilError(t, err, "no error expected here")
	content := dumper.DumpConfig(c, dumper.IndentedStyle)
	assert.Equal(t, content, string(data))
}

func TestParser_SkipBlock(t *testing.T) {
	t.Parallel()
	conf := `user root;
sendfile on;
tcp_nopush on;
map $http_upgrade $connection_upgrade {
	default upgrade;
	'' close;
	test '';
}
fake_map $http_upgrade $connection_upgrade {
	default upgrade;
	'' close;
	test '';
}
`

	p := NewStringParser(conf)
	_, err := p.Parse()
	assert.Error(t, err, "unknown directive 'fake_map' on line 9, column 1")

	p = NewStringParser(conf, WithCustomDirectives("fake_map"), WithSkipValidBlocks("fake_map"))
	_, err = p.Parse()
	assert.NilError(t, err, "no error expected here")
}

func TestParser_SkipBlockSub(t *testing.T) {
	t.Parallel()
	conf := `server {
		listen 80;
		fake_location {
			fake_root /var/www/html;
		}
	}
`
	p := NewStringParser(conf, WithSkipValidBlocks("server"))
	_, err := p.Parse()
	assert.NilError(t, err, "no error expected here")
}

func TestParser_TestFull(t *testing.T) {
	t.Parallel()
	p, err := NewParser("../testdata/full_conf/nginx.conf", WithIncludeParsing())
	assert.NilError(t, err, "no error expected here")

	_, err = p.Parse()
	assert.NilError(t, err, "no error expected here")
}

func TestParser_EventsParent(t *testing.T) {
	p := NewStringParser(`user www www;
worker_processes 5;
error_log logs/error.log;
pid logs/nginx.pid;
worker_rlimit_nofile 8192;
events { worker_connections 4096; } 
# http comment`)

	conf, err := p.Parse()
	assert.NilError(t, err, "no error expected here")

	events := conf.FindDirectives("events")
	assert.Assert(t, len(events) > 0, "cannot find events")
	mainBlock := events[0].GetParent()
	assert.Assert(t, mainBlock == nil, "the events block should not have a parent block")
}

func TestParser_ParentSubDirective1(t *testing.T) {
	p := NewStringParser(`user www www;
events { 
	worker_connections 4096; 
	use epoll;
	} 
# http comment`)

	conf, err := p.Parse()
	assert.NilError(t, err, "no error expected here")

	workerConnections := conf.FindDirectives("worker_connections")
	assert.Assert(t, len(workerConnections) == 1, "cannot find worker_connections")

	events := workerConnections[0].GetParent()
	allDire := events.GetBlock().GetDirectives()
	assert.Assert(t, len(allDire) == 2, "num of sub directive in events error")
}

func TestParser_ParentSubDirective2(t *testing.T) {
	p := NewStringParser(`user www www;
events { 
	worker_connections 4096; 
	use epoll;
	} 

http {
	upstream backend {
		server 127.0.0.1:8080;
		server 127.0.0.1:8081;
	}

	server {
		listen 80;
		location / {
			proxy_pass http://backend/;
		}
	}
}
`)
	conf, err := p.Parse()
	assert.NilError(t, err, "no error expected here")

	server := conf.FindDirectives("server")
	assert.Equal(t, len(server), 1, "num of server error")

	upstreams := conf.FindUpstreams()
	assert.Equal(t, len(upstreams), 1, "num of upstream error")
	uServers := upstreams[0].UpstreamServers
	assert.Equal(t, len(uServers), 2, "num of upstream server error")

	_, ok := uServers[0].Parent.(*config.Upstream)
	assert.Equal(t, ok, true, "cannot convert upstream to blcok")

}

func TestParser_ParentSubDirective3(t *testing.T) {
	p := NewStringParser(`user www www;
http {
	include mime.types;
	server {
		listen 80;
		location / {
			proxy_pass http://backend/;
		}
	}
}
`)
	conf, err := p.Parse()
	assert.NilError(t, err, "no error expected here")

	server := conf.FindDirectives("server")
	assert.Equal(t, len(server), 1, "num of server error")

	httpBlock, ok := server[0].GetParent().(*config.HTTP)
	assert.Equal(t, ok, true, "cannot convert server parent to http")

	includes := httpBlock.FindDirectives("include")

	assert.Equal(t, len(includes), 1, "cannot find include directive in http block")
}

func TestParser_ParentSubDirective4(t *testing.T) {
	p := NewStringParser(`user www www;
http {
	include mime.types;
	server {
		listen 80;
		location / {
			proxy_pass http://backend/;
		}
	}
}
`)
	conf, err := p.Parse()
	assert.NilError(t, err, "no error expected here")

	listens := conf.FindDirectives("listen")
	assert.Equal(t, len(listens), 1, "num of listen error")

	serverIBlock := listens[0].GetParent()
	server, ok := serverIBlock.(*config.Server)

	assert.Equal(t, ok, true, "cannot convert listen parent to server")

	_, ok = server.GetParent().(*config.HTTP)
	assert.Equal(t, ok, true, "cannot convert server parent to http")

}

func TestParser_ParentSubDirective5(t *testing.T) {
	p := NewStringParser(`user www www;
stream {
	upstream ssh_backend {
		server 192.168.1.10:22;
	}
	server {
		listen 2345;
		proxy_pass ssh_backend;
	}
}

`)
	conf, err := p.Parse()
	assert.NilError(t, err, "no error expected here")

	servers := conf.FindDirectives("server")
	assert.Equal(t, len(servers), 1, "num of server error")

}

func TestParser_KeepDataInMultiLine01(t *testing.T) {
	p := NewStringParser(`log_format main '$remote_addr - $remote_user [$time_local] "$request" '
    '$status $body_bytes_sent "$http_referer" '
    '"$http_user_agent" "$http_x_forwarded_for"';`)
	conf, err := p.Parse()
	assert.NilError(t, err, "no error expected here")
	s := dumper.DumpConfig(conf, dumper.IndentedStyle)
	assert.Equal(t, `log_format main '$remote_addr - $remote_user [$time_local] "$request" '
    '$status $body_bytes_sent "$http_referer" '
    '"$http_user_agent" "$http_x_forwarded_for"';`, s)

}

func TestParser_KeepDataInMultiLine02(t *testing.T) {
	p := NewStringParser(`events {
	worker_connections
	4096;}`)
	conf, err := p.Parse()
	assert.NilError(t, err, "no error expected here")
	s := dumper.DumpConfig(conf, dumper.IndentedStyle)
	assert.Equal(t, `events {
    worker_connections
        4096;
}`, s)

}

func TestParser_QuotedString_ISSUE65(t *testing.T) {
	p := NewStringParser(`log_format json_analytics escape=json '{' # json start
	'"msec": "$msec", ' # request unixtime in seconds with a milliseconds resolution
    '"connection": "$connection", ' # connection serial number
    '"connection_requests": "$connection_requests", ' # number of requests made in connection
    '}'; # inline comment
error_log off; # error_log inline comment`)
	conf, err := p.Parse()
	assert.NilError(t, err, "no error expected here")
	s := dumper.DumpConfig(conf, dumper.IndentedStyle)
	assert.Equal(t, `log_format json_analytics escape=json '{' # json start
    '"msec": "$msec", ' # request unixtime in seconds with a milliseconds resolution
    '"connection": "$connection", ' # connection serial number
    '"connection_requests": "$connection_requests", ' # number of requests made in connection
    '}';# inline comment
error_log off;# error_log inline comment`, s)
}

func TestParser_LuaError(t *testing.T) {
	t.Parallel()
	p := NewStringParser(`location / {
        content_by_lua_block { -- comment
local foo = if -- comment }
    }`)
	c, err := p.Parse()
	assert.NilError(t, err, "no error expected here")
	s := dumper.DumpConfig(c, dumper.IndentedStyle)

	assert.Equal(t, `location / {
    content_by_lua_block {
-- comment
local foo = if -- comment
    }
}`, s)
}
