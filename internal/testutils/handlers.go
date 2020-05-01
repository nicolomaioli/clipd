package testutils

import "net/http"

// SpyHandler is a http.Handler that writes "called" to the http.ResponseWriter and returns
type SpyHandler struct{}

func (h SpyHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("called"))
}
