package main

import (
	"app/xhttp"
	"archive/zip"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path"
	"slices"
	"strings"
	"time"

	bufra "github.com/avvmoto/buf-readerat"
	"github.com/google/go-github/v88/github"
	"github.com/snabb/httpreaderat"
	"golang.org/x/mod/module"
	"golang.org/x/sync/singleflight"
)

//go:embed all:public
var publicEmbedFS embed.FS
var Public = func() fs.FS {
	fsys, err := fs.Sub(publicEmbedFS, "public")
	if err != nil {
		panic(err)
	}
	return fsys
}()

var Handler = func() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServerFS(Public))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(map[string]any{"status": "pass"})
		if err != nil {
			xhttp.InternalServerError(w, r, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/health+json")
		w.Write(b)
	})
	mux.Handle("GET /{owner}/{repo}/", func() http.Handler {
		mux := xhttp.NewRegexpMux()
		mux.HandleCompileFunc("GET", `^/(?:[^/]*)/(?:[^/]*)/(?P<escapedPath>.*?)/@v/list$`, list)
		mux.HandleCompileFunc("GET", `^/(?:[^/]*)/(?:[^/]*)/(?P<escapedPath>.*?)/@latest$`, latest)
		mux.HandleCompileFunc("GET", `^/(?:[^/]*)/(?:[^/]*)/(?P<escapedPath>.*?)/@v/(?P<escapedVersion>.*?)(?P<ext>.info|mod|zip)$`, details)
		return mux
	}())
	return mux
}()

var getReleaseByTagGroup = &singleflight.Group{}

func getReleaseByTag(ctx context.Context, owner, repo, tag string) (*github.RepositoryRelease, error) {
	key := owner + "/" + repo + "#" + tag
	v, err, _ := getReleaseByTagGroup.Do(key, func() (any, error) {
		release, _, err := GitHubClient.Repositories.GetReleaseByTag(ctx, owner, repo, tag)
		return release, err
	})
	return v.(*github.RepositoryRelease), err
}

var listReleasesGroup = &singleflight.Group{}

func listReleases(ctx context.Context, owner, repo string) ([]*github.RepositoryRelease, error) {
	key := owner + "/" + repo
	v, err, _ := listReleasesGroup.Do(key, func() (any, error) {
		releases := []*github.RepositoryRelease{}
		var r *github.RepositoryRelease
		var err error
		for r, err = range GitHubClient.Repositories.ListReleasesIter(ctx, owner, repo, nil) {
			if r != nil {
				releases = append(releases, r)
			}
			if err != nil {
				break
			}
		}
		return releases, err
	})
	return v.([]*github.RepositoryRelease), err
}

type Info struct {
	Version string
	Time    time.Time `json:",omitzero"`
}

func list(w http.ResponseWriter, r *http.Request) {
	owner := r.PathValue("owner")
	repo := r.PathValue("repo")
	escapedPath := r.PathValue("escapedPath")
	q := r.URL.Query()
	subdirectory := q.Get("subdirectory")

	pathValue, err := module.UnescapePath(escapedPath)
	if err != nil {
		xhttp.BadRequest(w, r, err.Error())
		return
	}
	_ = pathValue // unused?

	var prefix string
	if subdirectory != "" {
		prefix = subdirectory + "/"
	}
	releases, err := listReleases(r.Context(), owner, repo)

	versions := make([]string, 0, len(releases))
	for _, r := range releases {
		v, found := strings.CutPrefix(r.GetTagName(), prefix)
		if !found {
			continue
		}

		// i := slices.IndexFunc(r.GetAssets(), func(a *github.ReleaseAsset) bool {
		// 	return a.GetName() == v+".zip"
		// })
		// if i < 0 {
		// 	continue
		// }

		versions = append(versions, v)
	}
	if len(versions) == 0 {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "%s\n", strings.Join(versions, "\n"))
}

func latest(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func details(w http.ResponseWriter, r *http.Request) {
	owner := r.PathValue("owner")
	repo := r.PathValue("repo")
	escapedPath := r.PathValue("escapedPath")
	escapedVersion := r.PathValue("escapedVersion")
	ext := r.PathValue("ext")
	q := r.URL.Query()
	subdirectory := q.Get("subdirectory")
	name := q.Get("name")
	if name == "" {
		xhttp.BadRequest(w, r, "no ?name= value")
		return
	}

	pathValue, pathErr := module.UnescapePath(escapedPath)
	version, versionErr := module.UnescapeVersion(escapedVersion)
	err := errors.Join(pathErr, versionErr)
	if err == nil {
		err = module.Check(pathValue, version)
	}
	if err != nil {
		xhttp.BadRequest(w, r, err.Error())
		return
	}

	var prefix string
	if subdirectory != "" {
		prefix = subdirectory + "/"
	}
	release, err := getReleaseByTag(r.Context(), owner, repo, prefix+version)
	if err != nil {
		xhttp.InternalServerError(w, r, err.Error())
		return
	}

	switch ext {
	case ".info":
		b, err := json.Marshal(Info{
			Version: version,
			Time:    release.GetPublishedAt().Time,
		})
		if err != nil {
			xhttp.InternalServerError(w, r, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	case ".mod", ".zip":
		i := slices.IndexFunc(release.GetAssets(), func(a *github.ReleaseAsset) bool {
			return a.GetName() == name
		})
		if i < 0 {
			http.Error(w, fmt.Sprintf("asset %q not found", name), http.StatusNotFound)
			return
		}
		asset := release.GetAssets()[i]

		u := asset.GetBrowserDownloadURL()
		if ext == ".zip" {
			http.Redirect(w, r, u, http.StatusTemporaryRedirect)
			return
		}

		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			panic(err)
		}

		httpReaderAt, err := httpreaderat.New(nil, req, nil)
		if err != nil {
			xhttp.InternalServerError(w, r, err.Error())
			return
		}
		bufferedReaderAt := bufra.NewBufReaderAt(httpReaderAt, 1_000_000)

		z, err := zip.NewReader(bufferedReaderAt, httpReaderAt.Size())
		if err != nil {
			xhttp.InternalServerError(w, r, err.Error())
			return
		}

		err = func() error {
			f, err := z.Open(path.Join(pathValue+"@"+version, "go.mod"))
			if err != nil {
				return err
			}
			defer f.Close()

			w.Header().Set("Content-Type", "text/plain;charset=utf-8")
			_, err = io.Copy(w, f)
			return err
		}()
		if err != nil {
			xhttp.InternalServerError(w, r, err.Error())
			return
		}
	default:
		panic(fmt.Sprintf("%q is not .(info|mod|zip)", ext))
	}
}
