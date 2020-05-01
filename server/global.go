package server

import (
	"sync"

	"github.com/rs/zerolog"
)

// This file contains global variables and constants for the server package

// Global logger
var logger *zerolog.Logger

// mem holds clipd's memory
var mem = map[string]string{}

// Lock and Unlock mem
var memMut = &sync.RWMutex{}

// Default register name
const defaultRegister = "default"
