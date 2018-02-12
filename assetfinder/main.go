package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	flag.Parse()

	domain := flag.Arg(0)
	if domain == "" {
		fmt.Println("no domain specified")
		return
	}

	// TODO:
	//   https://crt.sh/?q=%25.%s&output=json
	//   crt.sh's return format isn't actually valid JSON;
	//   it's a series of JSON objects with no delimiter
	//   between them. I mean seriously, why would you
	//   do that?!
	sources := []fetchFn{
		fetchCertSpotter,
		fetchHackerTarget,
		fetchThreatCrowd,
	}

	// track what we've already printed to avoid duplicates
	printed := make(map[string]bool)

	for _, source := range sources {
		names, err := source(domain)
		if err != nil {
			log.Fatal(err)
		}

		for _, n := range names {
			n = cleanDomain(n)
			if _, ok := printed[n]; ok {
				continue
			}
			fmt.Println(n)
			printed[n] = true
		}
	}
}

type fetchFn func(string) ([]string, error)

func fetchThreatCrowd(domain string) ([]string, error) {
	out := make([]string, 0)

	raw, err := httpGet(
		fmt.Sprintf("https://www.threatcrowd.org/searchApi/v2/domain/report/?domain=%s", domain),
	)
	if err != nil {
		return out, err
	}

	wrapper := struct {
		Subdomains []string `json:"subdomains"`
	}{}
	err = json.Unmarshal(raw, &wrapper)
	if err != nil {
		return out, err
	}

	out = append(out, wrapper.Subdomains...)

	return out, nil
}

func fetchHackerTarget(domain string) ([]string, error) {
	out := make([]string, 0)

	raw, err := httpGet(
		fmt.Sprintf("https://api.hackertarget.com/hostsearch/?q=%s", domain),
	)
	if err != nil {
		return out, err
	}

	sc := bufio.NewScanner(bytes.NewReader(raw))
	for sc.Scan() {
		parts := strings.SplitN(sc.Text(), ",", 2)
		if len(parts) != 2 {
			continue
		}

		out = append(out, parts[0])
	}

	return out, sc.Err()
}

func fetchCertSpotter(domain string) ([]string, error) {
	out := make([]string, 0)

	raw, err := httpGet(
		fmt.Sprintf("https://certspotter.com/api/v0/certs?domain=%s", domain),
	)
	if err != nil {
		return out, err
	}

	wrapper := []struct {
		DNSNames []string `json:"dns_names"`
	}{}
	err = json.Unmarshal(raw, &wrapper)
	if err != nil {
		return out, err
	}

	for _, w := range wrapper {
		out = append(out, w.DNSNames...)
	}

	return out, nil
}

func httpGet(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}

	raw, err := ioutil.ReadAll(res.Body)

	res.Body.Close()
	if err != nil {
		return []byte{}, err
	}

	return raw, nil
}

func cleanDomain(d string) string {
	// no idea what this is, but we can't clean it ¯\_(ツ)_/¯
	if len(d) < 2 {
		return d
	}

	if d[0] == '*' || d[0] == '%' {
		d = d[1:]
	}

	if d[0] == '.' {
		d = d[1:]
	}

	return d

}
