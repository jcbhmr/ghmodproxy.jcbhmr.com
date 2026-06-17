package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}
	l, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Listening on %s:%s\n", l.Addr().Network(), l.Addr())

	log.Fatal(http.Serve(l, nil))
}
