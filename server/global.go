package server

import (
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog"
)

// This file contains global variables and constants for the server package

// Global logger
var logger *zerolog.Logger

// memc holds clipd's memory
var memc = cache.New(24*time.Hour, 10*24*time.Hour)

// Default register name
const defaultRegister = "default"
