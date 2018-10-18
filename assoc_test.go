package saft_test

import (
	//	"github.com/johan-bolmsjo/saft"
	"fmt"
	"strings"
	"testing"
)

func TestAssoc_Find(t *testing.T) {
	elems, err := parse(t, `{a:c d:c a:?}`)
	checkParseError(t, err, "nil")

	if len(elems) != 1 {
		t.Fatalf("expected 1 parsed element; got %d", len(elems))
	}

	assoc, err := elems[0].ExpectAssoc()
	if err != nil {
		t.Fatalf("expected parsed association list, got error: %s", err)
	}

	var sb strings.Builder

	l := assoc.L.Find("a")
	for l != nil {
		p := &l[0]
		fmt.Fprintf(&sb, "%s: ", p.K.V)
		if s, ok := p.V.IsString(); ok {
			fmt.Fprintf(&sb, "%s", s.V)
		}
		sb.WriteByte('\n')
		l = l[1:].Find("a")
	}

	got := sb.String()
	want := "a: c\na: ?\n"
	if got != want {
		t.Fatalf("found association list pairs:\n%s\nwant:\n%s\n", got, want)
	}
}
