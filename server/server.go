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

// ClipdServer holds the clipd server
type ClipdServer struct {
	Logger *zerolog.Logger
	Cache  *cache.Cache
	Server *http.Server
}

// NewClipdServer creates a new Server from Config, and initializes the global logger
func NewClipdServer(config *Config) *ClipdServer {
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

	return &ClipdServer{
		Logger: lr,
		Cache:  cache,
		Server: &http.Server{
			Addr:         config.Addr,
			Handler:      requestLogger,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}
}

// CleanAfterShutdown can be called after http.Server.Shutdown
func (c *ClipdServer) CleanAfterShutdown() {
	c.Logger.Info().Msg("ClipdServer clean")
	c.Cache.Flush()
	c.Logger.Info().Msg("server shutdown complete")
}

// ListenAndServe starts the clipd server on the configured address
func (c *ClipdServer) ListenAndServe() error {
	c.Logger.Info().Str("addr", c.Server.Addr).Msg("server listening")
	return c.Server.ListenAndServe()
}
