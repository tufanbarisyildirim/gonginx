package dumper

import (
	"bytes"
	"log"
	"runtime"
	"strings"

	"github.com/imega/luaformatter/formatter"
	"github.com/tufanbarisyildirim/gonginx/config"
)

// DumpLuaBlock convert a lua block to a string
func DumpLuaBlock(b config.IBlock, style *Style) (luaCode string) {
	luaCode = b.GetCodeBlock()

	defer func() {
		// luaformatter may panic if the lua code is not valid
		if r := recover(); r != nil {
			buf := make([]byte, 1024)
			runtime.Stack(buf, false)
			log.Printf("%s\n%s", r, buf)
		}
	}()
	var buf bytes.Buffer

	if luaCode == "" {
		return ""
	}

	config := formatter.DefaultConfig()
	config.IndentSize = uint8(style.Indent / 4)
	config.Highlight = false

	// Replace # comments with -- comments temporarily
	luaCode = strings.ReplaceAll(luaCode, "#", "--")

	formatter.Format(config, []byte(luaCode), &buf)

	formatted := buf.String()

	// Add indentation to each line
	lines := bytes.Split([]byte(formatted), []byte("\n"))
	indentation := bytes.Repeat([]byte(" "), style.StartIndent)

	var indentedBuf bytes.Buffer
	for i, line := range lines {
		if len(line) > 0 {
			indentedBuf.Write(indentation)
			indentedBuf.Write(line)
		}
		if i < len(lines)-1 {
			indentedBuf.WriteByte('\n')
		}
	}

	formatted = indentedBuf.String()

	// Restore # comments
	formatted = strings.ReplaceAll(formatted, "--", "#")

	return strings.TrimRight(formatted, "\n")
}
