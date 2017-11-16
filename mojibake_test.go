package mojibake_test

import (
	"testing"

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
