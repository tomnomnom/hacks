package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const searchURL = "https://archive.org/advancedsearch.php?q=collection%3A%28UrlteamWebCrawls%29&fl%5B%5D=identifier&sort%5B%5D=addeddate+desc&sort%5B%5D=&sort%5B%5D=&rows=5&page=1&output=json"
const metaURL = "http://archive.org/metadata/%s"

type file struct {
	Name   string `json:"name"`
	Format string `json:"format"`
}

func main() {
	res, err := http.Get(searchURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()
	dec := json.NewDecoder(res.Body)

	wrapper := &struct {
		Response struct {
			Docs []struct {
				Identifier string `json:"identifier"`
			} `json:"docs"`
		} `json:"response"`
	}{}

	err = dec.Decode(wrapper)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, d := range wrapper.Response.Docs {
		files, err := getFiles(d.Identifier)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, f := range files {
			if f.Format != "ZIP" {
				continue
			}
			fmt.Printf("https://archive.org/download/%s/%s\n", d.Identifier, f.Name)
		}
	}
}

func getFiles(ident string) ([]file, error) {

	res, err := http.Get(fmt.Sprintf(metaURL, ident))
	if err != nil {
		return []file{}, err
	}
	defer res.Body.Close()

	wrapper := &struct {
		Files []file `json:"files"`
	}{}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(wrapper)

	return wrapper.Files, err

}
