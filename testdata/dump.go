//+build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/text/unicode/runenames"
)

func _main(path string) int {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", path, err)
		return 1
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", path, err)
		return 1
	}
	for i, r := range string(buf) {
		fmt.Printf("%04d: %s\n", i, runenames.Name(r))
	}
	fmt.Printf("%04d:\n", len(buf))
	return 0
}

func main() {
	os.Exit(_main(os.Args[1]))
}
