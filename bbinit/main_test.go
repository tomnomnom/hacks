package main

import "testing"

func TestAsset(t *testing.T) {

	cases := []struct {
		a          asset
		hasScheme  bool
		domain     string
		isWildcard bool
	}{
		{a: asset{Type: "URL", Identifier: "Http://example.com"}, hasScheme: true, domain: "example.com", isWildcard: false},
		{a: asset{Type: "URL", Identifier: "Https://example.com"}, hasScheme: true, domain: "example.com", isWildcard: false},
		{a: asset{Type: "URL", Identifier: "http://example.com"}, hasScheme: true, domain: "example.com", isWildcard: false},
		{a: asset{Type: "URL", Identifier: "https://example.com"}, hasScheme: true, domain: "example.com", isWildcard: false},

		{a: asset{Type: "URL", Identifier: "*example.com"}, hasScheme: false, domain: "example.com", isWildcard: true},
		{a: asset{Type: "URL", Identifier: "*.example.com"}, hasScheme: false, domain: "example.com", isWildcard: true},
		{a: asset{Type: "URL", Identifier: ".example.com"}, hasScheme: false, domain: "example.com", isWildcard: true},
		{a: asset{Type: "URL", Identifier: "%.example.com"}, hasScheme: false, domain: "example.com", isWildcard: true},
		{a: asset{Type: "URL", Identifier: "%example.com"}, hasScheme: false, domain: "example.com", isWildcard: true},
	}

	for _, c := range cases {

		if c.a.hasScheme() != c.hasScheme {
			t.Errorf("want %t for a.hasScheme() with ident %s; don't", c.hasScheme, c.a.Identifier)
		}

		actual, err := c.a.Domain()
		if err != nil {
			t.Errorf("want nil error for a.Domain() with ident %s; don't", c.a.Identifier)
		}
		if actual != c.domain {
			t.Errorf("want %s for a.Domain() with ident %s; have %s", c.domain, c.a.Identifier, actual)
		}

		if c.a.isWildcard() != c.isWildcard {
			t.Errorf("want %t for a.isWildcard() with ident %s; don't", c.isWildcard, c.a.Identifier)
		}
	}
}
