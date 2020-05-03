package server

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog"
)

// DefaultRegister holds the name of the default register
const DefaultRegister = "default"

// NewLogger instantiates a new global logger
func NewLogger(develop bool, level zerolog.Level) *zerolog.Logger {
	var output io.Writer

	if develop {
		output = zerolog.ConsoleWriter{Out: os.Stdout}
	} else {
		output = os.Stdout
	}

	logger := zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Logger()

	return &logger
}

// Config holds the configuration for the clipd Server
type Config struct {
	Addr     string
	Develop  bool
	LogLevel zerolog.Level
}

// Server holds the clipd server
type Server struct {
	Config Config
	Logger *zerolog.Logger
	router http.Handler
}

// NewServer creates a new Server from Config, and initializes the global logger
func NewServer(config Config) *Server {
	// Setup logger
	lr := NewLogger(config.Develop, config.LogLevel)
	cache := cache.New(24*time.Hour, 10*24*time.Hour)

	// Create router
	router := httprouter.New()
	controller := NewClipdController(lr, cache)
	router.POST("/clipd", controller.Yank)
	router.GET("/clipd", controller.Paste)
	router.GET("/clipd/:reg", controller.Paste)

	requestLogger := RequestLogger{
		Next:   router,
		Logger: lr,
	}

	return &Server{
		Config: config,
		Logger: lr,
		router: requestLogger,
	}
}

// ListenAndServe starts the clipd server on the configured address
func (s *Server) ListenAndServe() error {
	s.Logger.Info().Msgf("server listening at %s", s.Config.Addr)
	return http.ListenAndServe(s.Config.Addr, s.router)
}
