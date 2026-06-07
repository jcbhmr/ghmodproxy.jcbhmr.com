package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/mdlayher/vsock"
)

func listenAndServe(addr string, handler http.Handler) error {
	if v, _ := strconv.ParseBool(os.Getenv("DENO_DEPLOY")); v {
		// DENO_SERVE_ADDRESS is always "duplicate,vsock:-1:8080" on Deno Deploy.
		l, err := vsock.Listen(8080, nil)
		if err != nil {
			return err
		}
		return http.Serve(l, handler)
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Printf("Listening on http://%s/", l.Addr())
	return http.Serve(l, handler)
}

func init() {
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!\n", "Alan Turing")
	})
}

func main() {
	log.Fatal(listenAndServe(":8080", nil))
}
