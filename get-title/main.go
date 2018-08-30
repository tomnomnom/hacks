package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"

	"golang.org/x/net/html"
)

func main() {

	var concurrency = 20
	flag.IntVar(&concurrency, "c", 20, "Concurrency")
	flag.Parse()

	jobs := make(chan string)

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			for j := range jobs {
				resp, err := http.Get(j)
				if err != nil {
					continue
				}

				z := html.NewTokenizer(resp.Body)

				for {
					tt := z.Next()
					if tt == html.ErrorToken {
						break
					}

					t := z.Token()

					if t.Type == html.StartTagToken && t.Data == "title" {
						if z.Next() == html.TextToken {
							fmt.Printf("%s (%s)\n", z.Token().Data, j)
						}
					}

				}

				resp.Body.Close()

			}
			wg.Done()
		}()
	}

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		jobs <- sc.Text()
	}
	close(jobs)

	wg.Wait()

}
