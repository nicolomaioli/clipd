package server

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

var outBuf = new(bytes.Buffer)

func TestMain(m *testing.M) {
	lr := zerolog.New(outBuf)
	logger = &lr

	exit := m.Run()
	os.Exit(exit)
}

func TestYank(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "It updates the default register",
			body:           "{\"content\":\"test\"}",
			expectedStatus: 200,
		},
		{
			name:           "It updates the given register",
			body:           "{\"reg\":\"abc\",\"content\":\"test\"}",
			expectedStatus: 200,
		},
		{
			name:           "It returns 400 if the body is invalid json",
			body:           "invalid json",
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer outBuf.Reset()

			memMut.Lock()
			mem = make(map[string]string)
			memMut.Unlock()

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/clipd", bytes.NewReader([]byte(tt.body)))
			var p []httprouter.Param

			yank(w, r, p)
			if w.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var clip Clip
			_ = json.Unmarshal([]byte(tt.body), &clip)

			if clip.Reg == "" {
				clip.Reg = defaultRegister
			}

			var content string

			memMut.RLock()
			if v, ok := mem[clip.Reg]; ok {
				content = v
			}
			memMut.RUnlock()

			if clip.Content != content {
				t.Fatalf("expected content %q, got %q", clip.Content, content)
			}
		})
	}
}
