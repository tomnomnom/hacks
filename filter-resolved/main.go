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

	concurrency := 20
	flag.IntVar(&concurrency, "c", 20, "Set the concurrency level")

	flag.Parse()

	jobs := make(chan string)

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)

		go func() {
			for domain := range jobs {
				_, err := net.ResolveIPAddr("ip4", domain)
				if err != nil {
					continue
				}
				fmt.Println(domain)
			}
			wg.Done()
		}()
	}

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		jobs <- sc.Text()
	}
	close(jobs)

	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %s\n", err)
	}
	wg.Wait()

}
