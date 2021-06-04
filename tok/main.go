package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"unicode"
)

func main() {
	var minlength int
	flag.IntVar(&minlength, "min", 1, "min length of string to be output")
	var maxlength int
	flag.IntVar(&maxlength, "max", 25, "max length of string to be output")
	var alphaNumOnly bool
	flag.BoolVar(&alphaNumOnly, "alpha-num-only", false, "return only strings containing at least one letter and one number")

	var delimExceptions string
	flag.StringVar(&delimExceptions, "delim-exceptions", "", "don't use the characters provided as delimiters")

	flag.Parse()

	r := bufio.NewReader(os.Stdin)
	var out strings.Builder

	maybeURLEncoded := false
	includesLetters := false
	includesNumbers := false
	last := ""

	reset := func() {
		maybeURLEncoded = false
		includesLetters = false
		includesNumbers = false
		out.Reset()
	}

	for {
		r, _, err := r.ReadRune()
		if err != nil {
			break
		}

		l := unicode.In(r, unicode.L)
		if l {
			includesLetters = true
		}

		n := unicode.In(r, unicode.N)
		if n {
			includesNumbers = true
		}

		if !l && !n && !isDelimException(r, delimExceptions) {
			if out.Len() == 0 {
				continue
			}

			str := out.String()

			if out.Len() < minlength {
				reset()
				continue
			}
			if out.Len() > maxlength {
				reset()
				continue
			}

			if alphaNumOnly && (!includesLetters || !includesNumbers || str == last) {
				reset()
				continue
			}

			if maybeURLEncoded {
				dec, err := url.QueryUnescape(str)
				if err == nil {
					str = dec
				}
			}

			fmt.Println(str)
			last = str

			reset()
			continue
		}

		if r == '%' {
			maybeURLEncoded = true
		}

		out.WriteRune(r)
	}

}

func isDelimException(r rune, delims string) bool {
	for _, comp := range delims {
		if r == comp {
			return true
		}
	}

	return false
}
