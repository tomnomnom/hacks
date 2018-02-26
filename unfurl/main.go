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

	var unique bool
	flag.BoolVar(&unique, "u", false, "")
	flag.BoolVar(&unique, "unique", false, "")

	var verbose bool
	flag.BoolVar(&verbose, "v", false, "")
	flag.BoolVar(&verbose, "verbose", false, "")

	flag.Parse()

	mode := flag.Arg(0)
	fmtStr := flag.Arg(1)

	procFn, ok := map[string]urlProc{
		"keys":    keys,
		"values":  values,
		"domains": domains,
		"paths":   paths,
		"format":  format,
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
			if verbose {
				fmt.Fprintf(os.Stderr, "parse failure: %s\n", err)
			}
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

			if seen[val] && unique {
				continue
			}

			fmt.Println(val)

			// no point using up memory if we're outputting dupes
			if unique {
				seen[val] = true
			}
		}
	}

	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %s\n", err)
	}
}

type urlProc func(*url.URL, string) []string

func keys(u *url.URL, _ string) []string {
	out := make([]string, 0)
	for key, _ := range u.Query() {
		out = append(out, key)
	}
	return out
}

func values(u *url.URL, _ string) []string {
	out := make([]string, 0)
	for _, vals := range u.Query() {
		for _, val := range vals {
			out = append(out, val)
		}
	}
	return out
}

func domains(u *url.URL, f string) []string {
	return format(u, "%d")
}

func paths(u *url.URL, f string) []string {
	return format(u, "%p")
}

func format(u *url.URL, f string) []string {
	out := &bytes.Buffer{}

	inFormat := false
	for _, r := range f {

		if r == '%' && !inFormat {
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

func init() {
	flag.Usage = func() {
		h := "Format URLs provided on stdin\n\n"

		h += "Usage:\n"
		h += "  unfurl [OPTIONS] [MODE] [FORMATSTRING]\n\n"

		h += "Options:\n"
		h += "  -u, --unique   Only output unique values\n"
		h += "  -v, --verbose  Verbose mode (output URL parse errors)\n\n"

		h += "Modes:\n"
		h += "  keys     Keys from the query string (one per line)\n"
		h += "  values   Values from the query string (one per line)\n"
		h += "  domains  The hostname (e.g. sub.example.com)\n"
		h += "  paths    The request path (e.g. /users)\n"
		h += "  format   Specify a custom format (see below)\n\n"

		h += "Format Directives:\n"
		h += "  %%  A literal percent character\n"
		h += "  %s  The request scheme (e.g. https)\n"
		h += "  %d  The domain (e.g. sub.example.com)\n"
		h += "  %P  The port (e.g. 8080)\n"
		h += "  %p  The path (e.g. /users)\n"
		h += "  %q  The raw query string (e.g. a=1&b=2)\n"
		h += "  %f  The page fragment (e.g. page-section)\n\n"

		h += "Examples:\n"
		h += "  cat urls.txt | unfurl keys\n"
		h += "  cat urls.txt | unfurl format %s://%d%p?%q\n"

		fmt.Fprint(os.Stderr, h)
	}
}
