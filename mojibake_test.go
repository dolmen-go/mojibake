package mojibake_test

import (
	"testing"
	"unicode/utf8"

	"github.com/dolmen-go/mojibake"
)

// "magnï¿½tique"

func TestFixDoubleUTF8(t *testing.T) {
	for _, test := range []struct {
		in  string
		out string
	}{
		{"", ""},
		{"a", "a"},
		{"gÃ©lule", "gélule"},
		{"RÃ©fÃ©rence", "Référence"},
		{"CrÃ©Ã©", "Créé"},
		{"CrÃ©Ã©e", "Créée"},
		{"Créé", "Créé"},
		{"Créée", "Créée"},
		{"prÃ¨s", "près"},
		{"chÃ¢teau", "château"},
		{"Ã©nergie", "énergie"},
		{"Ã©vÃ¨nement", "évènement"},
		{"TÃªte", "Tête"},
		{"\u00f0\u009f\u0087\u00ab\u00f0\u009f\u0087\u00b7", "🇫🇷"},
	} {
		work := []byte(test.in)
		work, err := mojibake.FixDoubleUTF8(work)
		if err != nil && err != mojibake.ErrImpure {
			t.Errorf("%q: unexpected error %q", test.in, err)
		} else {
			if err == mojibake.ErrImpure {
				t.Logf("%q: %q", test.in, err)
			}
			out := string(work)
			if out != test.out {
				t.Errorf("%q: got %q, want %q", test.in, out, test.out)
				//t.Logf("%q -> %q", test.out, out)
			}
		}
	}
}

func TestDoubleUTF8Canaries(t *testing.T) {
	var buf [utf8.UTFMax]byte
	startBytes := make(map[byte]bool, 128)
	for r := rune(utf8.RuneSelf); r <= utf8.MaxRune; r++ {
		if !utf8.ValidRune(r) {
			continue
		}
		_ = utf8.EncodeRune(buf[:], r)
		if startBytes[buf[0]] {
			continue
		}
		b := buf[0]
		if b < mojibake.UTF8FirstByteMin {
			t.Errorf("%d %c: %d < %d", r, r, b, mojibake.UTF8FirstByteMin)
		}
		if b > mojibake.UTF8FirstByteMax {
			t.Errorf("%d %c: %d > %d", r, r, b, mojibake.UTF8FirstByteMax)
		}
		startBytes[b] = true
		n := utf8.EncodeRune(buf[:], rune(b))
		t.Logf("%d %c => %s %d", r, r, buf[:n], b)
	}
	t.Log(len(startBytes), "values")
}
