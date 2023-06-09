package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/golemon/parse"
)

func usage() {
	fmt.Println("usage: lemon infile [outfile]")
	os.Exit(1)
}

func fileNameWithoutExtension(filename string) string {
	base := filepath.Base(filename)
	extension := filepath.Ext(base)

	return base[0 : len(base)-len(extension)]
}

// TODO: parse command flag
func main() {
	if len(os.Args) < 2 {
		usage()
	}

	infile := os.Args[1]
	outfile := fileNameWithoutExtension(infile) + ".go"

	if len(os.Args) == 3 {
		outfile = os.Args[2]
	}

	lemon := parse.NewLemon(infile, outfile)
	lemon.Parse()

}
