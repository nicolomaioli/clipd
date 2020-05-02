package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/patrickmn/go-cache"
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
			name:           "It returns 400 if the body contains invalid json",
			body:           "invalid json",
			expectedStatus: 400,
		},
	}

	memc = cache.New(1*time.Hour, 1*time.Hour)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				defer outBuf.Reset()
				defer memc.Flush()
			}()

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

			if v, ok := memc.Get(clip.Reg); ok {
				content = v.(string)
			}

			if clip.Content != content {
				t.Fatalf("expected content %q, got %q", clip.Content, content)
			}
		})
	}
}

func TestPaste(t *testing.T) {
	tests := []struct {
		name            string
		insertInMemory  bool
		expectedReg     string
		expectedContent string
		expectedStatus  int
	}{
		{
			name:            "It copies from the default register",
			insertInMemory:  true,
			expectedContent: "test",
			expectedStatus:  200,
		},
		{
			name:            "It copies from a given register",
			insertInMemory:  true,
			expectedReg:     "abc",
			expectedContent: "test",
			expectedStatus:  200,
		},
		{
			name:           "It returns 404 if the register is empty",
			expectedStatus: 404,
		},
		{
			name:            "It escapes malformed json correctly",
			insertInMemory:  true,
			expectedContent: "\"test\"},{\"another key\":\"another value\"}",
			expectedStatus:  200,
		},
	}

	memc = cache.New(1*time.Hour, 1*time.Hour)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				defer outBuf.Reset()
				defer memc.Flush()
			}()

			if tt.insertInMemory {
				if tt.expectedReg == "" {
					tt.expectedReg = defaultRegister
				}

				memc.Set(tt.expectedReg, tt.expectedContent, 0)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", fmt.Sprintf("/clipd/%s", tt.expectedReg), nil)
			p := []httprouter.Param{{Key: "reg", Value: tt.expectedReg}}

			paste(w, r, p)
			if w.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var clip Clip
			body, _ := ioutil.ReadAll(w.Body)
			_ = json.Unmarshal(body, &clip)

			t.Logf("register %q, content %q", clip.Reg, clip.Content)
			t.Logf("expected register %q, content %q", tt.expectedReg, tt.expectedContent)

			if clip.Reg != tt.expectedReg || clip.Content != tt.expectedContent {
				t.Fatalf("expected to find content %q in register %q, got %q", tt.expectedContent, tt.expectedReg, clip.Content)
			}
		})
	}
}
