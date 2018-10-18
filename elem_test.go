package saft_test

import (
	"bufio"
	"bytes"
	"github.com/johan-bolmsjo/saft"
	"testing"
)

func getTestElem(t *testing.T, s string) saft.Elem {
	elems, err := saft.Parse(bufio.NewReader(bytes.NewBuffer([]byte(s))))
	if len(elems) != 1 {
		panic(err)
	}
	return elems[0]
}

const stringValue = `abc`

func getStringElem(t *testing.T) saft.Elem { return getTestElem(t, stringValue) }
func getListElem(t *testing.T) saft.Elem   { return getTestElem(t, `[]`) }
func getAssocElem(t *testing.T) saft.Elem  { return getTestElem(t, `{}`) }

func TestElem_String(t *testing.T) {
	stringElem := getStringElem(t)
	listElem := getListElem(t)

	s, err := stringElem.ExpectString()
	checkError(t, "stringElem.ExpectString()", err, "nil")
	if s.V != stringValue {
		t.Fatalf("stringElem.ExpectString() = %q; want %q", s.V, stringValue)
	}

	s, err = listElem.ExpectString()
	checkError(t, "listElem.ExpectString()", err, "1:0: expected string, found list")
	if s != nil {
		t.Fatalf("listElem.ExpectString() = %v; want nil", s)
	}

	if _, ok := stringElem.IsString(); !ok {
		t.Fatalf("stringElem.IsString() = false")
	}
	if _, ok := listElem.IsString(); ok {
		t.Fatalf("listElem.IsString() = true")
	}
}

func TestElem_List(t *testing.T) {
	stringElem := getStringElem(t)
	listElem := getListElem(t)
	assocElem := getAssocElem(t)

	_, err := listElem.ExpectList()
	checkError(t, "listElem.ExpectList()", err, "nil")

	_, err = stringElem.ExpectList()
	checkError(t, "stringElem.ExpectList()", err, "1:0: expected list, found string")

	_, err = assocElem.ExpectList()
	checkError(t, "assocElem.ExpectList()", err, "1:0: expected list, found assoc")

	if _, ok := listElem.IsList(); !ok {
		t.Fatalf("listElem.IsList() = false")
	}
	if _, ok := stringElem.IsList(); ok {
		t.Fatalf("stringElem.IsList() = true")
	}
}

func TestElem_Assoc(t *testing.T) {
	stringElem := getStringElem(t)
	assocElem := getAssocElem(t)

	_, err := assocElem.ExpectAssoc()
	checkError(t, "assocElem.ExpectAssoc()", err, "nil")

	_, err = stringElem.ExpectAssoc()
	checkError(t, "stringElem.ExpectAssoc()", err, "1:0: expected assoc, found string")

	if _, ok := assocElem.IsAssoc(); !ok {
		t.Fatalf("assocElem.IsAssoc() = false")
	}
	if _, ok := stringElem.IsAssoc(); ok {
		t.Fatalf("stringElem.IsAssoc() = true")
	}
}
