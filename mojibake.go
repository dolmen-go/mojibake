package mojibake

import (
	"errors"
	"html"
	"unicode/utf8"
)

func FixHTMLEntities(s string) string {
	return html.UnescapeString(s)
}

var ErrUTF8 = errors.New("UTF-8 encoding error")
var ErrImpure = errors.New("FixDoubleUTF8: skip (impure input)")

// FixDoubleUTF8 fixes double UTF-8 encoding in-place
func FixDoubleUTF8(buf []byte) ([]byte, error) {
	impure := false
	x := [4]byte{} // Buffer for decoding a rune. Max rune: U+10FFFF => 4 UTF-8 bytes
	first := -1
	last := -1
	i := 0
	// First pass: check that the buffer contains clean UTF-8
	// and that all runes above 127 are below 255 and can be combined
	// with the next one to make an UTF-8 char.
Main:
	for i < len(buf) {
		if buf[i] < 128 {
			i++
			continue
		}
		r, n := utf8.DecodeRune(buf[i:])
		if r == utf8.RuneError {
			return buf, ErrUTF8
		}
		i += n
		if r > 255 {
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
		if buf[i] < 128 {
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
