package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

var (
	client    *http.Client
	transport *http.Transport
	wg        sync.WaitGroup

	concurrency = 20
)

func main() {
	flag.Parse()

	var input io.Reader
	input = os.Stdin

	if flag.NArg() > 0 {
		file, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Printf("failed to open file: %s\n", err)
			os.Exit(1)
		}
		input = file
	}

	sc := bufio.NewScanner(input)

	client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        concurrency,
			MaxIdleConnsPerHost: concurrency,
			MaxConnsPerHost:     concurrency,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 5 * time.Second,
	}

	semaphore := make(chan bool, concurrency)

	for sc.Scan() {
		raw := sc.Text()
		wg.Add(1)
		semaphore <- true

		go func(raw string) {
			defer wg.Done()
			u, err := url.ParseRequestURI(raw)
			if err != nil {
				return
			}
			resp, err := fetchURL(u)
			if err != nil {
				return
			}
			if resp.StatusCode == http.StatusOK {
				fmt.Println(u)
			}
			<-semaphore
		}(raw)
	}

	wg.Wait()

	if sc.Err() != nil {
		fmt.Printf("error: %s\n", sc.Err())
	}
}

func fetchURL(u *url.URL) (*http.Response, error) {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "burl/0.1")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	io.Copy(ioutil.Discard, resp.Body)

	return resp, err
}
