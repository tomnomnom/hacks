package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

func main() {

	var depth int
	flag.IntVar(&depth, "depth", 4, "max recursion depth")
	flag.IntVar(&depth, "d", 4, "max recursion depth")

	flag.Parse()

	suffix := flag.Arg(0)
	wordListFile := flag.Arg(1)

	if suffix == "" {
		fmt.Fprintln(os.Stderr, "usage: ettu [--depth=<int>] <domain> [<wordfile>|-]")
		return
	}

	var f io.Reader
	var err error

	// default to stdin for the wordlist
	f = os.Stdin

	if wordListFile != "" && wordListFile != "-" {
		f, err = os.Open(wordListFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open word list: %s\n", err)
			return
		}
	}

	sc := bufio.NewScanner(f)

	words := make([]string, 0)
	for sc.Scan() {
		words = append(words, sc.Text())
	}

	out := make(chan string, 1000)

	go func() {
		for o := range out {
			fmt.Println(o)
		}
	}()

	brute(suffix, words, out, 1, depth)
}

func brute(suffix string, words []string, out chan string, depth, maxDepth int) {
	if depth > maxDepth {
		return
	}

	var wg sync.WaitGroup

	for _, w := range words {
		candidate := fmt.Sprintf("%s.%s", w, suffix)

		_, err := net.LookupHost(candidate)

		if err != nil {
			// need a DNSError so we can check some of its fields
			nerr, ok := err.(*net.DNSError)
			if !ok {
				_ = nerr
				continue
			}

			// why you makin' me do this, Go? :(
			if nerr.IsTimeout || nerr.Err == "no such host" {
				continue
			}
		}

		// recurse for any candidate that either resolves or
		// isn't specifically a 'no such host' error as it may
		// have other subdomains which resolve
		wg.Add(1)
		go func() {
			brute(candidate, words, out, depth+1, maxDepth)
			wg.Done()
		}()

		if err != nil {
			continue
		}

		// output any name that didn't have an error (i.e. it resolves)
		out <- candidate
	}

	wg.Wait()
}
