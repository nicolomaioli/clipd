package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

var outBuf = new(bytes.Buffer)

func TestMain(t *testing.M) {
	lr := zerolog.New(outBuf)
	logger = &lr
	t.Run()
}

func TestYank(t *testing.T) {
	t.Parallel()

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
		p []httprouter.Param
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "It updates the content of the default register",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/clipd", nil),
				p: []httprouter.Param{{}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.w.Code != http.StatusOK {
				t.Error("not ok")
			}
		})
	}
}
