package decode

import (
	"encoding/base64"
	"strings"
)

const b64Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
const b64Standard = "+/"
const b64Padding = "="
const b64URL = "-_"

const b64name = "b64"

// Base64 takes a decoder and an input string
type Base64 struct {
	dec   *decoder
	input string
}

// NewB64CodecC state machine to smartly decode a string with invalid chars
// and different variants
// nolint: gocyclo
func NewB64CodecC(in string) CodecC {
	const (
		runInvalid runType = iota
		runAlphabet
		runStandard
		runUrl
	)

	// emit writes into output what was read up until this point and move l.start to l.pos
	emit := func(d *decoder, t runType) {
		token := d.input[d.start : d.pos-d.width]
		if len(token) == 0 {
			return
		}

		token = strings.TrimRight(token, b64Padding)
		var decodefunc func(string) []byte

		switch t {
		case runAlphabet, runStandard:
			decodefunc = func(in string) []byte {
				if len(in) < 2 {
					return []byte(genInvalid(len(in)))
				}
				encoding := base64.RawStdEncoding
				buf, err := encoding.DecodeString(in)
				if err != nil {
					return []byte(err.Error())
				}
				return buf
			}

		case runUrl:
			decodefunc = func(in string) []byte {
				if len(in) < 2 {
					return []byte(genInvalid(len(in)))
				}
				encoding := base64.RawURLEncoding
				buf, err := encoding.DecodeString(in)
				if err != nil {
					return []byte(err.Error())
				}
				return buf
			}

		case runInvalid:
			decodefunc = func(in string) []byte {
				return []byte(genInvalid(len(in)))
			}
		}

		d.out.Write(decodefunc(token))
		d.start = d.pos - d.width
	}

	var (
		startState    stateFn
		alphabetState stateFn
		standardState stateFn
		urlState      stateFn
	)

	startState = func(d *decoder) stateFn {
		switch n := d.next(); {
		case strings.ContainsRune(b64Alphabet, n):
			emit(d, runInvalid)
			return alphabetState

		case strings.ContainsRune(b64Standard, n):
			emit(d, runInvalid)
			return standardState

		case strings.ContainsRune(b64URL, n):
			emit(d, runInvalid)
			return urlState

		case n == eof:
			emit(d, runInvalid)
			return nil

		default:
			return startState
		}
	}

	alphabetState = func(d *decoder) stateFn {
		switch n := d.next(); {
		case strings.ContainsRune(b64Alphabet+b64Padding, n):
			return alphabetState

		case strings.ContainsRune(b64Standard, n):
			return standardState

		case strings.ContainsRune(b64URL, n):
			return urlState

		case n == eof:
			emit(d, runAlphabet)
			return nil

		default:
			emit(d, runAlphabet)
			return startState
		}
	}

	standardState = func(d *decoder) stateFn {
		switch n := d.next(); {
		case strings.ContainsRune(b64Alphabet+b64Standard+b64Padding, n):
			return standardState

		case n == eof:
			emit(d, runStandard)
			return nil

		default:
			emit(d, runStandard)
			return startState
		}
	}

	urlState = func(d *decoder) stateFn {
		switch n := d.next(); {
		case strings.ContainsRune(b64Alphabet+b64URL+b64Padding, n):
			return urlState

		case n == eof:
			emit(d, runUrl)
			return nil

		default:
			emit(d, runUrl)
			return startState
		}
	}

	return &Base64{
		dec:   newDecoder(in, startState),
		input: in,
	}
}

// Name returns the name of the codec
func (b *Base64) Name() string {
	return b64name
}

// Decode a valid b64 string
func (b *Base64) Decode() (output string) {
	return string(b.dec.decode())
}

// Encode a string into b64 with StdEncodig set
func (b *Base64) Encode() (output string) {
	//TODO allow user to decide which encoder
	return base64.StdEncoding.EncodeToString([]byte(b.input))
}

// Check returns the percentage of valid b16 characters in the input string
func (b *Base64) Check() (acceptability float64) {
	var c int
	var tot int
	for _, r := range b.input {
		tot++
		if strings.ContainsRune(b64Alphabet+b64Standard+b64URL+b64Padding, r) {
			c++
		}
	}
	//Heuristic to consider uneven strings as less likely to be valid base64
	if delta := tot % 4; delta != 0 {
		tot += delta
	}
	return float64(c) / float64(tot)
}
