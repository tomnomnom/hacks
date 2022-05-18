package main

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aofei/mimesniffer"
	"github.com/nsf/jsondiff"
)

func main() {

	mimesniffer.Register("application/json", func(bs []byte) bool {
		return json.Valid(bs)
	})

	base, err := http.NewRequest("GET", "https://httpbin.org/anything", nil)
	fatalErr(err)

	candidate, err := http.NewRequest("GET", "https://httpbin.org/anything/foo", nil)
	fatalErr(err)

	candidate.Header.Add("X-SomeHeader", "some-value")
	candidate.Header.Add("X-SomeOtherHeader", "error and warning etc")

	client := newHTTPClient(false, 10, "")

	baseResp, err := client.Do(base)
	fatalErr(err)

	candidateResp, err := client.Do(candidate)
	fatalErr(err)

	differences, err := compare(baseResp, candidateResp)
	fatalErr(err)

	for _, d := range differences {
		fmt.Println(d.String())
	}

}

type diff struct {
	Kind         string
	Key          string
	BaseVal      string
	CandidateVal string
}

type diffs []diff

// b is base, c is candidate - abbreviated due to heavy use (:
func compare(b, c *http.Response) ([]diff, error) {

	out := make([]diff, 0)

	if b.Status != c.Status {
		out = append(out, diff{"status-line", "status", b.Status, c.Status})
	}

	if b.Proto != c.Proto {
		out = append(out, diff{"status-line", "proto", b.Proto, c.Proto})
	}

	// Different values, and in b but not c
	for header, vals := range b.Header {
		cVals := c.Header.Values(header)

		if len(cVals) == 0 {
			out = append(out, diff{"header-missing", header, strings.Join(vals, ", "), ""})
			continue
		}

		if slicesDiffer(vals, cVals) {
			out = append(out, diff{
				"header-values",
				header,
				strings.Join(vals, ", "),
				strings.Join(cVals, ", "),
			})
		}

	}

	// Header in c, but not b
	for header, vals := range c.Header {
		if b.Header.Get(header) == "" {
			out = append(out, diff{"header-added", header, "", strings.Join(vals, ", ")})
		}
	}

	// Body comparisons
	var err error
	var bBody []byte
	if b.Body != nil {
		bBody, err = ioutil.ReadAll(b.Body)
		if err != nil {
			return out, err
		}
	}

	var cBody []byte
	if c.Body != nil {
		cBody, err = ioutil.ReadAll(c.Body)
		if err != nil {
			return out, err
		}
	}

	// Length difference
	if len(bBody) != len(cBody) {
		out = append(out, diff{
			"body",
			"length",
			fmt.Sprintf("%d", len(bBody)),
			fmt.Sprintf("%d", len(cBody)),
		})
	}

	// Hash difference
	bHash := sha256.New()
	bHash.Write(bBody)
	bHashStr := fmt.Sprintf("%x", bHash.Sum(nil))

	cHash := sha256.New()
	cHash.Write(cBody)
	cHashStr := fmt.Sprintf("%x", cHash.Sum(nil))

	if bHashStr != cHashStr {
		out = append(out, diff{"body", "hash", bHashStr, cHashStr})
	}

	// MIME sniff
	bMIME := mimesniffer.Sniff(bBody)
	cMIME := mimesniffer.Sniff(cBody)
	if bMIME != cMIME {
		out = append(out, diff{"body", "mime", bMIME, cMIME})
	}

	// Content-type specific diffs

	// JSON type
	if bMIME == cMIME && bMIME == "application/json" {
		opts := jsondiff.DefaultJSONOptions()
		opts.SkipMatches = true
		jDiff, desc := jsondiff.Compare(bBody, cBody, &opts)
		switch jDiff {
		case jsondiff.NoMatch, jsondiff.SupersetMatch:
			out = append(out, diff{"body", "json", "", desc})
		case jsondiff.FirstArgIsInvalidJson:
			out = append(out, diff{"body", "json-fixed", string(bBody), string(cBody)})
		case jsondiff.SecondArgIsInvalidJson:
			out = append(out, diff{"body", "json-broken", string(bBody), string(cBody)})
		}
	}

	// Keywords
	keywords := []string{
		"error", "warn", "debug",
	}
	bBodyLower := strings.ToLower(string(bBody))
	cBodyLower := strings.ToLower(string(cBody))
	for _, k := range keywords {
		bCount := strings.Count(bBodyLower, k)
		cCount := strings.Count(cBodyLower, k)

		if bCount != cCount {
			out = append(out, diff{
				"keyword-count",
				k,
				fmt.Sprintf("%d", bCount),
				fmt.Sprintf("%d", cCount),
			})
		}

	}

	// HTML: tag count (for each tag type?)

	return out, nil
}

func slicesDiffer[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return true
	}

	for i, v := range a {
		if v != b[i] {
			return true
		}
	}

	return false
}

func inSlice[T comparable](s []T, v T) (bool, int) {
	for i, c := range s {
		if c == v {
			return true, i
		}
	}
	return false, -1
}

func (d diff) String() string {
	return fmt.Sprintf("%s[%s]: %s (base: %s)", d.Kind, d.Key, d.CandidateVal, d.BaseVal)
}

func fatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func newHTTPClient(keepAlives bool, timeout time.Duration, proxy string) *http.Client {

	tr := &http.Transport{
		MaxIdleConns:      30,
		IdleConnTimeout:   time.Second,
		DisableKeepAlives: !keepAlives,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   time.Second * timeout,
			KeepAlive: time.Second,
		}).DialContext,
	}

	if proxy != "" {
		if p, err := url.Parse(proxy); err == nil {
			tr.Proxy = http.ProxyURL(p)
		}
	}

	re := func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &http.Client{
		Transport:     tr,
		CheckRedirect: re,
		Timeout:       time.Second * 10,
	}

}
