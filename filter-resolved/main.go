package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
)

func main() {

	// configure the concurrency flag
	concurrency := 20
	flag.IntVar(&concurrency, "c", 20, "Set the concurrency level")

	// parse the flags
	flag.Parse()

	// jobs is a channel of strings. We'll send domains on the
	// channel so that a bunch of workers can receive them and
	// try to resolve them
	jobs := make(chan string)

	// A WaitGroup is useful if you have lots of goroutines
	// and you want to know when they're all done.
	var wg sync.WaitGroup

	// spin up a whole bunch of workers
	for i := 0; i < concurrency; i++ {
		// tell the waitgroup about the new worker
		wg.Add(1)

		// launch a goroutine that takes domains off the
		// jobs channel, tries to resolve them and outputs
		// them only if there was no error
		go func() {
			for domain := range jobs {
				_, err := net.ResolveIPAddr("ip4", domain)
				if err != nil {
					continue
				}
				fmt.Println(domain)
			}

			// when the jobs channel is closed the loop
			// above will stop; then we need to tell the
			// waitgroup that the worker is done
			wg.Done()
		}()
	}

	// open stdin as a scanner. That makes it super easy
	// to deal with line-delimited input
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		// send each line (a domain) on the jobs channel
		jobs <- sc.Text()
	}

	// as soon as we're done sending all the jobs we can
	// close the jobs channel. If we don't the workers
	// will never stop.
	close(jobs)

	// check there were no errors reading stdin (unlikely)
	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %s\n", err)
	}

	// wait for the workers to finish doing their thing
	wg.Wait()

}
