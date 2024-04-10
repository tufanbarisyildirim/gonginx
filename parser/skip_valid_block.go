package parser

import "strings"

var skipValidBlocks = `types
map
`

var SkipValidBlocks map[string]struct{} = map[string]struct{}{}

func init() {
	blocks := strings.Split(skipValidBlocks, "\n")
	for _, block := range blocks {
		SkipValidBlocks[strings.TrimSpace(block)] = struct{}{}
	}
}
