//+build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dolmen-go/mojibake"
)

func _main() int {
	dir, err := os.Open("corpus")
	if err != nil {
		fmt.Fprintf(os.Stderr, "corpus: %s\n", err)
		return 1
	}
	corpus, err := dir.Readdir(-1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "corpus: %s\n", err)
		dir.Close()
		return 1
	}
	dir.Close()

	const (
		extDouble  = ".double"
		extImpure  = ".impure"
		extNotUTF8 = ".broken"
		extSimple  = ".simple"
	)

	for _, file := range corpus {
		if !file.Mode().IsRegular() {
			continue
		}
		// File already has an extension
		if strings.Contains(file.Name(), ".") {
			continue
		}

		path := "corpus/" + file.Name()
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", path, err)
			continue
		}
		var ext string
		if func() bool {
			defer f.Close()
			defer func() {
				// Catch any exceptions
				if e := recover(); e != nil {
					fmt.Fprintf(os.Stderr, "%s: %v\n", path, err)
				}
			}()

			content, err := ioutil.ReadAll(f)
			c, err := mojibake.FixDoubleUTF8(content)
			switch {
			case err == mojibake.ErrImpure:
				ext = extImpure
			case err == mojibake.ErrUTF8:
				ext = extNotUTF8
			case len(c) < len(content):
				ext = extDouble
			default:
				ext = extSimple
			}
			return true
		}() {
			if err := os.Rename(path, path+ext); err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", path, err)
			}
		}
	}
	return 0
}

func main() {
	os.Exit(_main())
}
