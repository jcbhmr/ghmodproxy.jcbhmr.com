package main

// https://inadarei.github.io/rfc-healthcheck/

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/health+json")
		w.WriteHeader(200)
		fmt.Fprintln(w, `{"status":"pass"}`)
	})
}
