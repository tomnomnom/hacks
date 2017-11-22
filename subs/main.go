package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

// lol
const wildcardCheck = "lkj23lk52lkn23kjh23mnbzzxckjhasdqwe"

func main() {
	flag.Parse()

	domainsFile := flag.Arg(0)
	subdomainsFile := flag.Arg(1)

	if domainsFile == "" {
		domainsFile = "apexes"
	}

	if subdomainsFile == "" {
		subdomainsFile = "subdomains"
	}

	d, err := os.Open(domainsFile)
	if err != nil {
		log.Fatal(err)
	}

	s, err := os.Open(subdomainsFile)
	if err != nil {
		log.Fatal(err)
	}

	ds := bufio.NewScanner(d)
	ss := bufio.NewScanner(s)

	type domain struct {
		name string
		subs chan string
	}

	// make a slice of all the domains and spin up a worker for each one
	domains := make([]domain, 0, 0)
	wg := sync.WaitGroup{}

	for ds.Scan() {
		name := ds.Text()

		_, err := net.LookupHost(fmt.Sprintf("%s.%s", wildcardCheck, name))
		if err == nil {
			// There's a wildcard, don't bother
			continue
		}

		subs := make(chan string, 128)
		d := domain{name, subs}
		domains = append(domains, d)

		worker := func() {
			for sub := range d.subs {
				candidate := fmt.Sprintf("%s.%s", sub, d.name)
				if subdomainExists(candidate) {
					fmt.Println(candidate)
				}
			}
			wg.Done()
		}
		wg.Add(2)
		go worker()
		go worker()
	}

	// check for errors reading the list of domains
	if err := ds.Err(); err != nil {
		log.Println(err)
	}

	// look up each sub for every domain
	for ss.Scan() {
		sub := ss.Text()

		for _, d := range domains {
			d.subs <- sub
		}
	}

	if err := ss.Err(); err != nil {
		log.Fatal(err)
	}

	// we're done, close all the channels
	for _, d := range domains {
		close(d.subs)
	}

	wg.Wait()
}

func subdomainExists(subdomain string) bool {
	_, err := net.LookupHost(subdomain)
	if err != nil {
		return false
	}

	return true
}
