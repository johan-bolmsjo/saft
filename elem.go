package saft

import (
	"fmt"
)

// Elem is any of the Saft data types String, List or Assoc.
type Elem struct {
	any elem
}

type elem interface {
	Pos() LexPos
	elemType() elemType
}

type elemType int

const (
	elemTypeString elemType = iota
	elemTypeList
	elemTypeAssoc
)

func (et elemType) String() string {
	switch et {
	case elemTypeString:
		return "string"
	case elemTypeList:
		return "list"
	case elemTypeAssoc:
		return "assoc"
	}
	return "?"
}

// Pos returns the lexed position of the element.
// Useful for error reporting if the element does not match an expected type.
func (e Elem) Pos() LexPos {
	return e.any.Pos()
}

// IsString asserts that e is a String.
// Return values as regular type assertions.
func (e Elem) IsString() (t *String, ok bool) {
	t, ok = e.any.(*String)
	return
}

// IsList asserts that e is a List.
// Return values as regular type assertions.
func (e Elem) IsList() (t *List, ok bool) {
	t, ok = e.any.(*List)
	return
}

// IsAssoc asserts that e is an Assoc.
// Return values as regular type assertions.
func (e Elem) IsAssoc() (t *Assoc, ok bool) {
	t, ok = e.any.(*Assoc)
	return
}

func expectError(e *Elem, expected elemType) error {
	pos := e.Pos()
	return fmt.Errorf("%s: expected %s, found %s", &pos, expected, e.any.elemType())
}

// ExpectString asserts and expects that e is a String.
// Returns a string or an error containing positional information.
func (e Elem) ExpectString() (t *String, err error) {
	if t, _ = e.any.(*String); t == nil {
		err = expectError(&e, elemTypeString)
	}
	return
}

// ExpectList asserts and expects that e is a List.
// Returns a List or an error containing positional information.
func (e Elem) ExpectList() (t *List, err error) {
	if t, _ = e.any.(*List); t == nil {
		err = expectError(&e, elemTypeList)
	}
	return
}

// ExpectAssoc asserts and expects that e is an Assoc.
// Returns an Assoc or an error containing positional information.
func (e Elem) ExpectAssoc() (t *Assoc, err error) {
	if t, _ = e.any.(*Assoc); t == nil {
		err = expectError(&e, elemTypeAssoc)
	}
	return
}
