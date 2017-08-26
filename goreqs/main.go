package main

import (
	"fmt"
	"log"
)

func main() {
	req := RawRequest{
		transport: "tls",
		host:      "httpbin.org",
		port:      "443",
		request:   "GET /anything HTTP/1.1\r\n" + "Host: httpbin.org\r\n" + "Connection: close\r\n",
	}

	resp, err := Do(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", resp)
}
