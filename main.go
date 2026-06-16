package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math"
	"net"
	"net/http"

	"github.com/mdlayher/vsock"
	"golang.org/x/sync/errgroup"
)

func init() {
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("GET /", "Method", r.Method, "URL", r.URL)
		fmt.Fprintf(w, "Hello, %s!\n", "Alan Turing")
	})
}

func main() {
	g, ctx := errgroup.WithContext(context.TODO())
	g.Go(func() error {
		l, err := vsock.ListenContextID(math.MaxUint32, 8080, nil)
		if err != nil {
			return err
		}

		addr := l.Addr().(*vsock.Addr)
		fmt.Printf("Listening for HTTP on %s\n", addr)

		return http.Serve(l, nil)
	})
	g.Go(func() error {
		l, err := net.Listen("tcp", ":0")
		if err != nil {
			return err
		}

		addr := l.Addr().(*net.TCPAddr)
		fmt.Printf("Listening for HTTP on %s\n", addr)

		return http.Serve(l, nil)
	})
	_ = ctx
	err := g.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
