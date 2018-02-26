package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	flag.Parse()
	keys := flag.Args()

	z := html.NewTokenizer(os.Stdin)

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}

		t := z.Token()

		for _, a := range t.Attr {

			if a.Val == "" {
				continue
			}

			// default to all values
			if len(keys) == 0 {
				fmt.Println(a.Val)
				continue
			}

			// output values only for keys we want
			for _, k := range keys {
				if k == a.Key {
					fmt.Println(a.Val)
				}
			}
		}
	}

}
