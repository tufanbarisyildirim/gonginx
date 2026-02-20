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

func TestDumpLuaBlock_DisableFormatting(t *testing.T) {
	t.Parallel()

	original := "local  x=1\nreturn x\n"

	block := &config.Block{
		IsLuaBlock:  true,
		LiteralCode: original,
	}

	style := (&Style{
		StartIndent: 8,
		Indent:      4,
	}).WithLuaFormatting(false)

	got := DumpLuaBlock(block, style)

	assert.Equal(t, got, strings.TrimRight(original, "\n"))
}

func TestDumpLuaBlock_UsesCustomFormatter(t *testing.T) {
	t.Parallel()

	block := &config.Block{
		IsLuaBlock:  true,
		LiteralCode: "return 1",
	}

	called := false
	style := (&Style{
		StartIndent: 4,
		Indent:      4,
	}).WithLuaFormatter(func(_ string, _ *Style) (string, error) {
		called = true
		return "return 42", nil
	})

	got := DumpLuaBlock(block, style)

	assert.Assert(t, called)
	assert.Equal(t, got, "    return 42")
}
