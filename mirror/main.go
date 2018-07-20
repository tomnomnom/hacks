package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

func main() {
	flag.Parse()

	rawURL := flag.Arg(0)
	if rawURL == "" {
		fmt.Println("usage: mirror <url>")
		return
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := http.Get(rawURL)
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
				fmt.Printf("query string key '%s' with value '%s' reflected in response body (...%s...)\n", k, v, m[0])
			}
		}

	}
}
