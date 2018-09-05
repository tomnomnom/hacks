package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"
)

const searchURL = "https://archive.org/advancedsearch.php"
const metaURL = "http://archive.org/metadata/%s"

type file struct {
	Name   string `json:"name"`
	Format string `json:"format"`
}

func main() {
	flag.Parse()

	sinceStr := flag.Arg(0)
	if sinceStr == "" {
		fmt.Println("usage: urlteamdl <sinceISODate>")
		return
	}

	sinceTime, err := time.Parse("2006-01-02", sinceStr)
	if err != nil {
		fmt.Println("invalid date! try e.g. 2017-10-26")
		return
	}

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	since := sinceTime.Format("2006-01-02")
	today := time.Now().Format("2006-01-02")

	q := req.URL.Query()
	q.Add("q", fmt.Sprintf("collection:(UrlteamWebCrawls) AND addeddate:[%s TO %s]", since, today))
	q.Add("fl[]", "identifier")
	q.Add("sort[]", "addeddate desc")
	q.Add("rows", "500")
	q.Add("output", "json")
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
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
		files, err := getDownloadURLs(d.Identifier)
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

func getDownloadURLs(ident string) ([]file, error) {

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
