package server

import (
	"io"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

func newLogger(develop bool, level zerolog.Level) *zerolog.Logger {
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
	lr := newLogger(config.Develop, config.LogLevel)
	logger = lr

	// Create router
	router := httprouter.New()
	router.POST("/clipd", yank)
	router.GET("/clipd", paste)
	router.GET("/clipd/:reg", paste)

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
