package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cgi"
)

func init() {
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!\n", "Alan Turing")
	})
}

func main() {
	// port := uint16(7045)
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
	// log.Fatal(fcgi.Serve(l, nil))
	err := cgi.Serve(nil)
	if err != nil {
		log.Fatal(err)
	}
	// log.Fatal(http.ListenAndServe(":"+strconv.FormatUint(uint64(port), 10), nil))
}
