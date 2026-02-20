package dumper

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/imega/luaformatter/formatter"
	"github.com/tufanbarisyildirim/gonginx/config"
)

const hashCommentSentinel = "__GONGINX_HASH_COMMENT__"

var hashCommentRestoreRegex = regexp.MustCompile(`--\s*` + hashCommentSentinel + `\s*`)

// LuaFormatterFunc lets callers override Lua formatting implementation.
type LuaFormatterFunc func(code string, style *Style) (string, error)

// DumpLuaBlock convert a lua block to a string
func DumpLuaBlock(b config.IBlock, style *Style) (luaCode string) {
	luaCode = b.GetCodeBlock()

	if luaCode == "" {
		return ""
	}

	if style.DisableLuaFormatting {
		return strings.TrimRight(luaCode, "\n")
	}

	converted := convertHashComments(luaCode)
	formatted, err := formatLuaCode(converted, style)
	if err != nil {
		// Fallback to original code to preserve semantics when formatter cannot parse.
		return strings.TrimRight(luaCode, "\n")
	}

	formatted = restoreHashComments(formatted)
	return strings.TrimRight(indentLuaCode(formatted, style.StartIndent), "\n")
}

func formatLuaCode(luaCode string, style *Style) (formatted string, err error) {
	if style.LuaFormatter != nil {
		return style.LuaFormatter(luaCode, style)
	}

	defer func() {
		// luaformatter may panic if lua code is not valid.
		if r := recover(); r != nil {
			err = fmt.Errorf("lua formatter panic: %v", r)
		}
	}()

	var buf bytes.Buffer
	cfg := formatter.DefaultConfig()
	cfg.IndentSize = uint8(style.Indent / 4)
	cfg.Highlight = false

	err = formatter.Format(cfg, []byte(luaCode), &buf)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func indentLuaCode(code string, indent int) string {
	if indent <= 0 {
		return code
	}

	lines := bytes.Split([]byte(code), []byte("\n"))
	indentation := bytes.Repeat([]byte(" "), indent)

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
	return indentedBuf.String()
}

func convertHashComments(code string) string {
	lines := strings.SplitAfter(code, "\n")
	if len(lines) == 1 {
		return convertHashCommentsInLine(code)
	}

	var out strings.Builder
	for _, line := range lines {
		out.WriteString(convertHashCommentsInLine(line))
	}
	return out.String()
}

func convertHashCommentsInLine(line string) string {
	commentIdx := findHashCommentIndex(line)
	if commentIdx < 0 {
		return line
	}

	var out strings.Builder
	out.WriteString(line[:commentIdx])
	out.WriteString("-- ")
	out.WriteString(hashCommentSentinel)

	// Preserve any content following # and keep one space before it when needed.
	if commentIdx+1 < len(line) {
		rest := line[commentIdx+1:]
		if len(rest) > 0 && rest[0] != ' ' && rest[0] != '\t' && rest[0] != '\n' && rest[0] != '\r' {
			out.WriteByte(' ')
		}
		out.WriteString(rest)
	}
	return out.String()
}

func restoreHashComments(code string) string {
	restored := hashCommentRestoreRegex.ReplaceAllString(code, "# ")
	restored = strings.ReplaceAll(restored, "# \n", "#\n")
	return strings.TrimSuffix(restored, "# ")
}

func findHashCommentIndex(line string) int {
	inSingleQuote := false
	inDoubleQuote := false
	escaped := false
	longBracketLevel := -1

	for i := 0; i < len(line); i++ {
		ch := line[i]

		if longBracketLevel >= 0 {
			if width, ok := longBracketCloseWidth(line, i, longBracketLevel); ok {
				i += width - 1
				longBracketLevel = -1
			}
			continue
		}

		if inSingleQuote {
			if escaped {
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == '\'' {
				inSingleQuote = false
			}
			continue
		}

		if inDoubleQuote {
			if escaped {
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == '"' {
				inDoubleQuote = false
			}
			continue
		}

		if ch == '\'' {
			inSingleQuote = true
			continue
		}
		if ch == '"' {
			inDoubleQuote = true
			continue
		}

		if level, width, ok := longBracketOpenWidth(line, i); ok {
			longBracketLevel = level
			i += width - 1
			continue
		}

		if ch == '#' && isHashCommentStart(line, i) {
			return i
		}
	}
	return -1
}

func isHashCommentStart(line string, idx int) bool {
	prevNonSpace, hasPrevNonSpace := previousNonSpace(line, idx)
	if !hasPrevNonSpace {
		return true
	}

	nextNonSpace, hasNextNonSpace := nextNonSpace(line, idx+1)
	if !hasNextNonSpace {
		return true
	}

	if idx+1 < len(line) && isSpace(line[idx+1]) {
		return true
	}

	if isOperatorByte(prevNonSpace) {
		return false
	}

	if isLuaKeyword(previousWord(line, idx)) {
		return false
	}

	// Inline comment likely follows an expression and is separated by whitespace.
	if idx > 0 && isSpace(line[idx-1]) && !isOperatorByte(nextNonSpace) {
		return true
	}

	return false
}

func previousNonSpace(s string, idx int) (byte, bool) {
	for i := idx - 1; i >= 0; i-- {
		if isSpace(s[i]) {
			continue
		}
		return s[i], true
	}
	return 0, false
}

func nextNonSpace(s string, idx int) (byte, bool) {
	for i := idx; i < len(s); i++ {
		if isSpace(s[i]) {
			continue
		}
		return s[i], true
	}
	return 0, false
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

func isOperatorByte(ch byte) bool {
	switch ch {
	case '=', '+', '-', '*', '/', '%', '^', '#', '<', '>', '~', '&', '|', ':', ',', '(', '{', '[':
		return true
	default:
		return false
	}
}

func previousWord(line string, idx int) string {
	j := idx - 1
	for j >= 0 && isSpace(line[j]) {
		j--
	}
	if j < 0 {
		return ""
	}

	end := j
	for j >= 0 && isWordByte(line[j]) {
		j--
	}

	return strings.ToLower(line[j+1 : end+1])
}

func isWordByte(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9') ||
		ch == '_'
}

func isLuaKeyword(word string) bool {
	switch word {
	case "if", "then", "do", "while", "repeat", "until", "for", "in", "function", "local", "return", "and", "or", "not", "elseif", "else":
		return true
	default:
		return false
	}
}

func longBracketOpenWidth(line string, idx int) (level int, width int, ok bool) {
	if idx >= len(line) || line[idx] != '[' {
		return 0, 0, false
	}

	j := idx + 1
	level = 0
	for j < len(line) && line[j] == '=' {
		level++
		j++
	}

	if j < len(line) && line[j] == '[' {
		return level, j - idx + 1, true
	}
	return 0, 0, false
}

func longBracketCloseWidth(line string, idx, level int) (width int, ok bool) {
	if idx >= len(line) || line[idx] != ']' {
		return 0, false
	}

	j := idx + 1
	currLevel := 0
	for j < len(line) && line[j] == '=' {
		currLevel++
		j++
	}

	if currLevel == level && j < len(line) && line[j] == ']' {
		return j - idx + 1, true
	}
	return 0, false
}
