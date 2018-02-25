package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
)

func main() {
	flag.Parse()

	mode := flag.Arg(0)

	procFn, ok := map[string]urlProc{
		"querykeys":   queryKeys,
		"queryvalues": queryValues,
		"paths":       paths,
		"domains":     domains,
	}[mode]

	if !ok {
		fmt.Fprintf(os.Stderr, "unknown mode: %s\n", mode)
		return
	}

	sc := bufio.NewScanner(os.Stdin)

	seen := make(map[string]bool)

	for sc.Scan() {
		u, err := url.Parse(sc.Text())
		if err != nil {
			// TODO: add this back in with a verbose flag
			// fmt.Fprintf(os.Stderr, "parse failure: %s\n", err)
			continue
		}

		// some urlProc functions return multiple things,
		// so it's just easier to always get a slice and
		// loop over it instead of having two kinds of
		// urlProc functions.
		for _, val := range procFn(u) {

			// you do see empty values sometimes
			if val == "" {
				continue
			}

			// TODO: add a mode that outputs duplicates
			if seen[val] {
				continue
			}

			fmt.Println(val)
			seen[val] = true
		}
	}

	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %s\n", err)
	}
}

type urlProc func(*url.URL) []string

func queryKeys(u *url.URL) []string {
	out := make([]string, 0)
	for key, _ := range u.Query() {
		out = append(out, key)
	}
	return out
}

func queryValues(u *url.URL) []string {
	out := make([]string, 0)
	for _, vals := range u.Query() {
		for _, val := range vals {
			out = append(out, val)
		}
	}
	return out
}

func paths(u *url.URL) []string {
	return []string{u.EscapedPath()}
}

func domains(u *url.URL) []string {
	return []string{u.Hostname()}
}
