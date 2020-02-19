package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSimple(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		qs := r.URL.Query()
		fmt.Fprintf(w, "hello, %s", qs.Get("name"))
	}))

	defer ts.Close()

	r, err := checkReflected(ts.URL + "?name=tom")
	t.Logf("params reflected: %#v", r)

	if err != nil {
		t.Fatalf("expected nil error from checkReflected(), have %s", err)
	}

	if len(r) != 1 {
		t.Errorf("wanted length 1 for returned keys, have %d", len(r))
	}
}
