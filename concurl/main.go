package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func main() {
	var concurrency int
	flag.IntVar(&concurrency, "c", 20, "Concurrency level")

	var delay int
	flag.IntVar(&delay, "d", 5000, "Delay between requests to the same domain")

	var outputDir string
	flag.StringVar(&outputDir, "o", "out", "Output directory")

	flag.Parse()

	// channel to send URLs to workers
	jobs := make(chan string)

	rl := newRateLimiter(time.Duration(delay * 1000000))

	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)

		go func() {
			for u := range jobs {

				// get the domain for use in the path
				// and for rate limiting
				domain := "unknown"
				parsed, err := url.Parse(u)
				if err == nil {
					domain = parsed.Hostname()
				}

				// rate limit requests to the same domain
				rl.Block(domain)

				// we need the silent flag to get rid
				// of the progress output
				args := []string{"--silent", u}

				// pass all the arguments on to curl
				args = append(args, flag.Args()...)
				cmd := exec.Command("curl", args...)

				out, err := cmd.CombinedOutput()
				if err != nil {
					fmt.Printf("failed to get output: %s\n", err)
					continue
				}

				// use a hash of the URL and the arguments as the filename
				filename := fmt.Sprintf("%x", sha1.Sum([]byte(u+strings.Join(args, " "))))
				p := filepath.Join(outputDir, domain, filename)

				if _, err := os.Stat(path.Dir(p)); os.IsNotExist(err) {
					err = os.MkdirAll(path.Dir(p), 0755)
					if err != nil {
						fmt.Printf("failed to create output dir: %s\n", err)
						continue
					}
				}

				// include the command at the top of the output file
				buf := &bytes.Buffer{}
				buf.WriteString("cmd: curl ")
				buf.WriteString(strings.Join(args, " "))
				buf.WriteString("\n------\n\n")
				buf.Write(out)

				err = ioutil.WriteFile(p, buf.Bytes(), 0644)
				if err != nil {
					fmt.Printf("failed to save output: %s\n", err)
					continue
				}

				fmt.Printf("%s %s\n", p, u)
			}

			wg.Done()
		}()
	}

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		// send each line (a domain) on the jobs channel
		jobs <- sc.Text()
	}

	close(jobs)
	wg.Wait()
}
