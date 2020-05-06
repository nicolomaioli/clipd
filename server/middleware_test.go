package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
)

type SpyHandler struct {
	Status int
}

func (h SpyHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(h.Status)
	w.Write([]byte("called"))
}

func TestRequestLogger_ServeHTTP(t *testing.T) {
	type args struct {
		status int
		method string
		path   string
	}

	testOutBuf := new(bytes.Buffer)
	lr := zerolog.New(testOutBuf)

	tests := []struct {
		name string
		args args
	}{
		{
			name: "It calls Next and logs the correct status 201",
			args: args{
				status: 201,
				method: "POST",
				path:   "/test/route",
			},
		},
		{
			name: "It calls Next and logs the correct status 404",
			args: args{
				status: 404,
				method: "GET",
				path:   "/test/route",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer testOutBuf.Reset()

			spy := SpyHandler{Status: tt.args.status}
			rl := RequestLogger{
				Next:   spy,
				Logger: &lr,
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.args.method, tt.args.path, nil)

			rl.ServeHTTP(w, r)

			tle := LogEntry{}
			err := json.Unmarshal(testOutBuf.Bytes(), &tle)
			if err != nil {
				t.Fatalf("could not unmarshal %q: %s", testOutBuf.String(), err)
			}

			if tt.args.status != tle.Status {
				t.Errorf("expected logged status %d, got %d", tt.args.status, tle.Status)
			}

			if tt.args.method != tle.Method {
				t.Errorf("expected logged method %q, got %q", tt.args.method, tle.Method)
			}

			if tt.args.path != tle.Path {
				t.Errorf("expected logged path %q, got %q", tt.args.path, tle.Path)
			}

			if w.Body.String() != "called" {
				t.Errorf("handler Next was not called")
			}
		})
	}
}
