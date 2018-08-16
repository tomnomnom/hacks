package main

import (
	"encoding/json"
	"fmt"
)

func fetchThreatCrowd(domain string) ([]string, error) {
	out := make([]string, 0)

	raw, err := httpGet(
		fmt.Sprintf("https://www.threatcrowd.org/searchApi/v2/domain/report/?domain=%s", domain),
	)
	if err != nil {
		return out, err
	}

	wrapper := struct {
		Subdomains []string `json:"subdomains"`
	}{}
	err = json.Unmarshal(raw, &wrapper)
	if err != nil {
		return out, err
	}

	out = append(out, wrapper.Subdomains...)

	return out, nil
}
