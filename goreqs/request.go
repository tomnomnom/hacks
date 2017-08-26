package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"
)

type Request interface {
	IsTLS() bool
	Host() string
	Port() string
	String() string
}

type RawRequest struct {
	transport string
	host      string
	port      string
	request   string
}

func (r RawRequest) IsTLS() bool {
	return r.transport == "tls"
}

func (r RawRequest) Host() string {
	return r.host
}

func (r RawRequest) Port() string {
	return r.port
}

func (r RawRequest) String() string {
	return r.request
}

func Do(req Request) (*Response, error) {
	var conn io.ReadWriter
	var connerr error

	// This needs timeouts because it's fairly likely
	// that something will go wrong :)
	if req.IsTLS() {
		roots, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
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
		return nil, connerr
	}

	fmt.Fprintf(conn, req.String())
	fmt.Fprintf(conn, "\r\n")

	return NewResponse(conn)

}
