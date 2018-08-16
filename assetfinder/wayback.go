package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func fetchWayback(domain string) ([]string, error) {

	res, err := http.Get(
		fmt.Sprintf("http://web.archive.org/cdx/search/cdx?url=*.%s/*&output=json&collapse=urlkey", domain),
	)
	if err != nil {
		return []string{}, err
	}

	raw, err := ioutil.ReadAll(res.Body)

	res.Body.Close()
	if err != nil {
		return []string{}, err
	}

	var wrapper [][]string
	err = json.Unmarshal(raw, &wrapper)

	out := make([]string, 0)

	skip := true
	for _, item := range wrapper {
		// The first item is always just the string "original",
		// so we should skip the first item
		if skip {
			skip = false
			continue
		}

		if len(item) < 3 {
			continue
		}

		u, err := url.Parse(item[2])
		if err != nil {
			continue
		}

		out = append(out, u.Hostname())
	}

	return out, nil
}
