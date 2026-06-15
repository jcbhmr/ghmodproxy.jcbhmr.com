package main

import (
	"fmt"
	"log"
	"math"
	"net"
	"net/http"

	"github.com/mdlayher/vsock"
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

	cid, err := vsock.ContextID()
	fmt.Printf("cid=%v, err=%v\n", cid, err)
	var l net.Listener
	l, err = vsock.ListenContextID(math.MaxUint32, 8080, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = http.Serve(l, nil)
	if err != nil {
		log.Fatal(err)
	}
}
