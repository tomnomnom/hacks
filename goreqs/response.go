package main

import "strings"

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
