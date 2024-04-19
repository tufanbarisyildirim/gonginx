package main

import (
	"fmt"

	"github.com/tufanbarisyildirim/gonginx/dumper"
	"github.com/tufanbarisyildirim/gonginx/parser"
)

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
