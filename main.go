package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"

	"github.com/mdlayher/vsock"
	"golang.org/x/sync/errgroup"
)

func init() {
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!\n", "Alan Turing")
	})
}

func main() {
	// port := uint16(8080)
	// if v, ok := os.LookupEnv("PORT"); ok {
	// 	port64, err := strconv.ParseUint(v, 10, 16)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	port = uint16(port64)
	// }

	// l, err := net.Listen("tcp", ":"+strconv.FormatUint(uint64(port), 10))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// addr := l.Addr().(*net.TCPAddr)
	// fmt.Printf("%d\n", addr.Port)

	// err = http.Serve(l, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	g, _ := errgroup.WithContext(context.TODO())
	g.Go(func() error {
		var l net.Listener
		l, err := vsock.ListenContextID(math.MaxUint32, 8080, nil)
		if err != nil {
			return err
		}

		addr := l.Addr().(*vsock.Addr)
		fmt.Printf("Listening on http://%s\n", addr)

		return http.Serve(l, nil)
	})
	// g.Go(func() error {
	// 	l, err := net.Listen("tcp", ":8000")
	// 	if err != nil {
	// 		if errors.Is(err, )
	// 		return err
	// 	}

	// 	addr := l.Addr().(*net.TCPAddr)
	// 	fmt.Printf("Listening on http://%s\n", addr)

	// 	return http.Serve(l, nil)
	// })
	err := g.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
