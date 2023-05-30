//go:build go1.18
// +build go1.18

package mojibake

import (
	"testing"
	"unicode/utf8"
)

func FuzzDoubleUTF8(f *testing.F) {

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 20 {
			t.Skip()
			return
		}

		isValidUTF8 := utf8.Valid(data)

		_, err := FixDoubleUTF8(data)

		if (err != ErrUTF8) != isValidUTF8 {
			t.Fatalf("UTF-8 validation error: %s", err)
		}
	})
}
