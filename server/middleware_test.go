package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/nicolomaioli/clipd/server/internal/testutils"
	"github.com/rs/zerolog"
)

func TestRequestLogger_ServeHTTP(t *testing.T) {
	type args struct {
		status int
		method string
		path   string
	}

	// An inspectable bytes.Buffer Logger can write to
	outBuf := new(bytes.Buffer)
	lr := zerolog.New(outBuf)

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
			defer outBuf.Reset()

			spy := testutils.SpyHandler{Status: tt.args.status}
			rl := RequestLogger{
				Next:   spy,
				Logger: &lr,
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.args.method, tt.args.path, nil)

			rl.ServeHTTP(w, r)

			tle := LogEntry{}
			err := json.Unmarshal(outBuf.Bytes(), &tle)
			if err != nil {
				t.Fatalf("could not unmarshal %q: %s", outBuf.String(), err)
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

func TestResponseStats_Header(t *testing.T) {
	tests := []struct {
		name string
		r    *ResponseStats
		want http.Header
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Header(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResponseStats.Header() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponseStats_Write(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		r       *ResponseStats
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Write(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResponseStats.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ResponseStats.Write() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponseStats_WriteHeader(t *testing.T) {
	type args struct {
		statusCode int
	}
	tests := []struct {
		name string
		r    *ResponseStats
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.WriteHeader(tt.args.statusCode)
		})
	}
}
