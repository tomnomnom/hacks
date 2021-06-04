package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSimple(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		qs := r.URL.Query()
		fmt.Fprintf(w, "hello, %s", qs.Get("name"))
	}))

	defer ts.Close()

	r, err := checkReflected(ts.URL + "?name=Mr%20Naughty")
	t.Logf("params reflected: %#v", r)

	if err != nil {
		t.Fatalf("expected nil error from checkReflected(), have %s", err)
	}

	if len(r) != 1 {
		t.Errorf("wanted length 1 for returned keys, have %d", len(r))
	}
}

func TestAppend(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		qs := r.URL.Query()
		fmt.Fprintf(w, "hello, %s", qs.Get("name"))
	}))

	defer ts.Close()

	r, err := checkAppend(ts.URL+"?name=Mr%20Naughty", "name", "somerandomvalue")

	if err != nil {
		t.Fatalf("expected nil error from checkAppend(), have %s", err)
	}

	if !r {
		t.Errorf("wanted checkAppend() to return true, but it didn't")
	}
}
