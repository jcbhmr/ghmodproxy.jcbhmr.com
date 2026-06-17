package main

import (
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/mod/module"
)

func init() {
	http.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		// https://inadarei.github.io/rfc-healthcheck/
		w.Header().Set("Content-Type", "application/health+json")
		fmt.Fprintln(w, `{"status":"pass"}`)
	})

	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintln(w, "TODO")
	})

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if before, after, found := strings.Cut(r.URL.Path, "/@v/"); found {
			pathEscaped := strings.TrimPrefix(before, "/")

			var handler func(w http.ResponseWriter, r *http.Request, v module.Version, pq parsedQuery)
			var versionEscaped string
			var found bool
			if versionEscaped, found = strings.CutSuffix(after, ".info"); found {
				handler = handleInfo
			} else if versionEscaped, found = strings.CutSuffix(after, ".mod"); found {
				handler = handleMod
			} else if versionEscaped, found = strings.CutSuffix(after, ".zip"); found {
				handler = handleZip
			} else {
				http.Error(w, fmt.Sprintf("%s does not match *.{info,mod,zip}", after), http.StatusBadRequest)
				return
			}

			pathValue, err := module.UnescapePath(pathEscaped)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			version, err := module.UnescapeVersion(versionEscaped)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err := module.Check(pathValue, version); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			pq, err := parseQuery(r.URL.Query(), pathValue)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			handler(w, r, module.Version{Path: pathValue, Version: version}, pq)
			return
		} else if before, found := strings.CutSuffix(r.URL.Path, "/@latest"); found {
			pathEscaped := strings.TrimPrefix(before, "/")

			pathValue, err := module.UnescapePath(pathEscaped)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			pq, err := parseQuery(r.URL.Query(), pathValue)

			handleLatest(w, r, pathValue, pq)
			return
		} else {
			http.Error(w, fmt.Sprintf("%s does not match */@v/* or */@latest", r.URL), http.StatusBadRequest)
			return
		}
	})
}

type parsedQuery struct {
	Owner  string
	Repo   string
	Prefix string
	Stem   string
}

var ownerRegexp = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9]|-(?=[a-zA-Z0-9])){0,38}$`)
var repoRegexp = regexp.MustCompile(`^[a-zA-Z0-9._-]{1,100}$`)

func parseQuery(q url.Values, pathValue string) (pq parsedQuery, err error) {
	owner := q.Get("owner")
	if owner == "" {
		return parsedQuery{}, fmt.Errorf("owner not present")
	}
	if !ownerRegexp.MatchString(pq.Owner) {
		return parsedQuery{}, fmt.Errorf("owner %s does not match %s: %w", owner, ownerRegexp, strconv.ErrSyntax)
	}

	repo := q.Get("repo")
	if repo == "" {
		return parsedQuery{}, fmt.Errorf("repo not present")
	}
	if !repoRegexp.MatchString(repo) {
		return parsedQuery{}, fmt.Errorf("repo %s does not match %s: %w", repo, repoRegexp, strconv.ErrSyntax)
	}

	prefix := q.Get("prefix")
	// Assume all prefixes are valid.

	stem := q.Get("asset")
	if stem == "" {
		stem = path.Base(pathValue)
	}
	if !fs.ValidPath(stem) || strings.Contains(stem, "/") {
		return parsedQuery{}, fmt.Errorf("stem %s is not a valid file name: %w", stem, strconv.ErrSyntax)
	}

	return parsedQuery{
		Owner:  owner,
		Repo:   repo,
		Prefix: prefix,
		Stem:   stem,
	}, nil
}

type info struct {
	Version string
	Time    time.Time `json:"omitzero"`
}

func handleInfo(w http.ResponseWriter, r *http.Request, v module.Version, pq parsedQuery) {
	location := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s%s/%s.info", url.PathEscape(pq.Owner), url.PathEscape(pq.Repo), url.PathEscape(pq.Prefix), url.PathEscape(v.Version), url.PathEscape(pq.Stem))
	http.Redirect(w, r, location, http.StatusTemporaryRedirect)
}

func handleMod(w http.ResponseWriter, r *http.Request, v module.Version, pq parsedQuery) {
	location := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s%s/%s.mod", url.PathEscape(pq.Owner), url.PathEscape(pq.Repo), url.PathEscape(pq.Prefix), url.PathEscape(v.Version), url.PathEscape(pq.Stem))
	http.Redirect(w, r, location, http.StatusTemporaryRedirect)
}

func handleZip(w http.ResponseWriter, r *http.Request, v module.Version, pq parsedQuery) {
	location := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s%s/%s.zip", url.PathEscape(pq.Owner), url.PathEscape(pq.Repo), url.PathEscape(pq.Prefix), url.PathEscape(v.Version), url.PathEscape(pq.Stem))
	http.Redirect(w, r, location, http.StatusTemporaryRedirect)
}

func handleLatest(w http.ResponseWriter, r *http.Request, pathValue string, pq parsedQuery) {
	http.Error(w, "not implemented", 404)
}

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
