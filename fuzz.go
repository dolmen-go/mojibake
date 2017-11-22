//+build gofuzz

// Build: go-fuzz-build -tags gofuzz github.com/dolmen-go/mojibake
// Run:   go-fuzz -bin mojibake-fuzz.zip -workdir testdata

package mojibake

import (
	"fmt"
	"unicode/utf8"
)

func Fuzz(data []byte) int {
	if len(data) > 20 {
		return -1
	}

	isValidUTF8 := utf8.Valid(data)

	out, err := FixDoubleUTF8(data)

	if (err != ErrUTF8) != isValidUTF8 {
		panic(fmt.Sprintf("UTF-8 validation error: %s", err))
	}

	switch {
	case len(data) > 17:
		return -1
	case err == nil && len(out) == len(data):
		return 0
	case len(data) < 8: // priority to short strings
		return 1
	default:
		return 0
	}
}
