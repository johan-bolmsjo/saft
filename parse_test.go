package saft_test

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/johan-bolmsjo/saft"
	"strings"
	"testing"
)

func parse(t *testing.T, input string) ([]saft.Elem, error) {
	t.Logf("Input:\n%s", input)
	return saft.Parse(bufio.NewReader(bytes.NewBuffer([]byte(input))))
}

func checkElems(t *testing.T, elems []saft.Elem, want string) {
	got := strings.TrimSpace(elemsToString(elems))
	want = strings.TrimSpace(want)
	if got != want {
		t.Fatalf("got parsed elems:\n%s\nwant:\n%s\n", got, want)
	}
}

func errorString(err error) string {
	if err != nil {
		return err.Error()
	}
	return "nil"
}

func checkError(t *testing.T, ctx string, err error, want string) {
	got := errorString(err)
	if got != want {
		t.Fatalf("%s error = %q; want %q", ctx, got, want)
	}
}

func checkParseError(t *testing.T, err error, want string) {
	checkError(t, "parse()", err, want)
}

func elemsToString(elems []saft.Elem) string {
	var sb strings.Builder
	for _, e := range elems {
		elemToSb(&sb, e, 0)
	}
	return sb.String()
}

func elemToSb(sb *strings.Builder, elem saft.Elem, indentLevel int) {
	if s, ok := elem.IsString(); ok {
		fmt.Fprintf(sb, "%q ", s.V)
	} else if l, ok := elem.IsList(); ok {
		sb.WriteByte('[')
		for _, elem = range l.L {
			elemToSb(sb, elem, indentLevel+1)
		}
		sb.WriteByte(']')
	} else if a, ok := elem.IsAssoc(); ok {
		sb.WriteByte('{')
		for _, pair := range a.L {
			fmt.Fprintf(sb, "%q:", pair.K.V)
			elemToSb(sb, pair.V, indentLevel+1)
			sb.WriteByte(' ')
		}
		sb.WriteByte('}')
	}

	if indentLevel == 0 {
		sb.WriteByte(' ')
	}
}

const Q = "`"

func TestParse_RootString(t *testing.T) {
	elems, err := parse(t, `a "a b" `+Q+`a
b`+Q)
	checkElems(t, elems, `"a"  "a b"  "a\nb"`)
	checkParseError(t, err, "nil")
}

func TestParse_RootList(t *testing.T) {
	elems, err := parse(t, `
[] [[]] [[][]]
[a a] [a "a a" `+Q+`a
a`+Q+`]
[a[a[a]]] [a [a [a]]]
[a {a:a a:[a{a:a}]}]`)

	checkElems(t, elems, `[] [[]] [[][]] ["a" "a" ] ["a" "a a" "a\na" ] ["a" ["a" ["a" ]]] ["a" ["a" ["a" ]]] ["a" {"a":"a"  "a":["a" {"a":"a"  }] }]`)
	checkParseError(t, err, "nil")
}

func TestParse_RootAssocList(t *testing.T) {
	elems, err := parse(t, `
{}
{a:b a:c}
{a: {x:y} b: [i j k]}
`)
	checkElems(t, elems, `{} {"a":"b"  "a":"c"  } {"a":{"x":"y"  } "b":["i" "j" "k" ] }`)
	checkParseError(t, err, "nil")
}

func TestParse_RootUnknownTypeError(t *testing.T) {
	elems, err := parse(t, `:`)
	checkElems(t, elems, ``)
	checkParseError(t, err, "1:0: expected string, list or association list")
}

func TestParse_ListUnknownTypeError(t *testing.T) {
	elems, err := parse(t, `[:`)
	checkElems(t, elems, ``)
	checkParseError(t, err, "1:1: expected string, list or association list")
}

func TestParse_RootUnterminatedListError(t *testing.T) {
	elems, err := parse(t, `[`)
	checkElems(t, elems, ``)
	checkParseError(t, err, "1:1: unterminated list")
}

func TestParse_RootUnterminatedAssocListError(t *testing.T) {
	elems, err := parse(t, `{`)
	checkElems(t, elems, ``)
	checkParseError(t, err, "1:1: unterminated association list")
}

func TestParse_RootUnterminatedAssocListPairError(t *testing.T) {
	elems, err := parse(t, `{a:`)
	checkElems(t, elems, ``)
	checkParseError(t, err, "1:3: unterminated association list pair")
}

func TestParse_RootAssocListPairColonError(t *testing.T) {
	elems, err := parse(t, `{a :`)
	checkElems(t, elems, ``)
	checkParseError(t, err, "1:1: key in association list pair must be immediately followed by colon")
}

func TestParse_RootAssocListKeyStringFormError(t *testing.T) {
	elems, err := parse(t, `{`+Q+`a`+Q+`:`)
	checkElems(t, elems, ``)
	checkParseError(t, err, "1:1: key in association list pair must be of symbol or interpreted string form")
}
