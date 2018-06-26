package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
)

func main() {

	servers := []string{
		//"209.244.0.3",
		//"209.244.0.4",
		//"64.6.64.6",
		//"64.6.65.6",
		"8.8.8.8",
		"8.8.4.4",
		"9.9.9.9",
		//	"149.112.112.112",
		//	"84.200.69.80",
		//	"84.200.70.40",
		//	"8.26.56.26",
		//	"8.20.247.20",
		//	"208.67.222.222",
		//	"208.67.220.220",
		//	"199.85.126.10",
		//	"199.85.127.10",
		//	"81.218.119.11",
		//	"209.88.198.133",
		//	"195.46.39.39",
		//	"195.46.39.40",
		//	"69.195.152.204",
		//	"23.94.60.240",
		//	"208.76.50.50",
		//	"208.76.51.51",
		//	"216.146.35.35",
		//	"216.146.36.36",
		//	"37.235.1.174",
		//	"37.235.1.177",
		//	"198.101.242.72",
		//	"23.253.163.53",
		//	"77.88.8.8",
		//	"77.88.8.1",
		//	"91.239.100.100",
		//	"89.233.43.71",
		//	"74.82.42.42",
		//	"109.69.8.51",
		//	"156.154.70.1",
		//	"156.154.71.1",
		"1.1.1.1",
		"1.0.0.1",
		//	"45.77.165.194",
	}

	rand.Seed(time.Now().Unix())

	type job struct{ domain, server string }
	jobs := make(chan job)

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)

		go func() {
			for j := range jobs {

				cname, err := getCNAME(j.domain, j.server)
				if err != nil {
					//fmt.Println(err)
					continue
				}

				if !resolves(cname) {
					fmt.Printf("%s does not resolve (pointed at by %s)\n", cname, j.domain)
				}
			}
			wg.Done()
		}()
	}

	sc := bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		target := strings.ToLower(strings.TrimSpace(sc.Text()))
		if target == "" {
			continue
		}
		server := servers[rand.Intn(len(servers))]

		jobs <- job{target, server}
	}
	close(jobs)

	wg.Wait()

}

func resolves(domain string) bool {
	_, err := net.LookupHost(domain)
	return err == nil
}

func getCNAME(domain, server string) (string, error) {
	c := dns.Client{}

	m := dns.Msg{}
	if domain[len(domain)-1:] != "." {
		domain += "."
	}
	m.SetQuestion(domain, dns.TypeCNAME)
	m.RecursionDesired = true

	r, _, err := c.Exchange(&m, server+":53")
	if err != nil {
		return "", err
	}

	if len(r.Answer) == 0 {
		return "", fmt.Errorf("no answers for %s", domain)
	}

	for _, ans := range r.Answer {
		if r, ok := ans.(*dns.CNAME); ok {
			return r.Target, nil
		}
	}
	return "", fmt.Errorf("no cname for %s", domain)

}
