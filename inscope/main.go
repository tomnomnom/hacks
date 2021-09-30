package main

import (
	"bufio"
	"flag"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type scopeChecker struct {
	patterns     []*regexp.Regexp
	antipatterns []*regexp.Regexp
}

func init() {
	flag.Usage = func() {
		h := []string{
			"Filters in scope and out of scope urls from stdin.",
			"",
			"Options:",
			"  -v, --inverse         Prints out of scope items",
			"",
		}

		fmt.Fprintf(os.Stderr, strings.Join(h, "\n"))
	}
}

func (s *scopeChecker) inScope(domain string) bool {

	// if it's a URL pull the hostname out to avoid matching
	// on part of the path or something like that
	if isURL(domain) {
		var err error
		domain, err = getHostname(domain)
		if err != nil {
			return false
		}
	}

	inScope := false
	for _, p := range s.patterns {
		if p.MatchString(domain) {
			inScope = true
			break
		}
	}

	for _, p := range s.antipatterns {
		if p.MatchString(domain) {
			return false
		}
	}
	return inScope
}

func newScopeChecker(r io.Reader) (*scopeChecker, error) {
	sc := bufio.NewScanner(r)
	s := &scopeChecker{
		patterns: make([]*regexp.Regexp, 0),
	}

	for sc.Scan() {
		p := strings.TrimSpace(sc.Text())
		if p == "" {
			continue
		}

		isAnti := false
		if p[0] == '!' {
			isAnti = true
			p = p[1:]
		}

		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}

		if isAnti {
			s.antipatterns = append(s.antipatterns, re)
		} else {
			s.patterns = append(s.patterns, re)
		}
	}

	return s, nil
}

func main() {
	var inverse bool
	flag.BoolVar(&inverse, "inverse", false, "")
	flag.BoolVar(&inverse, "v", false, "")

	flag.Parse()

	sf, err := openScopefile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening scope file: %s\n", err)
		return
	}

	checker, err := newScopeChecker(sf)
	sf.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing scope file: %s\n", err)
		return
	}

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		domain := strings.TrimSpace(sc.Text())

		inScope := checker.inScope(domain)
		if !inverse && inScope {
			fmt.Println(domain)
			continue
		}
		
		if inverse && !inScope {
			fmt.Println(domain)
		}
	}
}

func getHostname(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	return u.Hostname(), nil
}

func isURL(s string) bool {
	s = strings.TrimSpace(strings.ToLower(s))

	if len(s) < 6 {
		return false
	}

	return s[:5] == "http:" || s[:6] == "https:"
}

func openScopefile() (io.ReadCloser, error) {
	pwd, err := filepath.Abs(".")
	if err != nil {
		return nil, err
	}

	for {
		f, err := os.Open(filepath.Join(pwd, ".scope"))

		// found one!
		if err == nil {
			return f, nil
		}

		newPwd := filepath.Dir(pwd)
		if newPwd == pwd {
			break
		}
		pwd = newPwd
	}

	return nil, errors.New("unable to find .scope file in current directory or any parent directory")
}
