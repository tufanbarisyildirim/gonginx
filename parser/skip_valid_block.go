package parser

import "strings"

var skipValidBlocks = `types
map
`

// SkipValidBlocks defines a list of valid blocks to be skipped during initialization.
// This string is split by newline characters, with each trimmed block name added to the SkipValidBlocks mapping.
var SkipValidBlocks map[string]struct{} = map[string]struct{}{}

func init() {
	blocks := strings.Split(skipValidBlocks, "\n")
	for _, block := range blocks {
		SkipValidBlocks[strings.TrimSpace(block)] = struct{}{}
	}
}
