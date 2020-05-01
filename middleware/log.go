package middleware

import (
	"log"
	"net/http"
)

// RequestLogger is middleware that logs incoming requests
type RequestLogger struct {
	Next   http.Handler
	Logger *log.Logger
}

func (rl RequestLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rl.Logger.Printf("%s %s", r.Method, r.URL.Path)
	rl.Next.ServeHTTP(w, r)
}
