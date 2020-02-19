package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	flag.Parse()

	inputURL := flag.Arg(0)

	reflected, err := checkReflected(inputURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error from checkReflected: %s\n", err)
	}

	fmt.Printf("reflected: %#v\n", reflected)
}

func checkReflected(targetURL string) ([]string, error) {

	out := make([]string, 0)

	resp, err := http.Get(targetURL)
	if err != nil {
		return out, err
	}
	if resp.Body == nil {
		return out, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return out, err
	}

	body := string(b)

	u, err := url.Parse(targetURL)
	if err != nil {
		return out, err
	}

	for key, vv := range u.Query() {
		for _, v := range vv {
			if !strings.Contains(body, v) {
				continue
			}

			out = append(out, key)
		}
	}

	return out, nil
}
