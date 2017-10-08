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
		domainsFile = "domains"
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
		name      string
		wildcards []string
		subs      chan string
	}

	// make a slice of all the domains and spin up a worker for each one
	domains := make([]domain, 0, 0)
	wg := sync.WaitGroup{}

	for ds.Scan() {
		name := ds.Text()

		wildcards, _ := net.LookupHost(fmt.Sprintf("%s.%s", wildcardCheck, name))

		subs := make(chan string)
		d := domain{name, wildcards, subs}
		domains = append(domains, d)

		wg.Add(1)
		go func() {
			for sub := range d.subs {
				candidate := fmt.Sprintf("%s.%s", sub, d.name)
				if subdomainExists(candidate, d.wildcards) {
					fmt.Println(candidate)
				}
			}
			wg.Done()
		}()
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

func subdomainExists(subdomain string, wildcards []string) bool {
	addrs, err := net.LookupHost(subdomain)
	if err != nil {
		return false
	}

	for _, addr := range addrs {
		for _, wildcard := range wildcards {
			if addr == wildcard {
				return false
			}
		}
	}

	return true
}
