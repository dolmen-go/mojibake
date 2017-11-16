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

// FixDoubleUTF8 fixes double UTF-8 encoding in place
func FixDoubleUTF8(buf []byte) ([]byte, error) {
	impure := false
	x := [2]byte{}
	first := -1
	last := -1
	i := 0
	// First pass: check that the buffer contains clean UTF-8
	// and that all runes above 127 are below 255 and can be combined
	// with the next one to make an UTF-8 char.
	for i < len(buf) {
		if buf[i] < 127 {
			i++
			continue
		}
		r, n := utf8.DecodeRune(buf[i:])
		if r == utf8.RuneError {
			return buf, ErrUTF8
		}
		if r > 255 {
			impure = true
			i += n
			continue
		}
		if first < 0 {
			first = i
		}
		x[0] = byte(r)
		i += n
		r, n = utf8.DecodeRune(buf[i:])
		if n == 0 {
			return buf, nil
		}
		if r == utf8.RuneError {
			return buf, ErrUTF8
		}
		i += n
		last = i - 1
		if r < 128 || r > 255 {
			impure = true
			continue
		}
		x[1] = byte(r)
		r, n = utf8.DecodeRune(x[:])
		if n == 0 || r == utf8.RuneError {
			impure = true
		}
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
		if buf[i] < 127 {
			out = append(out, buf[i])
			i++
			continue
		}
		r, n := utf8.DecodeRune(buf[i:])
		x[0] = byte(r)
		i += n
		r, n = utf8.DecodeRune(buf[i:])
		x[1] = byte(r)
		i += n
		out = append(out, x[0], x[1])
	}
	// Append remaining buf[last+1:]
	if last < len(buf)-1 {
		out = append(out, buf[last+1:]...)
	}
	return out, nil
}
