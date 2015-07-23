package httpex

import (
	"net/http"
	"sync"
)

type ServeMux struct {
	mu    sync.RWMutex
	m     map[string]muxEntry
	hosts bool
}

type muxEntry struct {
	explicit bool
	h        http.Handler
	pattern  string
}

func NewServeMux() *ServeMux { return &ServeMux{m: make(map[string]muxEntry)} }

var DefaultServeMux = NewServeMux()

func (mux *ServeMux) Handler(r *Request) (h http.Handler, pattern string) {
	if r.Method != "CONNECT" {
		if p := cleanPath(r.URL.Path); p != r.URL.Path {
			_, pattern = mux.handler(r.Host, p)
			url := *r.URL
			url.Path = p
			return RedirectHandler(url.String(), StatusMovedPermanently), pattern
		}
	}

	return mux.handler(r.Host, r.URL.Path)
}
