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
		if got, want := os.Getenv("DENO_SERVE_ADDRESS"), "duplicate,vsock:-1:8080"; got != want {
			panic(fmt.Errorf("DENO_SERVE_ADDRESS got %v, want %v", got, want))
		}
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
	fmt.Printf("Hello %s!\n", "Alan Turing")
	log.Fatal(listenAndServe(":8080", nil))
}
