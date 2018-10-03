package main

import (
	"strings"
	"testing"
)

func TestSmoke(t *testing.T) {
	sf := strings.NewReader(`
		.*\.example\.com$
		^example\.com$
		.*\.example\.net$
		!.*outofscope\.example\.net$
	`)

	checker, err := newScopeChecker(sf)
	if err != nil {
		t.Fatalf("failed to make scope checker: %s", err)
	}

	cases := []struct {
		url     string
		inScope bool
	}{
		{"https://example.com/footle", true},
		{"https://inscope.example.com/some/path?foo=bar", true},
		{"example.com", true},
		{"http://sub.example.com", true},
		{"https://outofscope.example.net/bar", false},
		{"example.net", false},
	}

	for _, c := range cases {
		expected := c.inScope
		actual := checker.inScope(c.url)

		if actual != expected {
			t.Errorf("want %t for inScope(%s), have %t", expected, c.url, actual)
			t.Logf("%#v", checker)
		}
	}
}

func TestIsURL(t *testing.T) {
	if isURL("http://example.com/footle") == false {
		t.Errorf("http://example.com/footle should be a URL but isn't")
	}

	if isURL(" https://example.com/footle") == false {
		t.Errorf("https://example.com/footle should be a URL but isn't")
	}
}
