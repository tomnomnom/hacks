package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"net/url"
	"os"
)

func main() {
	flag.Parse()

	mode := flag.Arg(0)
	fmtStr := flag.Arg(1)

	procFn, ok := map[string]urlProc{
		"keys":   queryKeys,
		"values": queryValues,
		"format": format,
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
		for _, val := range procFn(u, fmtStr) {

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

type urlProc func(*url.URL, string) []string

func queryKeys(u *url.URL, _ string) []string {
	out := make([]string, 0)
	for key, _ := range u.Query() {
		out = append(out, key)
	}
	return out
}

func queryValues(u *url.URL, _ string) []string {
	out := make([]string, 0)
	for _, vals := range u.Query() {
		for _, val := range vals {
			out = append(out, val)
		}
	}
	return out
}

func format(u *url.URL, f string) []string {
	out := &bytes.Buffer{}

	inFormat := false
	for _, r := range f {

		if r == '%' {
			inFormat = true
			continue
		}

		if !inFormat {
			out.WriteRune(r)
			continue
		}

		switch r {
		case '%':
			out.WriteRune('%')
		case 's':
			out.WriteString(u.Scheme)
		case 'd':
			out.WriteString(u.Hostname())
		case 'P':
			out.WriteString(u.Port())
		case 'p':
			out.WriteString(u.EscapedPath())
		case 'q':
			out.WriteString(u.RawQuery)
		case 'f':
			out.WriteString(u.Fragment)
		default:
			// output untouched
			out.WriteRune('%')
			out.WriteRune(r)
		}

		inFormat = false
	}

	return []string{out.String()}
}
