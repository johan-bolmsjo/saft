package saft

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func lex(input string) ([]lexToken, error) {
	var tokens []lexToken

	lex := newLexer(strings.NewReader(input))
	var done bool
	for !done {
		tok, err := lex.readToken()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)
		done = tok.k == lexKindEof || err != nil
	}
	return tokens, nil
}

func tokensToString(tokens []lexToken) string {
	var sb strings.Builder
	for _, tok := range tokens {
		fmt.Fprintf(&sb, "%s: %s\n", tok.pos.String(), &tok)
	}
	return sb.String()
}

func checkTokens(t *testing.T, tokens []lexToken, want string) {
	got := strings.TrimSpace(tokensToString(tokens))
	want = strings.TrimSpace(want)
	if got != want {
		t.Fatalf("got lexed tokens:\n%s\nwant:\n%s\n", got, want)
	}
}

func checkError(t *testing.T, err error, want string) {
	got := "nil"
	if err != nil {
		got = err.Error()
	}
	if got != want {
		t.Fatalf("got error %q; want %q", got, want)
	}
}

// Test that space tokens are merged.
// e.g. [space, comment, space] -> [space]
func TestLex_MergeSpace(t *testing.T) {
	tokens, err := lex(`x //

x`)
	checkTokens(t, tokens, `
1:0: <symbol-string> "x"
1:1: <space>
3:0: <symbol-string> "x"
3:1: <eof>
`)
	checkError(t, err, "nil")
}

// Comments are note emitted as lexed tokens.
func TestLex_Comment(t *testing.T) {
	tokens, err := lex(`{foo:bar // Comment
foo: baz//Comment next to symbol string form
}`)
	checkTokens(t, tokens, `
1:0: {
1:1: <symbol-string> "foo"
1:4: :
1:5: <symbol-string> "bar"
1:8: <space>
2:0: <symbol-string> "foo"
2:3: :
2:4: <space>
2:5: <symbol-string> "baz"
2:44: <space>
3:0: }
3:1: <eof>
`)
	checkError(t, err, "nil")
}

// Test that leading and trailing space is emitted.
func TestLex_LeadingTrailingSpace(t *testing.T) {
	tokens, err := lex(` x `)
	checkTokens(t, tokens, `
1:0: <space>
1:1: <symbol-string> "x"
1:2: <space>
1:3: <eof>
`)
	checkError(t, err, "nil")
}

func TestLex_SyntaxTokens(t *testing.T) {
	tokens, err := lex(`{}[]:`)
	checkTokens(t, tokens, `
1:0: {
1:1: }
1:2: [
1:3: ]
1:4: :
1:5: <eof>
`)
	checkError(t, err, "nil")
}

func TestLex_SymbolString(t *testing.T) {
	tokens, err := lex(`/a/b/c`)
	checkTokens(t, tokens, `
1:0: <symbol-string> "/a/b/c"
1:6: <eof>
`)
	checkError(t, err, "nil")
}

func TestLex_InterpString(t *testing.T) {
	tokens, err := lex(`"interpreted string \n\r\t\\\""`)
	checkTokens(t, tokens, `
1:0: <interp-string> "interpreted string \n\r\t\\\""
1:31: <eof>
`)
	checkError(t, err, "nil")
}

func TestLex_InterpStringUnterminatedStringError(t *testing.T) {
	tokens, err := lex(`"`)
	checkTokens(t, tokens, ``)
	checkError(t, err, "1:1: unterminated string")

	tokens, err = lex(`"\`)
	checkTokens(t, tokens, ``)
	checkError(t, err, "1:2: unterminated string")
}

func TestLex_InterpStringUnknownEscapeSequenceError(t *testing.T) {
	tokens, err := lex(`"\x"`)
	checkTokens(t, tokens, ``)
	checkError(t, err, "1:2: unknown escape sequence")
}

func TestLex_InterpStringNewlineInStringError(t *testing.T) {
	tokens, err := lex(`"
"`)
	checkTokens(t, tokens, ``)
	checkError(t, err, "1:1: newline in string")
}

const Q = "`"

func TestLex_RawString(t *testing.T) {
	tokens, err := lex(Q + `raw string \x
` + Q)
	checkTokens(t, tokens, `
1:0: <raw-string> "raw string \\x\n"
2:1: <eof>
`)
	checkError(t, err, "nil")
}

func TestLex_RawStringUnterminatedStringError(t *testing.T) {
	tokens, err := lex(Q)
	checkTokens(t, tokens, ``)
	checkError(t, err, "1:1: unterminated string")
}

func TestLex_JoinedStringError(t *testing.T) {
	// There are many combinations of this but the same code is used to detect them all.
	var tbl = []struct{ input, error string }{
		{`a"a"`, "1:1: strings must be separated"},
		{`"a"a`, "1:3: strings must be separated"},
		{"a`a`", "1:1: strings must be separated"},
	}

	for i, td := range tbl {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tokens, err := lex(td.input)
			checkTokens(t, tokens, ``)
			checkError(t, err, td.error)
		})
	}
}

const TAB = "\t"

// Check that TAB count as 8 spaces and that the position is aligned to an even
// multiple of that number. e.g. SPC + SPC + TAB -> 8.
func TestLex_TabPos(t *testing.T) {
	tokens, err := lex(TAB + `1` + TAB + `1234567` + TAB)
	checkTokens(t, tokens, `
1:0: <space>
1:8: <symbol-string> "1"
1:9: <space>
1:16: <symbol-string> "1234567"
1:23: <space>
1:24: <eof>
`)
	checkError(t, err, "nil")
}
