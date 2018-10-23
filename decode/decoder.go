package decode

import (
	"bytes"
	"strings"
	"unicode/utf8"
)

const eof = -1

var invalid = '�'

// Pos is an integer that define a position inside a string
type Pos int
type runType int

func genInvalid(n int) (inv string) {
	return strings.Repeat(string(invalid), n)
}

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*decoder) stateFn

// decoder holds the state of the scanner.
type decoder struct {
	input string  // the string being scanned
	state stateFn // the next lexing function to enter
	pos   Pos     // current position in the input
	start Pos     // start position of this item
	width Pos     // width of last rune read from input
	out   *bytes.Buffer
}

// Constructs a new decoder
func newDecoder(input string, startState stateFn) *decoder {
	return &decoder{
		input: input,
		state: startState,
		out:   bytes.NewBuffer(nil),
	}
}

// next returns the next rune in the input.
func (l *decoder) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	return r
}

// decode runs the decode until EOF
func (l *decoder) decode() []byte {
	for l.state != nil {
		l.state = l.state(l)
	}
	return l.out.Bytes()
}
