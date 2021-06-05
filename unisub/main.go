package main

import (
	"flag"
	"fmt"
	"net/url"
	"strings"
)

func main() {

	var silent bool
	flag.BoolVar(&silent, "silent", false, "enable silnet mode")
	var urlEncodeOnly bool
	flag.BoolVar(&urlEncodeOnly, "encode", false, "display only url econded string")

	flag.Parse()
	
	if flag.NArg() < 1 {
		fmt.Println("usage: unisub <char>")
		return
	}
	
	char := flag.Arg(0)
		subs, ok := translations[rune(char[0])]
	if !ok {
		fmt.Println("no substitutions found")
		return
	}

	for _, s := range subs {
		if silent {
			fmt.Printf("%c\n", s)
		} else if urlEncodeOnly {
			fmt.Printf("%s\n", url.QueryEscape(string(s)))
		} else {
			fmt.Printf("fallback: %c %U %s\n", s, s, url.QueryEscape(string(s)))
		}		
	}

	for cp := 1; cp < 0x10FFFF; cp++ {
		s := rune(cp)
		if char == string(s) {
			continue
		}

		if strings.ToLower(string(s)) == char {
			if silent {
				fmt.Printf("%c\n", s)
			} else if urlEncodeOnly {
				fmt.Printf("%s\n", url.QueryEscape(string(s)))
			} else {
				fmt.Printf("toLower: %c %U %s\n", s, s, url.QueryEscape(string(s)))
			}
		}

		if strings.ToUpper(string(s)) == char {
			if silent {
				fmt.Printf("%c\n", s)
			} else if urlEncodeOnly {
				fmt.Printf("%s\n", url.QueryEscape(string(s)))
			} else {
				fmt.Printf("toUpper: %c %U %s\n", s, s, url.QueryEscape(string(s)))
			}
		}
	}
}
