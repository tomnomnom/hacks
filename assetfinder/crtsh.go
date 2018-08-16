package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func fetchCrtSh(domain string) ([]string, error) {
	resp, err := http.Get(
		fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain),
	)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	output := make([]string, 0)

	dec := json.NewDecoder(resp.Body)

	for {
		wrapper := struct {
			Name string `json:"name_value"`
		}{}

		err := dec.Decode(&wrapper)
		if err != nil {
			break
		}

		output = append(output, wrapper.Name)
	}
	return output, nil
}
