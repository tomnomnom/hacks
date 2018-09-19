package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"unicode"
)

var ASCIIChar = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x0020, 0x007F, 1},
	},
}

func main() {
	flag.Parse()

	filename := flag.Arg(0)
	if filename == "" {
		fmt.Fprintln(os.Stderr, "usage: b64d <filename>")
		return
	}

	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	content := string(b)

	// TODO: deal with urlencoded bits
	re := regexp.MustCompile("[^A-Za-z0-9+/][a-zA-Z0-9+/]+={0,2}")
	matches := re.FindAllString(content, -1)

	if matches == nil {
		return
	}

	for _, m := range matches {

		if len(m) < 7 {
			continue
		}

		// match has one extra char at the beginning
		if (len(m)-1)%4 != 0 {
			continue
		}
		decb, _ := base64.StdEncoding.DecodeString(m[1:])
		decoded := string(decb)
		if decoded == "" {
			continue
		}

		decoded = strings.Replace(decoded, "\n", " ", -1)

		containsNonASCII := false
		for _, r := range decoded {
			if !unicode.Is(ASCIIChar, r) {
				containsNonASCII = true
				break
			}
		}
		if containsNonASCII {
			continue
		}

		fmt.Println(decoded)
	}
}
