package server

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

// RequestLogger is a middleware that logs incoming requests.
type RequestLogger struct {
	Next   http.Handler
	Logger *zerolog.Logger
}

// LogEntry holds the information written to the logger
type LogEntry struct {
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	StartTime time.Time `json:"receivedTime"`
	Latency   int64     `json:"latencyMs"`
	Status    int       `json:"statusCode"`
}

// ResponseStats is a wrapper around http.ResponseWriter to hold information for logging
type ResponseStats struct {
	w      http.ResponseWriter
	status int
}

// Header implements http.ResponseWriter.Header
func (r *ResponseStats) Header() http.Header {
	return r.w.Header()
}

// Write implements http.ResponseWriter.Write
func (r *ResponseStats) Write(b []byte) (int, error) {
	return r.w.Write(b)
}

// WriteHeader implements http.ResponseWriter.WriteHeader
func (r *ResponseStats) WriteHeader(statusCode int) {
	r.status = statusCode
	r.w.WriteHeader(statusCode)
}

func (rl RequestLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	rs := &ResponseStats{w: w}
	rl.Next.ServeHTTP(rs, r)

	le := &LogEntry{
		Method:    r.Method,
		Path:      r.URL.Path,
		StartTime: start,
		Latency:   time.Since(start).Milliseconds(),
		Status:    rs.status,
	}

	rl.Logger.
		Info().
		Str("method", le.Method).
		Str("path", le.Path).
		Time("receivedTime", le.StartTime).
		Int64("latencyMs", le.Latency).
		Int("statusCode", le.Status).
		Send()
}
