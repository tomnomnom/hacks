package main

import (
	"flag"
	"fmt"
	"net/url"
)

func main() {
	flag.Parse()

	char := flag.Arg(0)
	if len(char) != 1 {
		fmt.Println("usage: unisub <char>")
		return
	}

	subs, ok := translations[rune(char[0])]
	if !ok {
		fmt.Println("no substitutions found")
		return
	}

	for _, s := range subs {
		fmt.Printf("%c %U %s\n", s, s, url.QueryEscape(string(s)))
	}
}
