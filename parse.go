package saft

import (
	"bufio"
	"github.com/johan-bolmsjo/errors"
	"io"
)

// Parse Saft document from reader. Zero or more root objects are returned in a slice.
func Parse(reader io.Reader) ([]Elem, error) {
	parser := newParser(reader)
	if elems := parser.parseRoot(); parser.err() == nil {
		return elems, nil
	}
	return nil, parser.err()

}

type parser struct {
	lexer      *lexer
	next, prev lexToken
	es         errors.Sink
}

func newParser(reader io.Reader) *parser {
	return &parser{lexer: newLexer(bufio.NewReader(reader))}
}

// Top level parser.
func (p *parser) parseRoot() []Elem {
	// Seed parser with the first lex token.
	p.next = p.lex()

	var elems []Elem
loop:
	for p.es.Ok() {
		switch {
		case p.accept(lexKindSpace):
		case p.isP((*lexToken).isString):
			elems = append(elems, Elem{p.parseString()})
		case p.is(lexKindLBracket):
			elems = append(elems, Elem{p.parseList()})
		case p.is(lexKindLBrace):
			elems = append(elems, Elem{p.parseAssoc()})
		case p.is(lexKindEof):
			break loop
		default:
			p.posError(errors.New("expected string, list or association list"), p.next.pos)
		}
	}
	return elems
}

func (p *parser) parseString() *String {
	p.consume() // Already matched as string
	return &String{pos: p.prev.pos, V: p.prev.s}
}

func (p *parser) parseList() *List {
	p.consume() // Already matched as list
	list := &List{pos: p.prev.pos}

loop:
	for p.es.Ok() {
		switch {
		case p.accept(lexKindSpace):
		case p.isP((*lexToken).isString):
			list.L = append(list.L, Elem{p.parseString()})
		case p.is(lexKindLBracket):
			list.L = append(list.L, Elem{p.parseList()})
		case p.is(lexKindLBrace):
			list.L = append(list.L, Elem{p.parseAssoc()})
		case p.accept(lexKindRBracket):
			break loop
		case p.is(lexKindEof):
			p.posError(errors.New("unterminated list"), p.next.pos)
		default:
			p.posError(errors.New("expected string, list or association list"), p.next.pos)
		}
	}

	return list
}

func (p *parser) parseAssoc() *Assoc {
	p.consume() // Already matched as association list
	assoc := &Assoc{pos: p.prev.pos}

loop:
	for p.es.Ok() {
		switch {
		case p.accept(lexKindSpace):
		case p.accept(lexKindRBrace):
			break loop
		case p.is(lexKindEof):
			p.posError(errors.New("unterminated association list"), p.next.pos)
		default:
			assoc.L = append(assoc.L, p.parsePair())
			if !p.is(lexKindRBrace) && !p.is(lexKindEof) {
				p.expect(lexKindSpace, "association list pairs must be separated by whitespace")
			}
		}
	}

	return assoc
}

func (p *parser) parsePair() Pair {
	keyPred := func(token *lexToken) bool {
		return token.k == lexKindSymbolString || token.k == lexKindInterpString
	}

	p.expectP(keyPred, "key in association list pair must be of symbol or interpreted string form")
	pair := Pair{K: String{pos: p.prev.pos, V: p.prev.s}}

	if !p.accept(lexKindColon) {
		p.posError(errors.New("key in association list pair must be immediately followed by colon"), p.prev.pos)
	}

	// Optional whitespace permitted between colon and value in association list pair.
	p.accept(lexKindSpace)

	switch {
	case p.isP((*lexToken).isString):
		pair.V = Elem{p.parseString()}
	case p.is(lexKindLBracket):
		pair.V = Elem{p.parseList()}
	case p.is(lexKindLBrace):
		pair.V = Elem{p.parseAssoc()}
	case p.is(lexKindRBrace) || p.is(lexKindEof):
		p.posError(errors.New("unterminated association list pair"), p.next.pos)
	default:
		p.posError(errors.New("expected string, list or association list"), p.next.pos)
	}

	// Pair is garbage if an error occurred but it does not really matter
	// since it will be thrown away in that case.
	return pair
}

// Inject error into the parser's error sink with positional information.
func (p *parser) posError(err error, pos LexPos) {
	if err != nil {
		p.es.Send(errors.Wrap(err, pos.String()))
	}
}

// Get any error injected into the error sink.
func (p *parser) err() error {
	return p.es.Cause()
}

// Get token from lexer and store any error in the error sink.
// Not supposed to be used by others than plumbing functions.
func (p *parser) lex() lexToken {
	tok, err := p.lexer.readToken()
	p.es.Send(err)
	return tok
}

// Get the next token from the lexer.
func (p *parser) consume() {
	p.prev = p.next
	if p.next.k != lexKindEof {
		p.next = p.lex()
	}
}

// Check if the next token is of the specified kind.
func (p *parser) is(k lexKind) bool {
	return p.next.k == k
}

// Same as is but with a predicate function.
func (p *parser) isP(pred func(*lexToken) bool) bool {
	return pred(&p.next)
}

// Accept the next token if it's of the specified kind.
// Returns true if the token was accepted.
func (p *parser) accept(k lexKind) (ok bool) {
	if p.next.k == k {
		ok = true
		p.consume()
	}
	return
}

// Expect the next token to be of the specified kind or inject an error
// containing msg into the parser error sink. Returns true if the token was as
// expected.
func (p *parser) expect(k lexKind, msg string) (ok bool) {
	if p.next.k != k {
		p.posError(errors.New(msg), p.next.pos)
		return false
	}
	p.consume()
	return true
}

// Same as expect but with a predicate function.
// Returns true if the token was as expected.
func (p *parser) expectP(pred func(*lexToken) bool, msg string) (ok bool) {
	if !pred(&p.next) {
		p.posError(errors.New(msg), p.next.pos)
		return false
	}
	p.consume()
	return true
}
