package server

import (
	"net/http"

	"github.com/rs/zerolog"
)

// RequestLogger is a middleware that logs incoming requests.
type RequestLogger struct {
	Next   http.Handler
	Logger *zerolog.Logger
}

func (rl RequestLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rl.Logger.Info().Str("method", r.Method).Str("url", r.URL.Path).Send()
	rl.Next.ServeHTTP(w, r)
}
