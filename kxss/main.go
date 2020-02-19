package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

type appendCheck struct {
	url   string
	param string
}

func main() {

	sc := bufio.NewScanner(os.Stdin)

	initialChecks := make(chan string, 10)
	appendChecks := make(chan appendCheck, 10)

	// initial check worker pool
	var wgInitial sync.WaitGroup
	for i := 0; i < 10; i++ {
		wgInitial.Add(1)
		go func() {
			for inputURL := range initialChecks {

				reflected, err := checkReflected(inputURL)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error from checkReflected: %s\n", err)
					continue
				}

				if len(reflected) == 0 {
					fmt.Printf("no params were reflected in %s\n", inputURL)
					continue
				}

				for _, param := range reflected {
					appendChecks <- appendCheck{inputURL, param}
				}

			}
			wgInitial.Done()
		}()
	}

	// append check worker pool
	var wgAppend sync.WaitGroup
	for i := 0; i < 10; i++ {
		wgAppend.Add(1)

		go func() {
			for c := range appendChecks {
				wasReflected, err := checkAppend(c.url, c.param)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error from checkAppend for url %s with param %s: %s", c.url, c.param, err)
					continue
				}

				if wasReflected {
					fmt.Printf("got reflection of appended param %s on %s\n", c.param, c.url)
				}
			}
			wgAppend.Done()
		}()

	}

	for sc.Scan() {
		initialChecks <- sc.Text()
	}

	close(initialChecks)
	wgInitial.Wait()
	close(appendChecks)
	wgAppend.Wait()

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

	qs := u.Query()
	val := qs.Get(param)
	if val == "" {
		return false, fmt.Errorf("can't append to non-existant param %s", param)
	}

	qs.Set(param, val+"lol this should be a random string")
	u.RawQuery = qs.Encode()

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
