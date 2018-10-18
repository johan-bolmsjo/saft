package saft

// Assoc is an association list.
type Assoc struct {
	pos LexPos
	L   Pairs // Key value pairs.
}

// Pos returns positional information useful for context dependent error reporting.
func (al *Assoc) Pos() LexPos {
	return al.pos
}

func (al *Assoc) elemType() elemType {
	return elemTypeAssoc
}

// Pairs is the list of pairs in an association list.
type Pairs []Pair

// Pair is a key value pair in an association list.
type Pair struct {
	K String // Key
	V Elem   // Value
}

// Find first pair in list of pairs with the specified key.
// Returns a list cut so that the found pair is first or nil if no pair was found.
func (lst Pairs) Find(key string) Pairs {
	return lst.FindP(func(pairKey string) bool { return key == pairKey })
}

// FindP is the same as Find but accepting a predicate function matching against pair keys.
// The predicate function should return true to indicate a match.
func (lst Pairs) FindP(pred func(key string) bool) Pairs {
	for i, v := range lst {
		if pred(v.K.V) {
			return lst[i:]
		}
	}
	return nil
}
