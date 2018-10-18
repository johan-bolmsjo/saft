package saft

// List is a regular list of elements.
type List struct {
	pos LexPos
	L   []Elem // List with elements.
}

// Pos returns positional information useful for context dependent error reporting.
func (l *List) Pos() LexPos {
	return l.pos
}

func (l *List) elemType() elemType {
	return elemTypeList
}
