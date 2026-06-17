package main


import (
	"fmt"
	"log/slog"
	"net/http"
)

func init() {
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("GET /", "Method", r.Method, "URL", r.URL)
		fmt.Fprintf(w, "Hello, %s!\n", "Alan Turing")
	})
}
