package jewel

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"sync"
)

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type HandlerFunc func(http.ResponseWriter, *http.Request)

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

// Helper handlers

// Error replies to the request with the specified error message and HTTP code.
// It does not otherwise end the request; the caller should ensure no further
// writes are done to w.
// The error message should be plain text.
func Error(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintln(w, error)
}

// NotFound replies to the request with an HTTP 404 not found error.
func NotFound(w http.ResponseWriter, r *http.Request) {
	Error(w, "404 page not found", http.StatusNotFound)
}

// NotFoundHandler returns a simple request handler
// that replies to each request with a ``404 page not found'' reply.
func NotFoundHandler() http.Handler { return HandlerFunc(NotFound) }

// Redirect to a fixed URL
type redirectHandler struct {
	url  string
	code int
}

func (rh *redirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, rh.url, rh.code)
}

// RedirectHandler returns a request handler that redirects
// each request it receives to the given url using the given
// status code.
//
// The provided code should be in the 3xx range and is usually
// StatusMovedPermanently, StatusFound or StatusSeeOther.
func RedirectHandler(url string, code int) http.Handler {
	return &redirectHandler{url, code}
}

var DefaultJewelServeMux = &JewelServeMux{matcher: DefaultJewelMatcher}
var DefaultJewelMatcher = NewJewelMatchSystem()

type JewelServeMux struct {
	matcher PatternMatcherInterface
	mu      sync.RWMutex
	hosts   bool
}

type PatternMatcherInterface interface {
	AddPattern(pattern string, handler http.Handler)
	Match(url string) (handler http.Handler, pattern string)
}

func (jmx *JewelServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h, pattern := jmx.Handler(r)
	if pattern != "" {
		log.Printf("New HTTP request [%s] %s ==> %s", r.Method, r.RequestURI, pattern)
	}
	h.ServeHTTP(w, r)
}

// Handle registers the handler for the given pattern.
// If a handler already exists for pattern, Handle panics.
func (jmx *JewelServeMux) Handle(pattern string, handler http.Handler) {
	jmx.mu.Lock()
	defer jmx.mu.Unlock()

	if pattern == "" {
		panic("http: invalid pattern " + pattern)
	}
	if handler == nil {
		panic("http: nil handler")
	}

	jmx.matcher.AddPattern(pattern, handler)

	if pattern[0] != '/' {
		jmx.hosts = true
	}
}

// Handler is the main implementation of ServeHTTP.
func (jmx *JewelServeMux) Handler(r *http.Request) (h http.Handler, pattern string) {
	if r.Method != "CONNECT" {
		if p := cleanPath(r.URL.Path); p != r.URL.Path {
			_, pattern = jmx.handler(r.Host, p)
			url := *r.URL
			url.Path = p
			return RedirectHandler(url.String(), http.StatusMovedPermanently), pattern
		}
	}

	return jmx.handler(r.Host, r.URL.Path)
}

// handler is the main implementation of Handler.
// The path is known to be in canonical form, except for CONNECT methods.
func (jmx *JewelServeMux) handler(host, path string) (h http.Handler, pattern string) {
	jmx.mu.RLock()
	defer jmx.mu.RUnlock()

	// Host-specific pattern takes precedence over generic ones
	if jmx.hosts {
		h, pattern = jmx.matcher.Match(host + path)
	}
	if h == nil {
		h, pattern = jmx.matcher.Match(path)
	}
	if h == nil {
		h, pattern = NotFoundHandler(), ""
	}
	return
}

// Return the canonical path for p, eliminating . and .. elements.
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}
	return np
}

func Handle(pattern string, handler http.Handler) { DefaultJewelServeMux.Handle(pattern, handler) }
