package dumper

import (
	"strings"
	"testing"

	"github.com/tufanbarisyildirim/gonginx/config"
	"gotest.tools/v3/assert"
)

func TestDumpLuaBlock_PreservesLiteralsAndHashComments(t *testing.T) {
	t.Parallel()

	block := &config.Block{
		IsLuaBlock: true,
		LiteralCode: `local s = "--keep"
local t = "#hash"
local u = '--stay'
# comment
return 1`,
	}

	got := DumpLuaBlock(block, &Style{
		StartIndent: 4,
		Indent:      4,
	})

	assert.Assert(t, strings.Contains(got, `local s = "--keep"`))
	assert.Assert(t, strings.Contains(got, `local t = "#hash"`))
	assert.Assert(t, strings.Contains(got, `local u = '--stay'`))
	assert.Assert(t, strings.Contains(got, `# comment`))
	assert.Assert(t, !strings.Contains(got, `"#keep"`))
	assert.Assert(t, !strings.Contains(got, hashCommentSentinel))
}

func TestDumpLuaBlock_PreservesHashOperator(t *testing.T) {
	t.Parallel()

	block := &config.Block{
		IsLuaBlock: true,
		LiteralCode: `local n = #arr
return n`,
	}

	got := DumpLuaBlock(block, &Style{
		StartIndent: 4,
		Indent:      4,
	})

	assert.Assert(t, strings.Contains(got, "local n = #arr"))
	assert.Assert(t, !strings.Contains(got, hashCommentSentinel))
}

func TestDumpLuaBlock_InvalidLuaFallsBackToOriginal(t *testing.T) {
	t.Parallel()

	original := `-- comment
local foo = if -- comment`

	block := &config.Block{
		IsLuaBlock:  true,
		LiteralCode: original,
	}

	got := DumpLuaBlock(block, &Style{
		StartIndent: 8,
		Indent:      4,
	})

	assert.Equal(t, got, original)
}
