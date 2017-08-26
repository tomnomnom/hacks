package main

import (
	"bufio"
	"io"
	"io/ioutil"
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

func NewResponse(conn io.Reader) (*Response, error) {

	r := bufio.NewReader(conn)
	resp := &Response{}

	s, err := r.ReadString('\n')
	if err != nil {
		return nil, err
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
			return nil, err
		}

		if length > 0 {
			b := make([]byte, length)
			_, err = io.ReadAtLeast(r, b, length)
			if err != nil {
				return nil, err
			}
			resp.body = b
		}

	} else {
		b, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}

		resp.body = b
	}

	return resp, nil
}
