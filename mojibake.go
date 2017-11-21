package mojibake

import (
	"errors"
	"html"
	"unicode/utf8"
)

func FixHTMLEntities(s string) string {
	return html.UnescapeString(s)
}

var (
	// ErrUTF8 is raised if the input has rune errors
	ErrUTF8 = errors.New("UTF-8 encoding error")

	// ErrImpure is raised if the string is not purely double UTF-8 encoded
	// Impurity criterias:
	// - some runes have values above 255
	// - some consecutive runes with value < 256 do not combine to make a valid rune
	ErrImpure = errors.New("FixDoubleUTF8: skip (impure input)")
)

const (
	// Limits of first byte of UTF-8 encoded Unicode codepoint outside the ASCII range
	UTF8FirstByteMin = 194 // \U00000080 in UTF-8 starts with byte 194
	UTF8FirstByteMax = 244 // \U0010FFFF in UTF-8 starts with byte 244
)

// FixDoubleUTF8 fixes double UTF-8 encoding issues in-place.
//
// All precautions are taken: nothing is changed if the input is not
// purely double encoded.
//
// In case of error, buf is not changed and is just returned.
// In case of success and double UTF-8 was found, the returned
// slice will be shorter than the input.
//
// Two errors may be returned:
//   - ErrUTF8: this is not a valid UTF-8 string
//   - ErrImpure: this is a valid UTF-8 string, but above, some rune do not make a
//     purely double encoded rune
func FixDoubleUTF8(buf []byte) ([]byte, error) {
	impure := false
	x := [utf8.UTFMax]byte{} // Buffer for decoding a rune. Max rune: U+10FFFF => 4 UTF-8 bytes
	first := -1
	last := -1
	i := 0
	// First pass: check that the buffer contains clean UTF-8
	// and that all runes above 127 are below 255 and can be combined
	// with the next one to make an UTF-8 char.
Main:
	for i < len(buf) {
		if buf[i] < utf8.RuneSelf {
			i++
			continue
		}
		r, n := utf8.DecodeRune(buf[i:])
		if r == utf8.RuneError {
			return buf, ErrUTF8
		}
		i += n
		if impure {
			continue
		}
		if r < UTF8FirstByteMin || r > UTF8FirstByteMax {
			impure = true
			continue
		}
		if first < 0 {
			first = i - n
		}
		x[0] = byte(r)
		j := 1
		for {
			r, n = utf8.DecodeRune(buf[i:])
			if n != 2 {
				impure = true
				break
			}
			if r == utf8.RuneError {
				return buf, ErrUTF8
			}
			i += n
			if r < 128 || r > 255 {
				impure = true
				continue Main
			}
			x[j] = byte(r)
			j++
			r, n = utf8.DecodeRune(x[:j])
			if n == j && r != utf8.RuneError {
				break
			}
			if j == len(x) || !utf8.ValidRune(r) {
				impure = true
				break
			}
		}
		last = i - 1
	}

	if last < 0 {
		return buf, nil
	}
	if impure {
		return buf, ErrImpure
	}

	// Second pass: fix in-place buf[first:last]
	out := buf[:first]
	i = first
	for i <= last {
		if buf[i] < utf8.RuneSelf {
			out = append(out, buf[i])
			i++
			continue
		}
		r, n := utf8.DecodeRune(buf[i:])
		i += n
		x[0] = byte(r)
		j := 1
		for {
			r, n = utf8.DecodeRune(buf[i:])
			i += n
			x[j] = byte(r)
			j++
			r, n = utf8.DecodeLastRune(x[:j])
			if n == j && r != utf8.RuneError {
				break
			}
		}
		out = append(out, x[:j]...)
	}
	// Append remaining buf[last+1:]
	if last < len(buf)-1 {
		out = append(out, buf[last+1:]...)
	}
	return out, nil
}
