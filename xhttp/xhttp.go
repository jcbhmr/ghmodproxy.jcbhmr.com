package xhttp

import (
	"io/fs"
	"net/http"
	"strings"
)

func MethodNotAllowed(w http.ResponseWriter, r *http.Request, allow []string) {
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
	}
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func MethodNotAllowedHandler(allow []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MethodNotAllowed(w, r, allow)
	})
}

func InternalServerError(w http.ResponseWriter, r *http.Request, error string) {
	http.Error(w, error, http.StatusInternalServerError)
}

func InternalServerErrorHandler(error string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		InternalServerError(w, r, error)
	})
}

func BadRequest(w http.ResponseWriter, r *http.Request, error string) {
	http.Error(w, error, http.StatusBadRequest)
}

func BadRequestHandler(error string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		BadRequest(w, r, error)
	})
}

func ServeFileFSHandler(fsys fs.FS, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, fsys, name)
	})
}
