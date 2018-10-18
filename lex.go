package saft

import (
	"fmt"
	"github.com/johan-bolmsjo/errors"
	"io"
	"strings"
	"unicode"
)

// lexer contains state for lexing from an io.RuneReader
type lexer struct {
	rr          io.RuneReader
	unreadRunes []lexRune // Uused for peeking
	pos         LexPos    // Current line and column based on read runes
	prevKind    lexKind   // Used to merge whitespace tokens
	eof         bool      // EOF seen?
	es          errors.Sink
}

// newLexer returns a new lexer lexing from the given input stream.
func newLexer(rr io.RuneReader) *lexer {
	return &lexer{rr: rr, pos: LexPos{Line: 1}}
}

// End of file is coded using a sentinel rune value.
const runeEof rune = -1

func (lex *lexer) emitPosError(err error, pos LexPos) {
	if err != nil {
		lex.eof = true
		lex.es.Send(errors.Wrap(err, pos.String()))
	}
}

func (lex *lexer) readRune() lexRune {
	// Pop unread runes first
	if l := len(lex.unreadRunes); l > 0 {
		lr := lex.unreadRunes[l-1]
		lex.unreadRunes = lex.unreadRunes[:l-1]
		return lr
	}

	lr := lexRune{runeEof, lex.pos}

	// EOF is sticky (EOF set on error as well)
	if !lex.eof {
		r, _, err := lex.rr.ReadRune()
		if err != nil {
			lex.eof = true
			if err != io.EOF {
				lex.emitPosError(err, lex.pos)
			}
		} else {
			lr.r = r
			lex.pos.update(r)
		}
	}

	return lr
}

func (lex *lexer) unreadRune(lr lexRune) {
	if lr.r != runeEof {
		lex.unreadRunes = append(lex.unreadRunes, lr)
	}
}

func (lex *lexer) peekRune() lexRune {
	lr := lex.readRune()
	lex.unreadRune(lr)
	return lr
}

// readToken returns the next token from the input stream.
// lexKindEof is returned as the last token, possibly with an error.
func (lex *lexer) readToken() (tok lexToken, err error) {
	defer func() {
		// Return any error from the error sink
		if err = lex.es.Cause(); err != nil {
			lex.prevKind = lexKindEof
			tok = lexToken{k: lex.prevKind, pos: lex.pos}
		}
	}()

	var lr lexRune
	lrTok := func(k lexKind) lexToken { return lexToken{k: k, pos: lr.pos} }

begin:
	switch r := lr.read(lex); {
	case r == runeEof:
		tok = lrTok(lexKindEof)
	case r == '/':
		var lr2 lexRune
		if lr2.read(lex) == '/' {
			lex.lexComment()
			goto begin
		}
		lr2.unread(lex)
		tok = lex.lexSymbolString(lr) // The only thing that match is a symbol string starting with '/'
	case isSpace(r):
		tok = lex.lexSpace(lr)
		if lex.prevKind == lexKindSpace {
			// "whitespace + comment + whitespace" would result in
			// two space tokens being emitted without this merging.
			goto begin
		}
	case r == ':':
		tok = lrTok(lexKindColon)
	case r == '{':
		tok = lrTok(lexKindLBrace)
	case r == '}':
		tok = lrTok(lexKindRBrace)
	case r == '[':
		tok = lrTok(lexKindLBracket)
	case r == ']':
		tok = lrTok(lexKindRBracket)
	case r == '"':
		tok = lex.lexInterpretedString(lr)
	case r == '`':
		tok = lex.lexRawString(lr)
	default:
		tok = lex.lexSymbolString(lr)
	}
	lex.prevKind = tok.k
	return
}

// lexSpace scans a run of space characters.
// One space has already been consumed and is passed as firstRune.
func (lex *lexer) lexSpace(firstRune lexRune) lexToken {
	var lr lexRune
	for isSpace(lr.read(lex)) {
	}
	lr.unread(lex)
	return lexToken{k: lexKindSpace, pos: firstRune.pos}
}

// Check that double strings like "a""b" a"b" etc are not allowed. There must
// always be some separator between two strings. This is caught in the lexer
// because the parser would produce confusing errors for some cases such as a"a"
// as a key in an association list.
func (lex *lexer) checkJoinedStrings(pos LexPos) (errToken lexToken, ok bool) {
	if lex.prevKind.isString() {
		lex.emitPosError(errors.New("strings must be separated"), pos)
		return lexToken{k: lexKindEof, pos: lex.pos}, false
	}
	return lexToken{}, true
}

func (lex *lexer) lexInterpretedString(firstRune lexRune) lexToken {
	if errToken, ok := lex.checkJoinedStrings(firstRune.pos); !ok {
		return errToken
	}

	var sb strings.Builder
	var lr lexRune

	parseEscape := false

runeLoop:
	for lr.read(lex) != runeEof {
		if parseEscape {
			switch lr.r {
			case '\\', '"':
				sb.WriteByte(byte(lr.r))
			case 'n':
				sb.WriteByte('\n')
			case 'r':
				sb.WriteByte('\r')
			case 't':
				sb.WriteByte('\t')
			default:
				lex.emitPosError(errors.New("unknown escape sequence"), lr.pos)
				return lexToken{k: lexKindEof, pos: lex.pos}
			}
			parseEscape = false
		} else {
			switch lr.r {
			case '\\':
				parseEscape = true
			case '\n', '\r':
				lex.emitPosError(errors.New("newline in string"), lr.pos)
				return lexToken{k: lexKindEof, pos: lex.pos}
			case '"':
				break runeLoop
			default:
				sb.WriteRune(lr.r)
			}
		}
	}

	if lr.r != '"' {
		lex.emitPosError(errors.New("unterminated string"), lr.pos)
		return lexToken{k: lexKindEof, pos: lex.pos}
	}

	return lexToken{k: lexKindInterpString, s: sb.String(), pos: firstRune.pos}
}

func (lex *lexer) lexRawString(firstRune lexRune) lexToken {
	if errToken, ok := lex.checkJoinedStrings(firstRune.pos); !ok {
		return errToken
	}

	var sb strings.Builder
	var lr lexRune

runeLoop:
	for lr.read(lex) != runeEof {
		if lr.r == '`' {
			break runeLoop
		} else {
			sb.WriteRune(lr.r)
		}
	}

	if lr.r != '`' {
		lex.emitPosError(errors.New("unterminated string"), lr.pos)
		return lexToken{k: lexKindEof, pos: lex.pos}
	}

	return lexToken{k: lexKindRawString, s: sb.String(), pos: firstRune.pos}
}

func (lex *lexer) lexSymbolString(firstRune lexRune) lexToken {
	if errToken, ok := lex.checkJoinedStrings(firstRune.pos); !ok {
		return errToken
	}

	predicate := func(r rune) bool {
		return !(isSpace(r) || r == '\\' || r == '`' || r == '"' ||
			r == '{' || r == '}' || r == '[' || r == ']' || r == ':' ||
			r == '/' || r == runeEof)
	}

	var sb strings.Builder
	sb.WriteRune(firstRune.r)

	var lr lexRune
keepScanning:
	for predicate(lr.read(lex)) {
		sb.WriteRune(lr.r)
	}

	// Comments need to terminate scanning but a single '/' is allowed
	if lr.r == '/' {
		if lex.peekRune().r != '/' {
			sb.WriteRune(lr.r)
			goto keepScanning
		}
	}

	lr.unread(lex)
	return lexToken{k: lexKindSymbolString, s: sb.String(), pos: firstRune.pos}
}

// lexComment scans characters until EOL or EOF.
// The comment marker '//' has already been consumed.
func (lex *lexer) lexComment() {
	var lr lexRune
	predicate := func(r rune) bool { return r != '\n' && r != runeEof }
	for predicate(lr.read(lex)) {
	}
	lr.unread(lex)
}

// lexToken is a lexed token emitted to the parser.
type lexToken struct {
	k   lexKind
	s   string // Content if applicable
	pos LexPos // Positional information
}

func (tok *lexToken) String() string {
	if tok.s == "" {
		return tok.k.String()
	}
	return fmt.Sprintf("%s %q", tok.k, tok.s)
}

func (tok *lexToken) isString() bool {
	return tok.k.isString()
}

// lexKind is the type of token emitted.
type lexKind int8

const (
	lexKindEof   lexKind = iota // End of input stream
	lexKindSpace                // Consecutive whitespace including new line
	lexKindSymbolString
	lexKindInterpString
	lexKindRawString
	lexKindColon
	lexKindLBrace
	lexKindRBrace
	lexKindLBracket
	lexKindRBracket
)

var lexKindItoa = map[lexKind]string{
	lexKindEof:          "<eof>",
	lexKindSpace:        "<space>",
	lexKindSymbolString: "<symbol-string>",
	lexKindInterpString: "<interp-string>",
	lexKindRawString:    "<raw-string>",
	lexKindColon:        ":",
	lexKindLBrace:       "{",
	lexKindRBrace:       "}",
	lexKindLBracket:     "[",
	lexKindRBracket:     "]",
}

func (kind lexKind) String() string {
	return lexKindItoa[kind]
}

func (kind lexKind) isString() bool {
	return kind == lexKindSymbolString || kind == lexKindInterpString || kind == lexKindRawString
}

type lexRune struct {
	r   rune
	pos LexPos
}

func (lr *lexRune) read(lex *lexer) rune {
	*lr = lex.readRune()
	return lr.r
}

func (lr *lexRune) unread(lex *lexer) {
	lex.unreadRune(*lr)
	lr.r = runeEof
}

// LexPos contain line and column information of lexed runes and tokens.
type LexPos struct {
	Line, Column int32
}

// update position based on input rune.
func (pos *LexPos) update(r rune) {
	switch r {
	case '\t': // Count tabs using a 8 character tab stop
		pos.Column += 8 - (pos.Column % 8)
	case '\n':
		pos.Line++
		pos.Column = 0
	default:
		pos.Column++
	}
}

// String implements the fmt.Stringer interface.
func (pos *LexPos) String() string {
	return fmt.Sprintf("%d:%d", pos.Line, pos.Column)
}

func isSpace(r rune) bool { return unicode.IsSpace(r) }
