package main

import (
	"app/healthcheck"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/v88/github"
	"golang.org/x/mod/module"
)

//go:embed index.html
var FS embed.FS

func ServeFileFSHandler(fsys fs.FS, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, fsys, name)
	})
}

func NewServer(l *slog.Logger, gh *github.Client) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		j, err := json.Marshal(healthcheck.Response{Status: healthcheck.StatusPass})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", healthcheck.ContentType)
		w.Write(j)
	})
	mux.Handle("GET /{$}", ServeFileFSHandler(FS, "index.html"))
	mux.HandleFunc("GET /{owner}/{repo}/", func(w http.ResponseWriter, r *http.Request) {
		repo, err := ParseRepo(r.PathValue("owner"), r.PathValue("repo"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		q := r.URL.Query()
		subdirectory := q.Get("subdirectory")

		http.StripPrefix(fmt.Sprintf("/%s/%s", repo.Owner, repo.Repo), &MyModuleProxyHandler{
			GitHub:       gh,
			Repo:         repo,
			Subdirectory: subdirectory,
		}).ServeHTTP(w, r)
	})
	return &http.Server{
		ErrorLog:          slog.NewLogLogger(l.Handler(), slog.LevelError),
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
		Handler:           mux,
	}
}

func ParseSubdirectory(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	if !fs.ValidPath(s) {
		return "", fs.ErrInvalid
	}
	return s, nil
}

type Info struct {
	Version string
	Time    time.Time `json:",omitzero"`
	Origin  *Origin   `json:",omitempty"`
}

type Origin struct {
	VCS       string `json:",omitempty"`
	URL       string `json:",omitempty"`
	Subdir    string `json:",omitempty"`
	Hash      string `json:",omitempty"`
	TagPrefix string `json:",omitempty"`
	TagSum    string `json:",omitempty"`
	Ref       string `json:",omitempty"`
	RepoSum   string `json:",omitempty"`
}

func ParsePath(pathEscaped string) (string, error) {
	path, err := module.UnescapePath(pathEscaped)
	if err != nil {
		return "", err
	}
	if err := module.CheckPath(path); err != nil {
		return "", err
	}
	return path, nil
}

func ParsePathVersion(pathEscaped string, versionEscaped string) (module.Version, error) {
	path, pathErr := module.UnescapePath(pathEscaped)
	version, versionErr := module.UnescapeVersion(versionEscaped)
	err := errors.Join(pathErr, versionErr)
	if err != nil {
		return module.Version{}, err
	}
	if err := module.Check(path, version); err != nil {
		return module.Version{}, err
	}
	return module.Version{Path: path, Version: version}, nil
}

type MyModuleProxyHandler struct {
	GitHub       *github.Client
	Repo         Repo
	Subdirectory string
	initOnce     sync.Once
	mux          *http.ServeMux
}

func (m *MyModuleProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.initOnce.Do(m.init)
	m.mux.ServeHTTP(w, r)
}

func (m *MyModuleProxyHandler) init() {
	m.mux = http.NewServeMux()
	m.mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if before, after, found := strings.Cut(r.URL.Path, "/@v/"); found {
			pathEscaped := strings.TrimPrefix(before, "/")

			if after == "list" {
				pathValue, err := ParsePath(pathEscaped)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				m.List(w, r, pathValue)
				return
			} else {
				var handlerFunc func(w http.ResponseWriter, r *http.Request, version module.Version)
				var versionEscaped string
				var found bool
				if versionEscaped, found = strings.CutSuffix(after, ".info"); found {
					handlerFunc = m.Info
				} else if versionEscaped, found = strings.CutSuffix(after, ".mod"); found {
					handlerFunc = m.Mod
				} else if versionEscaped, found = strings.CutSuffix(after, ".zip"); found {
					handlerFunc = m.Zip
				} else {
					http.Error(w, fmt.Sprintf("%q does not match *.{info,mod,zip}", after), http.StatusBadRequest)
					return
				}

				version, err := ParsePathVersion(pathEscaped, versionEscaped)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}

				handlerFunc(w, r, version)
				return
			}
		} else if before, found := strings.CutSuffix(r.URL.Path, "/@latest"); found {
			pathEscaped := strings.TrimPrefix(before, "/")

			pathValue, err := ParsePath(pathEscaped)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			m.Latest(w, r, pathValue)
			return
		} else {
			http.Error(w, fmt.Sprintf("%q does not match **/@v/* or **/@latest", r.URL), http.StatusBadRequest)
			return
		}
	})
}

func (m *MyModuleProxyHandler) releaseTags(ctx context.Context) ([]string, error) {
	releaseTags := []string{}
	for r, err := range m.GitHub.Repositories.ListReleasesIter(ctx, m.Repo.Owner, m.Repo.Repo, nil) {
		if err != nil {
			return nil, err
		}
		releaseTags = append(releaseTags, *r.TagName)
	}
	return releaseTags, nil
}

func (m *MyModuleProxyHandler) origin(ctx context.Context, tag string) (Origin, error) {
	hash, _, err := m.GitHub.Repositories.GetCommitSHA1(ctx, m.Repo.Owner, m.Repo.Repo, tag, "")
	if err != nil {
		return Origin{}, err
	}

	return Origin{
		VCS:  "git",
		URL:  fmt.Sprintf("https://github.com/%s/%s", url.PathEscape(m.Repo.Owner), url.PathEscape(m.Repo.Repo)),
		Hash: hash,
		Ref:  fmt.Sprintf("refs/tags/%s", tag),
	}, nil
}

func (m *MyModuleProxyHandler) List(w http.ResponseWriter, r *http.Request, path string) {
	releaseTags, err := m.releaseTags(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prefix := ""
	if m.Subdirectory != "" {
		prefix = m.Subdirectory + "/"
	}

	versions := []string{}
	for _, t := range releaseTags {
		versions = append(versions, strings.TrimPrefix(t, prefix))
	}

	s := strings.Join(versions, "\n") + "\n"

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, s)
}

func (m *MyModuleProxyHandler) Latest(w http.ResponseWriter, r *http.Request, path string) {
	http.NotFound(w, r)
}

func (m *MyModuleProxyHandler) Info(w http.ResponseWriter, r *http.Request, version module.Version) {
	prefix := ""
	if m.Subdirectory != "" {
		prefix = m.Subdirectory + "/"
	}

	origin, err := m.origin(r.Context(), version.Version)

	b, err := json.Marshal(Info{
		Version: "",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (m *MyModuleProxyHandler) Mod(w http.ResponseWriter, r *http.Request, version module.Version) {

}

func (m *MyModuleProxyHandler) Zip(w http.ResponseWriter, r *http.Request, version module.Version) {

}
