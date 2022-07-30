package main

import (
	"bufio"
	"strings"
	"flag"
	"fmt"
	"os"
)

func main() {
	var flip bool
	flag.BoolVar(&flip, "f", false, "")
	flag.BoolVar(&flip, "flip", false, "")

	var trim bool
	flag.BoolVar(&trim, "t", false, "")
	flag.BoolVar(&trim, "trim", false, "")

	var separator string
	flag.StringVar(&separator, "s", "", "")
	flag.StringVar(&separator, "separator", "", "")

	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	prefixFile, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	suffixFile, err := os.Open(flag.Arg(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	// use 'a' and 'b' because which is the prefix
	// and which is the suffix depends on if we're in
	// flip mode or not.
	fileA := prefixFile
	fileB := suffixFile

	if flip {
		fileA, fileB = fileB, fileA
	}

	a := bufio.NewScanner(fileA)
	for a.Scan() {
		// rewind file B so we can scan it again
		fileB.Seek(0, 0)

		b := bufio.NewScanner(fileB)
		for b.Scan() {
			aText := a.Text()
			bText := b.Text()

			if trim {
				aText = strings.TrimSpace(aText)
				bText = strings.TrimSpace(bText)
			}

			if flip {
				fmt.Printf("%s%s%s\n", bText, separator, aText)
			} else {
				fmt.Printf("%s%s%s\n", aText, separator, bText)
			}
		}
	}

}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Combine the lines from two files in every combination\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  comb [OPTIONS] <prefixfile> <suffixfile>\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -f, --flip             Flip mode (order by suffix)\n")
		fmt.Fprintf(os.Stderr, "  -s, --separator <str>  String to place between prefix and suffix\n")
		fmt.Fprintf(os.Stderr, "  -t, --trim             Trim strings\n")
	}
}
