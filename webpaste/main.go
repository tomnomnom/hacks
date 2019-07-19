package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type payload struct {
	Token string `json:"token"`
	Data  string `json:"data"`
}

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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if p.Token != token {
		http.Error(w, "lol no", http.StatusUnauthorized)
		return
	}

	fmt.Println(p.Data)
}

func main() {
	http.HandleFunc("/", payloadHandler)
	http.ListenAndServe(":8443", nil)
}
