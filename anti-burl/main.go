package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var deadHosts = make(map[string]struct{})
var deadMutex = &sync.Mutex{}
var aliveHosts = make(map[string]struct{})
var aliveMutex = &sync.Mutex{}
var statusCodes = flag.String("s", "200", "Status Codes to accept separated by a comma")
var inputFile = flag.String("w", "", "File containing URLS")
var codes []string

func main() {
	flag.Parse()
	if strings.Compare(*statusCodes, "200") != 0 {
		codes = strings.Split(*statusCodes, ",")
	} else {
		codes = append(codes, "200")
	}
	var input io.Reader
	input = os.Stdin

	if strings.Compare(*inputFile, "") != 0 {
		file, err := os.Open(*inputFile)
		if err != nil {
			fmt.Printf("failed to open file: %s\n", err)
			os.Exit(1)
		}
		input = file
	}

	sc := bufio.NewScanner(input)

	urls := make(chan string, 128)
	concurrency := 12
	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			for raw := range urls {

				u, err := url.ParseRequestURI(raw)
				if err != nil {
					continue
				}

				if !resolves(u) {
					continue
				}

				resp, err := fetchURL(u)
				if err != nil {
					continue
				}

				if contains(codes, resp.StatusCode) {
					fmt.Printf("%s\n", u)
				}
			}
			wg.Done()
		}()
	}

	for sc.Scan() {
		urls <- sc.Text()
	}
	close(urls)

	if sc.Err() != nil {
		fmt.Printf("error: %s\n", sc.Err())
	}

	wg.Wait()
}
func contains(status []string, code int) bool {
	for _, stat := range status {
		stat, _ := strconv.Atoi(stat)
		if stat == code {
			return true
		}
	}
	return false
}
func resolves(u *url.URL) bool {
	aliveMutex.Lock()
	if _, ok := aliveHosts[u.Hostname()]; ok {
		return true
	}
	aliveMutex.Unlock()

	deadMutex.Lock()
	if _, ok := deadHosts[u.Hostname()]; ok {
		return false
	}
	deadMutex.Unlock()

	addrs, _ := net.LookupHost(u.Hostname())
	if len(addrs) == 0 {
		deadMutex.Lock()
		deadHosts[u.Hostname()] = struct{}{}
		deadMutex.Unlock()
	} else {
		aliveMutex.Lock()
		aliveHosts[u.Hostname()] = struct{}{}
		aliveMutex.Unlock()
	}
	return len(addrs) != 0
}

func fetchURL(u *url.URL) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{
		Transport: tr,
		Timeout:   5 * time.Second,
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Set("User-Agent", "burl/0.1")

	resp, err := client.Do(req)
	if resp != nil {
		resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	return resp, err
}
