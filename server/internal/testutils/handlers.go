package testutils

import "net/http"

// SpyHandler is a http.Handler that writes "called" to the http.ResponseWriter and returns
type SpyHandler struct {
	Status int
}

func (h SpyHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(h.Status)
	w.Write([]byte("called"))
}
