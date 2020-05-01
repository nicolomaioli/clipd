package middleware

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nicolomaioli/clipd/internal/testutils"
)

func TestRequestLogger_ServeHTTP(t *testing.T) {
	t.Parallel()

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}

	// An inspectable bytes.Buffer Logger can write to
	outBuf := new(bytes.Buffer)

	tests := []struct {
		name string
		rl   RequestLogger
		args args
	}{
		{
			name: "It logs the incoming request and calls Next",
			rl: RequestLogger{
				Next:   testutils.SpyHandler{},
				Logger: log.New(outBuf, "", 0),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/copy/reg", nil),
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

			if tt.args.w.Body.String() != "called" {
				t.Fatal("handler Next was not called")
			}

			// Reset buffer for next test
			outBuf.Reset()
		})
	}
}
