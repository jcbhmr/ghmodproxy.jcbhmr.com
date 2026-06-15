package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sys/unix"
)

func init() {
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Hello stdout!\n")
		log.Printf("Hello log!")
		slog.Info("GET /", "r.Method", r.Method, "r.URL", r.URL, "r.Header", r.Header)
		fmt.Fprintf(w, "Hello, %s!\n", "Alan Turing")
	})
}

func main() {
	g, _ := errgroup.WithContext(context.TODO())
	g.Go(func() error {
		l, err := net.Listen("tcp", os.Getenv("HOST")+":"+os.Getenv("PORT"))
		if err != nil {
			return err
		}

		addr := l.Addr().(*net.TCPAddr)
		fmt.Printf("Listening on http://%s\n", addr)

		p, err := os.FindProcess(os.Getppid())
		if err != nil {
			log.Fatal(err)
		}
		err = p.Signal(unix.SIGUSR2)
		if err != nil {
			log.Fatal(err)
		}

		return http.Serve(l, nil)
	})
	err := g.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
