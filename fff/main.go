package main

import (
	"bufio"
	"crypto/sha1"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

func main() {

	var keepAlives bool
	flag.BoolVar(&keepAlives, "keep-alive", false, "use HTTP keep-alives")

	var delayMs int
	flag.IntVar(&delayMs, "delay", 100, "delay between issuing requests (ms)")

	flag.Parse()

	delay := time.Duration(delayMs * 1000000)
	client := newClient(keepAlives)
	prefix := "out"

	var wg sync.WaitGroup

	sc := bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		rawURL := sc.Text()
		wg.Add(1)

		time.Sleep(delay)

		go func() {
			defer wg.Done()

			req, err := http.NewRequest("GET", rawURL, nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to create request: %s\n", err)
				return
			}

			resp, err := client.Do(req)
			if err != nil {
				fmt.Fprintf(os.Stderr, "request failed: %s\n", err)
				return
			}
			defer resp.Body.Close()

			hash := sha1.Sum([]byte(rawURL))
			p := path.Join(prefix, req.URL.Hostname(), fmt.Sprintf("%x.body", hash))
			err = os.MkdirAll(path.Dir(p), 0750)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to create dir: %s\n", err)
				return
			}

			f, err := os.Create(p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to create file: %s\n", err)
				return
			}

			// TODO: write response headers to a second file
			_, err = io.Copy(f, resp.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to write file contents: %s\n", err)
				return
			}

			fmt.Printf("%s: %s\n", p, rawURL)

		}()
	}

	wg.Wait()

}

func newClient(keepAlives bool) *http.Client {

	tr := &http.Transport{
		MaxIdleConns:      30,
		IdleConnTimeout:   time.Second,
		DisableKeepAlives: !keepAlives,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
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
