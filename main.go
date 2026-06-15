package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

func init() {
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!\n", "Alan Turing")
	})
}

func main() {
	port := uint16(8080)
	if v, ok := os.LookupEnv("PORT"); ok {
		port64, err := strconv.ParseUint(v, 10, 16)
		if err != nil {
			log.Fatal(err)
		}
		port = uint16(port64)
	}

	l, err := net.Listen("tcp", ":"+strconv.FormatUint(uint64(port), 10))
	if err != nil {
		log.Fatal(err)
	}
	addr := l.Addr().(*net.TCPAddr)
	fmt.Printf("%d\n", addr.Port)

	err = http.Serve(l, nil)
	if err != nil {
		log.Fatal(err)
	}
}
