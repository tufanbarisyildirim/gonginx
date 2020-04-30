package main

import (
	"github.com/tufanbarisyildirim/gonginx/dumper"
	"github.com/tufanbarisyildirim/gonginx/parser"
)

func main() {
	p, err := parser.NewParser("./full-example/formatting/raw.conf")
	if err != nil {
		panic(err)
	}
	c := p.Parse()
	c.FilePath = "./full-example/formatting/formatted.conf" //move this to savefile method.
	c.SaveToFile(&dumper.Style{
		SortDirectives:    true,
		SpaceBeforeBlocks: true,
		StartIndent:       0,
		Indent:            4,
	})
}
