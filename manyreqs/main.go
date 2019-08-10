package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	var paramCount int
	flag.IntVar(&paramCount, "p", 40, "params per request")

	var concurrency int
	flag.IntVar(&concurrency, "c", 20, "concurrency level")

	flag.Parse()

	paramPatterns, err := getParamPatterns()
	if err != nil {
		log.Fatal(err)
	}

	headerPatterns, err := getHeaderPatterns()
	if err != nil {
		log.Fatal(err)
	}

	urls := make(chan string)

	// workers
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)

		c := getClient()
		go func() {
			defer wg.Done()

			for u := range urls {
				sendRequests(c, u, headerPatterns, paramPatterns, paramCount)
			}
		}()
	}

	sc := bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		urls <- sc.Text()
	}
	close(urls)

	wg.Wait()

}

func getClient() *http.Client {
	tr := &http.Transport{
		MaxIdleConns:    30,
		IdleConnTimeout: time.Second,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   time.Second * 10,
			KeepAlive: time.Second,
		}).DialContext,
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

func sendRequests(c *http.Client, u string, headerPatterns, paramPatterns map[string]string, paramCount int) {

	params, err := getParams(u, paramPatterns)
	headers, err := getHeaders(u, headerPatterns)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	count := 0
	chunks := make([]string, 0)
	var buf strings.Builder
	for k, v := range params {
		buf.WriteString(k)
		buf.WriteRune('=')
		buf.WriteString(v)
		buf.WriteRune('&')

		count++
		if count > paramCount {
			count = 0
			chunks = append(chunks, buf.String())
			buf.Reset()
		}
	}

	for _, chunk := range chunks {

		req, err := http.NewRequest("GET", u+"?"+chunk, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		for h, v := range headers {
			req.Header.Set(h, v)
		}

		resp, err := c.Do(req)
		if resp != nil {
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "error requesting %s: %s\n", u, err)
			return
		}
	}

}

func getHeaders(raw string, patterns map[string]string) (map[string]string, error) {
	out := make(map[string]string)

	u, err := url.Parse(raw)
	if err != nil {
		return out, err
	}

	for h, v := range patterns {
		out[h] = fmt.Sprintf(v, u.Hostname())
	}

	return out, nil
}

func getHeaderPatterns() (map[string]string, error) {
	out := make(map[string]string)

	f, err := os.Open("headers")
	if err != nil {
		return out, err
	}

	sc := bufio.NewScanner(f)

	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		out[parts[0]] = strings.TrimSpace(parts[1])
	}

	return out, nil
}

func getParamPatterns() (map[string]string, error) {
	out := make(map[string]string)

	f, err := os.Open("params")
	if err != nil {
		return out, err
	}

	sc := bufio.NewScanner(f)

	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		out[parts[0]] = strings.TrimSpace(parts[1])
	}

	return out, nil
}

func getParams(raw string, patterns map[string]string) (map[string]string, error) {
	out := make(map[string]string)

	u, err := url.Parse(raw)
	if err != nil {
		return out, err
	}

	for h, v := range patterns {
		out[h] = fmt.Sprintf(v, u.Hostname())
	}

	return out, nil
}
