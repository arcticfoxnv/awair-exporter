package main

import (
	"awair-exporter/awair"
	"log"
	"net/http"
	"os"
)

func main() {
	client := awair.NewClient(os.Getenv("AWAIR_ACCESS_TOKEN"))
	e := NewExporterHTTP(client)
	m := http.NewServeMux()
	m.HandleFunc("/data/latest", e.serveLatest)
	m.HandleFunc("/meta/usage", e.serveUsage)
	s := &http.Server{Addr: ":8080", Handler: m}

	log.Println("Starting HTTP listener on", s.Addr)
	s.ListenAndServe()
}
