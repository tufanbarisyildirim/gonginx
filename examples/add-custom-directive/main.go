package main

import (
	"fmt"
	"log"

	"github.com/tufanbarisyildirim/gonginx/config"
	"github.com/tufanbarisyildirim/gonginx/dumper"
	"github.com/tufanbarisyildirim/gonginx/parser"
)

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
		Parameters: []config.Parameter{{Value: directiveValue}},
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
