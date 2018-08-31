package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {

	z := html.NewTokenizer(os.Stdin)

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}

		t := z.Token()

		if t.Type == html.CommentToken {
			d := strings.Replace(t.Data, "\n", " ", -1)
			d = strings.TrimSpace(d)
			if d == "" {
				continue
			}
			fmt.Println(d)
		}

	}

}
