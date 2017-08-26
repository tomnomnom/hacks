package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
)

func main() {
	req := RawRequest{
		transport: "tls",
		host:      "httpbin.org",
		port:      "443",
		request:   "GET /anything HTTP/1.1\r\n" + "Host: httpbin.org\r\n" + "Connection: close\r\n",
	}

	var conn io.ReadWriter
	var connerr error

	// This needs timeouts because it's fairly likely
	// that something will go wrong :)
	if req.IsTLS() {
		roots, err := x509.SystemCertPool()
		if err != nil {
			log.Fatal(err)
		}
		conn, connerr = tls.Dial(
			"tcp",
			fmt.Sprintf("%s:%s", req.Host(), req.Port()),
			&tls.Config{RootCAs: roots},
		)

	} else {
		conn, connerr = net.Dial(
			"tcp",
			fmt.Sprintf("%s:%s", req.Host(), req.Port()),
		)
	}

	if connerr != nil {
		log.Fatal(connerr)
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
