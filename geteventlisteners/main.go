package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/chromedp/chromedp"
)

func main() {
	flag.Parse()

	// default to stdin unless we have an arg to use
	var input io.Reader
	input = os.Stdin
	if flag.Arg(0) != "" {
		input = strings.NewReader(flag.Arg(0))
	}
	sc := bufio.NewScanner(input)

	parent, pcancel := chromedp.NewContext(context.Background())
	defer pcancel()

	for sc.Scan() {
		ctx, cancel := chromedp.NewContext(parent)
		requestURL := sc.Text()

		var res map[string][]string

		err := chromedp.Run(ctx,
			chromedp.Navigate(requestURL),
			chromedp.EvaluateAsDevTools(`
			var listeners = getEventListeners(window)

			for (let i in listeners){
				listeners[i] = listeners[i].map(l => {
					return l.listener.toString()
				})
			}

			listeners`,
				&res),
		)

		if err != nil {
			cancel()
			fmt.Fprintf(os.Stderr, "%s\n", err)
			continue
		}

		// TODO: provide option to write one output file per request
		fmt.Printf("// %s\n", requestURL)
		fmt.Println("(function(){")
		for event, listeners := range res {
			seen := make(map[string]bool)

			for i, l := range listeners {
				if seen[l] {
					continue
				}
				seen[l] = true

				suffix := strconv.Itoa(i + 1)
				if suffix == "1" {
					suffix = ""
				}

				fmt.Printf("let on%s%s = %s\n\n", event, suffix, l)
			}
		}
		fmt.Println("})()\n")

		cancel()
	}

}
