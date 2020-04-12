# Gonginx
![reportcard](https://goreportcard.com/badge/github.com/tufanbarisyildirim/gonginx) [![Actions Status](https://github.com/tufanbarisyildirim/gonginx/workflows/Go/badge.svg)](https://github.com/tufanbarisyildirim/gonging/actions)


Gonginx is an Nginx configuration parser helps you to parse, edit, regenerate your nginx config in your go applications.

## Supporting Blocks/Directives
Generated a to-do/feature list from a full nginx config examle to track how is going.


- [ ] `user       www www;  ## Default: nobody`
- [ ] `worker_processes  5;  ## Default: 1`
- [ ] `error_log  logs/error.log;`
- [ ] `pid        logs/nginx.pid;`
- [ ] `worker_rlimit_nofile 8192;`

- [ ] `events {`
  - [ ] `worker_connections  4096;  ## Default: 1024`
`}`

- [ ] `http {`
  - [x] `include    fastcgi.conf;`
  - [ ] `index    index.html index.htm index.php;`

  - [ ] ```log_format   main '$remote_addr - $remote_user [$time_local]  $status '
    '"$request" $body_bytes_sent "$http_referer" '
    '"$http_user_agent" "$http_x_forwarded_for"';```
  - [ ] `access_log   logs/access.log  main;`
  - [ ] `sendfile     on;`
  - [ ] `tcp_nopush   on;`
  - [ ] `server_names_hash_bucket_size 128; # this seems to be required for some vhosts`

  - [ ] `server { # php/fastcgi`
    - [ ] `listen       80;`
    - [ ] `server_name  domain1.com www.domain1.com;`
    - [ ] `access_log   logs/domain1.access.log  main;`
    - [ ] `root         html;`

    - [ ] `location ~ \.php$ {`
      - [ ] `fastcgi_pass   127.0.0.1:1025;`
    `}`
    - [ ] `location / {`
      - [ ] `proxy_pass      http://127.0.0.1:8080;`
    `}`
  `}`

  - [ ] `upstream big_server_com {`
    - [x] `server 127.0.0.3:8000;`
    - [ ] `server 127.0.0.3:8001 weight=5;`
  `}`
`}`


## Limitations
-- a list is coming soon, but it will support most of core functions of nginx

# [Contributing](CONTRIBUTING.md)

## License
[MIT License](LICENSE)