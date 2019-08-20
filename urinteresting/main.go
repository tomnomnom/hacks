package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
)

// Ideas:
//   More than, say, 3 query string parameteres (exluding utm_*?)
//   Popular app names (phpmyadmin etc) in path
//	 Filenames from configfiles list / seclist
//   dev/stage/test in path or hostname
//   jenkins, graphite etc in hostname or path

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
			exts := []string{
				".php",
				".phtml",
				".asp",
				".aspx",
				".asmx",
				".ashx",
				".cgi",
				".pl",
				".json",
				".xml",
				".rb",
				".py",
				".sh",
				".yaml",
				".yml",
				".toml",
				".ini",
				".md",
				".mkd",
				".do",
				".jsp",
				".jspa",
			}

			p := strings.ToLower(u.EscapedPath())
			for _, e := range exts {
				if strings.HasSuffix(p, e) {
					return true
				}
			}

			return false
		},

		// path bits
		func(u *url.URL) bool {
			p := strings.ToLower(u.EscapedPath())
			return strings.Contains(p, "ajax") ||
				strings.Contains(p, "jsonp") ||
				strings.Contains(p, "admin") ||
				strings.Contains(p, "include") ||
				strings.Contains(p, "src") ||
				strings.Contains(p, "redirect") ||
				strings.Contains(p, "proxy") ||
				strings.Contains(p, "test") ||
				strings.Contains(p, "tmp") ||
				strings.Contains(p, "temp")
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

		if isBoringStaticFile(u) {
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

// qsCheck looks a key=value pair from a query
// string and returns true if it looks interesting
func qsCheck(k, v string) bool {
	k = strings.ToLower(k)
	v = strings.ToLower(v)

	// the super-common utm_referrer etc
	// are rarely interesting
	if strings.HasPrefix(k, "utm_") {
		return false
	}

	// value checks
	return strings.HasPrefix(v, "http") ||
		strings.Contains(v, "{") ||
		strings.Contains(v, "[") ||
		strings.Contains(v, "/") ||
		strings.Contains(v, "\\") ||
		strings.Contains(v, "<") ||
		strings.Contains(v, "(") ||
		// shoutout to liveoverflow ;)
		strings.Contains(v, "eyj") ||

		// key checks
		strings.Contains(k, "redirect") ||
		strings.Contains(k, "debug") ||
		strings.Contains(k, "password") ||
		strings.Contains(k, "passwd") ||
		strings.Contains(k, "file") ||
		strings.Contains(k, "fn") ||
		strings.Contains(k, "template") ||
		strings.Contains(k, "include") ||
		strings.Contains(k, "require") ||
		strings.Contains(k, "url") ||
		strings.Contains(k, "uri") ||
		strings.Contains(k, "src") ||
		strings.Contains(k, "href") ||
		strings.Contains(k, "func") ||
		strings.Contains(k, "callback")
}

func isBoringStaticFile(u *url.URL) bool {
	exts := []string{
		// OK, so JS could be interesting, but 99% of the time it's boring.
		".js",

		".html",
		".htm",
		".svg",
		".eot",
		".ttf",
		".woff",
		".woff2",
		".png",
		".jpg",
		".jpeg",
		".gif",
		".ico",
	}

	p := strings.ToLower(u.EscapedPath())
	for _, e := range exts {
		if strings.HasSuffix(p, e) {
			return true
		}
	}

	return false
}
