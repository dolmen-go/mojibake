//+build gofuzz

// Build: go-fuzz-build -tags gofuzz github.com/dolmen-go/mojibake
// Run:   go-fuzz -bin mojibake-fuzz.zip -workdir testdata

package mojibake

func Fuzz(data []byte) int {
	FixDoubleUTF8(data)
	return 0
}
