package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
)

type payload struct {
	Token string   `json:"token"`
	Lines []string `json:"lines"`
}

// lol globals
var bus chan []string

func payloadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/plain")

	if r.Method != "POST" {
		fmt.Fprint(w, "not gonna happen sorry")
		return
	}

	token := os.Getenv("WEBPASTE_TOKEN")

	d := json.NewDecoder(r.Body)
	p := &payload{}
	err := d.Decode(p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode JSON: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if p.Token != token {
		fmt.Fprintf(os.Stderr, "got invalid token")
		http.Error(w, "lol no", http.StatusUnauthorized)
		return
	}

	bus <- p.Lines
}

func main() {
	var unique bool
	flag.BoolVar(&unique, "u", false, "only print unique lines")

	var port string
	flag.StringVar(&port, "p", "8080", "port to listen on")

	var address string
	flag.StringVar(&address, "a", "0.0.0.0", "address to listen on")

	flag.Parse()

	if os.Getenv("WEBPASTE_TOKEN") == "" {
		fmt.Fprintln(os.Stderr, "WEBPASTE_TOKEN is not set")
		return
	}

	bus = make(chan []string)

	go func() {
		seen := make(map[string]bool)

		for ss := range bus {
			for _, s := range ss {
				if unique && seen[s] {
					continue
				}
				seen[s] = true
				fmt.Println(s)
			}
		}
	}()

	http.HandleFunc("/", payloadHandler)
	http.ListenAndServe(address+":"+port, nil)
}
