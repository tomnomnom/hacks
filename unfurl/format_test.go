package main

import (
	"net/url"
	"testing"
)

func TestFormat(t *testing.T) {
	cases := []struct {
		url      string
		format   string
		expected string
	}{
		{"https://example.com/foo", "%d", "example.com"},
		{"https://example.com/foo", "%d%p", "example.com/foo"},
		{"https://example.com/foo", "%s://%d%p", "https://example.com/foo"},

		{"https://example.com:8080/foo", "%d", "example.com"},
		{"https://example.com:8080/foo", "%P", "8080"},

		{"https://example.com/foo?a=b&c=d", "%p", "/foo"},
		{"https://example.com/foo?a=b&c=d", "%q", "a=b&c=d"},

		{"https://example.com/foo#bar", "%f", "bar"},
		{"https://example.com#bar", "%f", "bar"},

		{"https://example.com#bar", "foo%%bar", "foo%bar"},
		{"https://example.com#bar", "%s://%%", "https://%"},
	}

	for _, c := range cases {
		u, err := url.Parse(c.url)
		if err != nil {
			t.Fatal(err)
		}

		actual := format(u, c.format)

		if actual[0] != c.expected {
			t.Errorf("want %s for format(%s, %s); have %s", c.expected, c.url, c.format, actual)
		}
	}
}
