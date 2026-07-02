package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/google/go-github/v88/github"
)

var Port uint16
var GitHubClient *github.Client

func Parse() {
	if s := os.Getenv("PORT"); s != "" {
		port64, err := strconv.ParseUint(s, 0, 16)
		Port = uint16(port64)
		if err != nil {
			log.Fatalf("PORT=%q invalid: %v", s, err)
		}
	} else {
		Port = 80
	}

	if s := os.Getenv("GITHUB_TOKEN"); s != "" {
		var err error
		GitHubClient, err = github.NewClient(github.WithAuthToken(s))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		var err error
		GitHubClient, err = github.NewClient()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	Parse()
	l, err := net.Listen("tcp", ":"+strconv.FormatUint(uint64(Port), 10))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening on %s", l.Addr())
	log.Fatal(http.Serve(l, nil))
}
