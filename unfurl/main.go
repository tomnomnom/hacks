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

	var procFn urlProc
	switch mode {

	case "querykeys":
		procFn = queryKeys

	case "paths":
		procFn = paths

	case "domains":
		procFn = domains

	default:
		fmt.Fprintf(os.Stderr, "unknown mode: %s\n", mode)
		return
	}

	sc := bufio.NewScanner(os.Stdin)

	seen := make(map[string]bool)

	for sc.Scan() {
		u, err := url.Parse(sc.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse failure: %s\n", err)
			continue
		}

		// some urlProc functions return multiple things,
		// so it's just easier to always get a slice and
		// loop over it instead of having two kinds of
		// urlProc functions.
		for _, val := range procFn(u) {
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

func paths(u *url.URL) []string {
	return []string{u.EscapedPath()}
}

func domains(u *url.URL) []string {
	return []string{u.Hostname()}
}
