package main

import (
	"fmt"
	"log"

	"github.com/tufanbarisyildirim/gonginx/config"
	"github.com/tufanbarisyildirim/gonginx/dumper"
	"github.com/tufanbarisyildirim/gonginx/parser"
)

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
			if listen.GetParameters()[0].GetValue() == oldPort {
				listenDirective := listen.(*config.Directive)
				listenDirective.Parameters[0].SetValue(newPort)
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
