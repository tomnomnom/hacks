package main

type Request interface {
	Transport() string
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

func (r RawRequest) Transport() string {
	return r.transport
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
