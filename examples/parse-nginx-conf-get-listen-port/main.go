package main

import (
	"fmt"

	"github.com/tufanbarisyildirim/gonginx/parser"
)

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
			for _, port := range listenPorts {
				ports = append(ports, port.GetValue())
			}
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
