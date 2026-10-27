package main

import (
	"log"
	"net/http"
	"observatorio.local/api/internal/health"
	"os"
	"time"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "healthcheck" {
		client := &http.Client{Timeout: 2 * time.Second}
		response, err := client.Get("http://127.0.0.1:8080/healthz")
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			os.Exit(1)
		}
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", health.Handler)
	server := &http.Server{Addr: ":8080", Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(server.ListenAndServe())
}
