package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/ditashi/jsbeautifier-go/jsbeautifier"
)

type filterArgs []string

func (f *filterArgs) Set(val string) error {
	*f = append(*f, val)
	return nil
}

func (f filterArgs) String() string {
	return "string"
}

func (f filterArgs) Includes(search string) bool {
	search = strings.ToLower(search)
	for _, filter := range f {
		filter = strings.ToLower(filter)
		if filter == search {
			return true
		}
	}
	return false
}

func main() {
	var filters filterArgs
	flag.Var(&filters, "filter", "")
	flag.Var(&filters, "f", "")

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

		buf := &strings.Builder{}
		first := true
		for event, listeners := range res {

			if len(filters) > 0 && !filters.Includes(event) {
				continue
			}

			if first {
				fmt.Fprintf(buf, "// %s\n", requestURL)
				buf.WriteString("(function(){")
				first = false
			}

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

				fmt.Fprintf(buf, "    let on%s%s = %s\n\n", event, suffix, l)
			}
		}

		if first {
			// we didn't find any matching event listeners
			cancel()
			continue
		}

		buf.WriteString("})()")

		raw := buf.String()
		options := jsbeautifier.DefaultOptions()
		out, err := jsbeautifier.Beautify(&raw, options)

		if err != nil {
			out = raw
		}
		fmt.Println(requestURL)

		// TODO: organise files into one dir per domain
		fn := genFilename(requestURL)
		f, err := os.Create(fn)
		fmt.Fprintf(f, "%s\n", out)
		f.Close()

		cancel()
	}

}

func genFilename(u string) string {
	re := regexp.MustCompile("[^a-zA-Z0-9_.-]")
	fn := re.ReplaceAllString(u, "-")

	// remove multiple dashes in a row
	re = regexp.MustCompile("-+")
	fn = re.ReplaceAllString(fn, "-")

	return fmt.Sprintf("%s.js", fn)
}
