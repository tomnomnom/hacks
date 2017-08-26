package main

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
