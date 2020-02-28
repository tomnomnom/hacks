package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type appendCheck struct {
	url   string
	param string
}

var transport = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: time.Second,
		DualStack: true,
	}).DialContext,
}

var httpClient = &http.Client{
	Transport: transport,
}

func main() {

	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	sc := bufio.NewScanner(os.Stdin)

	initialChecks := make(chan string, 40)
	appendChecks := make(chan appendCheck, 40)
	charChecks := make(chan appendCheck, 40)

	// initial check worker pool
	var wgInitial sync.WaitGroup
	for i := 0; i < 40; i++ {
		wgInitial.Add(1)
		go func() {
			for inputURL := range initialChecks {

				reflected, err := checkReflected(inputURL)
				if err != nil {
					//fmt.Fprintf(os.Stderr, "error from checkReflected: %s\n", err)
					continue
				}

				if len(reflected) == 0 {
					// TODO: wrap in verbose mode
					//fmt.Printf("no params were reflected in %s\n", inputURL)
					continue
				}

				for _, param := range reflected {
					appendChecks <- appendCheck{inputURL, param}
				}

			}
			wgInitial.Done()
		}()
	}

	// append check worker pool
	var wgAppend sync.WaitGroup
	for i := 0; i < 40; i++ {
		wgAppend.Add(1)

		go func() {
			for c := range appendChecks {
				wasReflected, err := checkAppend(c.url, c.param, "iy3j4h234hjb23234")
				if err != nil {
					fmt.Fprintf(os.Stderr, "error from checkAppend for url %s with param %s: %s", c.url, c.param, err)
					continue
				}

				if wasReflected {
					charChecks <- appendCheck{c.url, c.param}
				}
			}
			wgAppend.Done()
		}()

	}

	// char check worker pool
	var wgChar sync.WaitGroup
	for i := 0; i < 40; i++ {
		wgChar.Add(1)

		go func() {
			for c := range charChecks {
				for _, char := range []string{"\"", "'", "<", ">"} {
					wasReflected, err := checkAppend(c.url, c.param, "aprefix"+char+"asuffix")
					if err != nil {
						fmt.Fprintf(os.Stderr, "error from checkAppend for url %s with param %s with %s: %s", c.url, c.param, char, err)
						continue
					}

					if wasReflected {
						fmt.Printf("got reflection of appended param %s with %s in value on %s\n", c.param, char, c.url)
					}
				}
			}
			wgChar.Done()
		}()

	}

	for sc.Scan() {
		initialChecks <- sc.Text()
	}

	// this is silly. Need to refactor into something a bit less repetitive
	close(initialChecks)
	wgInitial.Wait()
	close(appendChecks)
	wgAppend.Wait()
	close(charChecks)
	wgChar.Wait()

}

func checkReflected(targetURL string) ([]string, error) {

	out := make([]string, 0)

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return out, err
	}

	// temporary. Needs to be an option
	req.Header.Add("User-Agent", "User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.100 Safari/537.36")

	resp, err := httpClient.Do(req)
	if err != nil {
		return out, err
	}
	if resp.Body == nil {
		return out, err
	}
	defer resp.Body.Close()

	// always read the full body so we can re-use the tcp connection
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return out, err
	}

	// nope (:
	if strings.HasPrefix(resp.Status, "3") {
		return out, nil
	}

	// also nope
	ct := resp.Header.Get("Content-Type")
	if ct != "" && !strings.Contains(ct, "html") {
		return out, nil
	}

	body := string(b)

	u, err := url.Parse(targetURL)
	if err != nil {
		return out, err
	}

	for key, vv := range u.Query() {
		for _, v := range vv {
			if !strings.Contains(body, v) {
				continue
			}

			out = append(out, key)
		}
	}

	return out, nil
}

func checkAppend(targetURL, param, suffix string) (bool, error) {
	u, err := url.Parse(targetURL)
	if err != nil {
		return false, err
	}

	qs := u.Query()
	val := qs.Get(param)
	//if val == "" {
	//return false, nil
	//return false, fmt.Errorf("can't append to non-existant param %s", param)
	//}

	qs.Set(param, val+suffix)
	u.RawQuery = qs.Encode()

	reflected, err := checkReflected(u.String())
	if err != nil {
		return false, err
	}

	for _, r := range reflected {
		if r == param {
			return true, nil
		}
	}

	return false, nil
}
