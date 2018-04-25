package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
)

type urlCheck func(*url.URL) bool

func main() {

	checks := []urlCheck{
		// query string stuff
		func(u *url.URL) bool {

			interesting := 0
			for k, vv := range u.Query() {
				for _, v := range vv {
					if qsCheck(k, v) {
						interesting++
					}
				}
			}
			return interesting > 0
		},

		// extensions
		func(u *url.URL) bool {
			return strings.HasSuffix(u.EscapedPath(), ".php") ||
				strings.HasSuffix(u.EscapedPath(), ".asp") ||
				strings.HasSuffix(u.EscapedPath(), ".aspx")
		},

		// path bits
		func(u *url.URL) bool {
			return strings.Contains(u.EscapedPath(), "ajax") ||
				strings.Contains(u.EscapedPath(), "jsonp")
		},

		// non-standard port
		func(u *url.URL) bool {
			return (u.Port() != "80" && u.Port() != "443" && u.Port() != "")
		},
	}

	seen := make(map[string]bool)

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {

		u, err := url.Parse(sc.Text())
		if err != nil {
			//fmt.Fprintf(os.Stderr, "failed to parse url %s [%s]\n", sc.Text(), err)
			continue
		}
		if isStaticFile(u) {
			continue
		}

		// Go's maps aren't ordered, but we want to use all the param names
		// as part of the key to output only unique requests. To do that, put
		// them into a slice and then sort it.
		pp := make([]string, 0)
		for p, _ := range u.Query() {
			pp = append(pp, p)
		}
		sort.Strings(pp)

		key := fmt.Sprintf("%s%s?%s", u.Hostname(), u.EscapedPath(), strings.Join(pp, "&"))

		// Only output each host + path + params combination once
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = true

		interesting := 0

		for _, check := range checks {
			if check(u) {
				interesting++
			}
		}

		if interesting > 0 {
			fmt.Println(sc.Text())
		}

	}

}

func qsCheck(k, v string) bool {
	// the super-common utm_referrer etc
	// are rarely interesting
	if strings.HasPrefix(k, "utm_") {
		return false
	}

	return strings.HasPrefix(v, "http") ||
		strings.Contains(v, "/") ||
		strings.Contains(k, "redirect") ||
		strings.Contains(k, "debug") ||
		strings.Contains(k, "callback")
}

func isStaticFile(u *url.URL) bool {
	exts := []string{
		".html",
		".htm",
		".svg",
		".eot",
		".ttf",
		".js",
		".png",
		".jpg",
		".jpeg",
		".gif",
	}

	for _, e := range exts {
		if strings.HasSuffix(u.EscapedPath(), e) {
			return true
		}
	}

	return false
}
