package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type payload struct {
	Token string `json:"token"`
	File  string `json:"file"`
	Data  string `json:"data"`
}

func payloadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/plain")

	if r.Method != "POST" {
		http.Error(w, "no.", http.StatusMethodNotAllowed)
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

	// base name only plz
	p.File = filepath.Base(p.File)

	f, err := os.OpenFile(p.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = f.WriteString(p.Data + "\n")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("wrote to %s\n", p.File)

	fmt.Sprint(w, "did it\n")
}

func main() {
	http.HandleFunc("/", payloadHandler)

	log.Println("Go!")
	http.ListenAndServe(":8443", nil)
}
