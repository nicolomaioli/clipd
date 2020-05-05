package server

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog"
)

type badReader int

func (badReader) Read([]byte) (int, error) {
	return 0, errors.New("test error")
}

func TestNewClipdController(t *testing.T) {
	testOutBuf := new(bytes.Buffer)
	lr := zerolog.New(testOutBuf)
	cc := cache.New(1*time.Millisecond, 1*time.Millisecond)

	type args struct {
		l *zerolog.Logger
		c *cache.Cache
	}
	tests := []struct {
		name string
		args args
		want *ClipdController
	}{
		{
			name: "It instantiates correctly",
			args: args{
				l: &lr,
				c: cc,
			},
			want: &ClipdController{
				logger: &lr,
				cache:  cc,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClipdController(tt.args.l, tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClipdController() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClipdController_Yank(t *testing.T) {
	testOutBuf := new(bytes.Buffer)
	lr := zerolog.New(testOutBuf)
	cc := cache.New(1*time.Minute, 10*time.Minute)
	controller := NewClipdController(&lr, cc)

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
		p httprouter.Params
	}

	type wants struct {
		statusCode int
		readCache  bool
		register   string
		content    []byte
	}

	tests := []struct {
		name  string
		c     *ClipdController
		args  args
		wants wants
	}{
		{
			name: "It yanks to the default register",
			c:    controller,
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/clipd", bytes.NewBuffer([]byte("test content"))),
				p: []httprouter.Param{{}},
			},
			wants: wants{
				statusCode: 200,
				readCache:  true,
				register:   DefaultRegister,
				content:    []byte("test content"),
			},
		},
		{
			name: "It yanks to a named register",
			c:    controller,
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/clipd/abc", bytes.NewBuffer([]byte("test content"))),
				p: []httprouter.Param{{
					Key:   "reg",
					Value: "abc",
				}},
			},
			wants: wants{
				statusCode: 200,
				readCache:  true,
				register:   "abc",
				content:    []byte("test content"),
			},
		},
		{
			name: "It returns 500 if reading the buffer errors out",
			c:    controller,
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/clipd", badReader(0)),
				p: []httprouter.Param{{}},
			},
			wants: wants{
				statusCode: 500,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				testOutBuf.Reset()
				cc.Flush()
			}()

			tt.c.Yank(tt.args.w, tt.args.r, tt.args.p)

			if tt.args.w.Code != tt.wants.statusCode {
				t.Errorf("Expected status %d, got %d", tt.wants.statusCode, tt.args.w.Code)
			}

			if tt.wants.readCache {
				if v, ok := cc.Get(tt.wants.register); !ok {
					t.Errorf("Expected content in register %q", tt.wants.register)
				} else {
					vByte := v.([]byte)

					if eq := bytes.Compare(vByte, tt.wants.content); eq != 0 {
						t.Errorf("Expected default register to contain %q, got %q", string(tt.wants.content), string(vByte))
					}
				}
			}
		})
	}
}

func TestClipdController_Paste(t *testing.T) {
	testOutBuf := new(bytes.Buffer)
	lr := zerolog.New(testOutBuf)
	cc := cache.New(1*time.Minute, 10*time.Minute)
	controller := NewClipdController(&lr, cc)

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
		p httprouter.Params
	}

	type wants struct {
		statusCode int
		writeCache bool
		register   string
		content    []byte
	}

	tests := []struct {
		name  string
		c     *ClipdController
		args  args
		wants wants
	}{
		{
			name: "It pastes from the default register",
			c:    controller,
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/clipd", nil),
				p: []httprouter.Param{{}},
			},
			wants: wants{
				statusCode: 200,
				writeCache: true,
				register:   DefaultRegister,
				content:    []byte("test content"),
			},
		},
		{
			name: "It pastes from a named register",
			c:    controller,
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/clipd/abc", nil),
				p: []httprouter.Param{{
					Key:   "reg",
					Value: "abc",
				}},
			},
			wants: wants{
				statusCode: 200,
				writeCache: true,
				register:   "abc",
				content:    []byte("test content"),
			},
		},
		{
			name: "It pastes from a named register",
			c:    controller,
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/clipd/abc", nil),
				p: []httprouter.Param{{
					Key:   "reg",
					Value: "abc",
				}},
			},
			wants: wants{
				statusCode: 404,
				register:   "abc",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				testOutBuf.Reset()
				cc.Flush()
			}()

			if tt.wants.writeCache {
				cc.Set(tt.wants.register, tt.wants.content, 0)
			}

			tt.c.Paste(tt.args.w, tt.args.r, tt.args.p)

			if tt.args.w.Code != tt.wants.statusCode {
				t.Errorf("Expected status %d, got %d", tt.wants.statusCode, tt.args.w.Code)
			}

			if eq := bytes.Compare(tt.args.w.Body.Bytes(), tt.wants.content); eq != 0 {
				t.Errorf("Expected response body to be %q, got %q", string(tt.wants.content), tt.args.w.Body.String())
			}
		})
	}
}
