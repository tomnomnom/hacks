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
	if flag.Arg(0) == "" {
		input = os.Stdin
	}

	sc := bufio.NewScanner(input)

	for sc.Scan() {
		u, err := url.Parse(sc.Text())
		if err != nil {
			fmt.Println(err)
			return
		}

		resp, err := http.Get(u.String())
		if err != nil {
			fmt.Println(err)
			return
		}

		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		body := string(b)

		for k, vv := range u.Query() {
			for _, v := range vv {

				// a fairly shonky way to get a few chars of context either side of the match
				re, err := regexp.Compile("(.{0,6}" + regexp.QuoteMeta(v) + ".{0,6})")
				if err != nil {
					fmt.Fprintf(os.Stderr, "regexp compile error: %s", err)
				}

				matches := re.FindAllStringSubmatch(body, -1)

				for _, m := range matches {
					fmt.Printf("%s: query string key '%s' with value '%s' reflected in response body (...%s...)\n", u, k, v, m[0])
				}
			}

		}
	}
}
