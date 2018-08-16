package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func fetchVirusTotal(domain string) ([]string, error) {

	apiKey := os.Getenv("VT_API_KEY")
	if apiKey == "" {
		// swallow not having an API key, just
		// don't fetch
		return []string{}, nil
	}

	resp, err := http.Get(fmt.Sprintf(
		"https://www.virustotal.com/vtapi/v2/domain/report?domain=%s&apikey=%s",
		domain, apiKey,
	))
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	wrapper := struct {
		Subdomains []string `json:"subdomains"`
	}{}

	err = dec.Decode(&wrapper)
	return wrapper.Subdomains, err
}
