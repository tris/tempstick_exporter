package main

import (
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultPort = "9974"
)

func main() {
	http.HandleFunc("/scrape", scrapeHandler)
	http.Handle("/metrics", promhttp.Handler()) // default Go metrics

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
