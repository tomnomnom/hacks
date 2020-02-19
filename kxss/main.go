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

	if len(reflected) == 0 {
		fmt.Println("no params were reflected; stopping")
		return
	}

	for _, r := range reflected {
		worked, err := checkAppend(inputURL, r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error requesting %s with modified %s param: %s", inputURL, r, err)
			continue
		}
		if worked {
			fmt.Printf("got reflection of param '%s' on %s", r, inputURL)
		}
	}
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

func checkAppend(targetURL, param string) (bool, error) {
	u, err := url.Parse(targetURL)
	if err != nil {
		return false, err
	}

	val := u.Query().Get(param)
	if val == "" {
		return false, fmt.Errorf("can't append to non-existant param %s", param)
	}

	u.Query().Set(param, val+"lol this should be a random string")

	reflected, err := checkReflected(u.String())
	if err != nil {
		return false, err
	}

	for _, r := range reflected {
		if r == param {
			return true, nil
		}
	}

	return false, nil
}
