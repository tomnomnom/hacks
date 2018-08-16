package main

import (
	"encoding/json"
	"fmt"
)

func fetchCertSpotter(domain string) ([]string, error) {
	out := make([]string, 0)

	raw, err := httpGet(
		fmt.Sprintf("https://certspotter.com/api/v0/certs?domain=%s", domain),
	)
	if err != nil {
		return out, err
	}

	wrapper := []struct {
		DNSNames []string `json:"dns_names"`
	}{}
	err = json.Unmarshal(raw, &wrapper)
	if err != nil {
		return out, err
	}

	for _, w := range wrapper {
		out = append(out, w.DNSNames...)
	}

	return out, nil
}
