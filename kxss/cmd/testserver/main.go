package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		qs := r.URL.Query()
		log.Printf("req: %#v", qs)
		w.Header().Add("Content-Type", "text/html")
		fmt.Fprintf(w, "Hello, %s!\n", qs.Get("name"))
		fmt.Fprintf(w, "I hear you're %s years old!\n", qs.Get("age"))
		fmt.Fprint(w, "My name is Buck and I'm here to greet you.\n")
	})

	http.ListenAndServe("127.0.0.1:5566", nil)
}
