package middleware

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"testing"
)

type EmptyHandler struct{}

func (h EmptyHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
	return
}

type EmptyWriter struct{}

func (w EmptyWriter) Header() http.Header {
	return http.Header{}
}

func (w EmptyWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (w EmptyWriter) WriteHeader(int) {
	return
}

func TestRequestLogger_ServeHTTP(t *testing.T) {
	t.Parallel()

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}

	outBuf := new(bytes.Buffer)

	tests := []struct {
		name string
		rl   RequestLogger
		args args
	}{
		{
			name: "It logs the GET request",
			rl: RequestLogger{
				Next:   EmptyHandler{},
				Logger: log.New(outBuf, "", 0),
			},
			args: args{
				w: EmptyWriter{},
				r: &http.Request{
					Method: "GET",
					URL: &url.URL{
						Path: "/",
					},
				},
			},
		},
		{
			name: "It logs the POST request",
			rl: RequestLogger{
				Next:   EmptyHandler{},
				Logger: log.New(outBuf, "", 0),
			},
			args: args{
				w: EmptyWriter{},
				r: &http.Request{
					Method: "POST",
					URL: &url.URL{
						Path: "/clipd/reg",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.rl.ServeHTTP(tt.args.w, tt.args.r)
			want := fmt.Sprintf("%s %s\n", tt.args.r.Method, tt.args.r.URL.Path)
			got := outBuf.String()
			if got != want {
				t.Fatalf("expected %q, got %q", want, got)
			}
			outBuf.Reset()
		})
	}
}
