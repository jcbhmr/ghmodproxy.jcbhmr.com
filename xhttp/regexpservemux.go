package xhttp

import (
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strconv"
)

type RegexpMux struct {
	routes []*regexpMuxRoute
}

func NewRegexpMux() *RegexpMux {
	return &RegexpMux{
		routes: []*regexpMuxRoute{},
	}
}

func (mux *RegexpMux) Handle(method string, pattern *regexp.Regexp, handler http.Handler) {
	mux.routes = append(mux.routes, &regexpMuxRoute{
		Method:  method,
		Pattern: pattern,
		Handler: handler,
	})
}

func (mux *RegexpMux) HandleFunc(method string, pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	mux.Handle(method, pattern, http.HandlerFunc(handler))
}

func (mux *RegexpMux) HandleCompile(method string, pattern string, handler http.Handler) {
	mux.Handle(method, regexp.MustCompile(pattern), handler)
}

func (mux *RegexpMux) HandleCompileFunc(method string, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mux.Handle(method, regexp.MustCompile(pattern), http.HandlerFunc(handler))
}

func (mux *RegexpMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	patternMatches := []*regexpMuxRoute{}
	for _, route := range mux.routes {
		if route.Pattern.MatchString(r.URL.Path) {
			patternMatches = append(patternMatches, route)
		}
	}
	if len(patternMatches) == 0 {
		http.NotFound(w, r)
		return
	}

	allow := []string{}
	patternAndMethodMatches := []*regexpMuxRoute{}
	for _, route := range patternMatches {
		if route.Method != "" {
			methods := []string{route.Method}
			if route.Method == http.MethodGet {
				methods = append(methods, http.MethodHead)
			}
			if slices.Contains(methods, r.Method) {
				allow = append(allow, methods...)
			} else {
				continue
			}
		}
		patternAndMethodMatches = append(patternAndMethodMatches, route)
	}

	if len(patternAndMethodMatches) == 0 {
		MethodNotAllowed(w, r, allow)
		return
	}
	route := patternAndMethodMatches[0]
	if len(patternAndMethodMatches) > 1 {
		panic(fmt.Errorf(">1 matches %#v for %s %q", patternAndMethodMatches, r.Method, r.URL.Path))
	}

	match := route.Pattern.FindStringSubmatch(r.URL.Path)
	for i, name := range route.Pattern.SubexpNames() {
		r.SetPathValue(strconv.FormatInt(int64(i), 10), match[i])
		if name != "" {
			r.SetPathValue(name, match[i])
		}
	}

	route.Handler.ServeHTTP(w, r)
}

type regexpMuxRoute struct {
	Method  string
	Pattern *regexp.Regexp
	Handler http.Handler
}
