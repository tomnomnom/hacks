package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

func main() {
	sc := bufio.NewScanner(os.Stdin)

	var wg sync.WaitGroup

	for sc.Scan() {
		domain := sc.Text()

		wg.Add(1)
		go func() {

			if _, err := net.LookupHost(domain); err == nil {
				fmt.Println(domain)
			}

			wg.Done()
		}()
	}
	wg.Wait()
}
