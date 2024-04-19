package main

import (
	"fmt"
	"os"

	"github.com/tufanbarisyildirim/gonginx/dumper"
	"github.com/tufanbarisyildirim/gonginx/parser"
)

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

	// dump config to file whit indented style
	dumpAndWriteConfigFile(fullConf, "./nginx-temp2.conf")
}
