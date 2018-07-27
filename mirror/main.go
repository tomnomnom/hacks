package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func main() {
	flag.Parse()

	var input io.Reader
	input = strings.NewReader(strings.Join(flag.Args(), "\n"))
	if flag.NArg() == 0 {
		input = os.Stdin
	}

	sc := bufio.NewScanner(input)

	for sc.Scan() {
		u, err := url.Parse(sc.Text())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		resp, err := http.Get(u.String())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		body := string(b)

		for k, vv := range u.Query() {
			for _, v := range vv {

				// short strings are so likely to show up in the response
				// that it's best just to skip over them to avoid too many
				// false positives. There should be a flag to control this.
				if len(v) < 4 {
					continue
				}

				// a fairly shonky way to get a few chars of context either side of the match
				// but it helps avoid trying to find the locations of all the matches in the
				// body, and then getting the context on either side, with all the bounds
				// checking etc that would need to be done for that.
				re, err := regexp.Compile("(.{0,6}" + regexp.QuoteMeta(v) + ".{0,6})")
				if err != nil {
					fmt.Fprintf(os.Stderr, "regexp compile error: %s", err)
				}

				matches := re.FindAllStringSubmatch(body, -1)

				for _, m := range matches {
					fmt.Printf("%s: '%s=%s' reflected in response body (...%s...)\n", u, k, v, m[0])
				}
			}

		}
	}
}
