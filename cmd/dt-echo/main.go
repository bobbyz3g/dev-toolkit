// dt-echo is a http echo server, it will
// echo back the request body and headers
package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
)

func main() {
	listen := flag.String("listen", "0.0.0.0:8080", "listen address")
	flag.Parse()
	http.HandleFunc("/", echo)

	log.Println("Server started on ", *listen)
	err := http.ListenAndServe(*listen, nil)
	if err != nil {
		log.Fatal(err)
	}
}

type EchoBody struct {
	RemoteAddr string      `json:"remote_addr"`
	Method     string      `json:"method"`
	URL        string      `json:"url"`
	Headers    http.Header `json:"headers"`
}

func echo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body := EchoBody{
		RemoteAddr: r.RemoteAddr,
		Method:     r.Method,
		URL:        r.URL.String(),
		Headers:    r.Header,
	}

	b, err := json.Marshal(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
