package main

import (
	"fmt"

	"github.com/tufanbarisyildirim/gonginx/dumper"
	"github.com/tufanbarisyildirim/gonginx/parser"
)

func main() {
	p := parser.NewStringParser(`user www www;
worker_processes 5;
error_log logs/error.log;
pid logs/nginx.pid;
worker_rlimit_nofile 8192;
events { worker_connections 4096; } 
# http comment
http {
# include comment
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
# comment: location
location ~ \.php$ {
fastcgi_pass 127.0.0.1:1025; } } 
map $http_upgrade $connection_upgrade {
	default upgrade;
	'' close;
	test '';
}
# comment: server
server {
listen 80;
server_name domain2.com www.domain2.com;
access_log logs/domain2.access.log main;
location ~ ^/(images|javascript|js|css|flash|media|static)/ {
root /var/www/virtual/big.server.com/htdocs;
expires 30d;
} location / { proxy_pass http://127.0.0.1:8080; } }
# comment: big_server_com
# comment: upstream big_server_com
upstream big_server_com {
server 127.0.0.3:8000 weight=5; # inline comment: server 127.0.0.3:8000 weight=5
server 127.0.0.3:8001 weight=5; # inline comment: server 127.0.0.3:8000 weight=5
server 192.168.0.1:8000;
server 192.168.0.1:8001;
}
# comment: server
server { # comment: listen
listen 80;
server_name big.server.com;
# comment: access_log
access_log logs/big.server.access.log main;
location / { proxy_pass http://big_server_com; } } }`)

	c, err := p.Parse()
	if err != nil {
		panic(err)
	}
	fmt.Println(dumper.DumpConfig(c, dumper.IndentedStyle))

}
