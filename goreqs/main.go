package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
)

type Response struct {
	rawStatus string
	headers   []string
	body      []byte
}

func (r *Response) AddHeader(header string) {
	r.headers = append(r.headers, header)
}

func (r Response) Header(search string) string {
	search = strings.ToLower(search)

	for _, header := range r.headers {

		p := strings.SplitN(header, ":", 2)
		if len(p) != 2 {
			continue
		}

		if strings.ToLower(p[0]) == search {
			return strings.TrimSpace(p[1])
		}
	}
	return ""
}

func main() {
	req := RawRequest{
		transport: "tcp",
		host:      "httpbin.org",
		port:      "80",
		request:   "GET /anything HTTP/1.1\r\n" + "Host: httpbin.org\r\n" + "Connection: close\r\n",
	}

	conn, err := net.Dial(
		req.Transport(),
		fmt.Sprintf("%s:%s", req.Host(), req.Port()),
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(conn, req.String())
	fmt.Fprintf(conn, "\r\n")

	resp := &Response{}

	r := bufio.NewReader(conn)
	s, err := r.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	resp.rawStatus = strings.TrimSpace(s)

	for {
		line, err := r.ReadString('\n')
		line = strings.TrimSpace(line)

		if err != nil || line == "" {
			break
		}

		resp.AddHeader(line)
	}

	if cl := resp.Header("Content-Length"); cl != "" {
		length, err := strconv.Atoi(cl)

		if err != nil {
			log.Fatal(err)
		}

		if length > 0 {
			b := make([]byte, length)
			_, err = io.ReadAtLeast(r, b, length)
			if err != nil {
				log.Fatal(err)
			}
			resp.body = b
		}

	} else {
		b, err := ioutil.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}

		resp.body = b
	}

	fmt.Printf("%#v\n", resp)

}
