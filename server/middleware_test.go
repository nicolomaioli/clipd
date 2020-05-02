package server

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"

	"github.com/nicolomaioli/clipd/server/internal/testutils"
)

func TestRequestLogger_ServeHTTP(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}

	// An inspectable bytes.Buffer Logger can write to
	outBuf := new(bytes.Buffer)
	lr := zerolog.New(outBuf)

	tests := []struct {
		name string
		rl   RequestLogger
		args args
	}{
		{
			name: "It logs the incoming request and calls Next",
			rl: RequestLogger{
				Next:   testutils.SpyHandler{},
				Logger: &lr,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/copy/reg", nil),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer outBuf.Reset()

			tt.rl.ServeHTTP(tt.args.w, tt.args.r)
			want := fmt.Sprintf("{\"level\":\"info\",\"method\":\"%s\",\"url\":\"%s\"}\n", tt.args.r.Method, tt.args.r.URL.Path)
			got := outBuf.String()

			if got != want {
				t.Fatalf("expected %q, got %q", want, got)
			}

			if tt.args.w.Body.String() != "called" {
				t.Fatal("handler Next was not called")
			}
		})
	}
}
