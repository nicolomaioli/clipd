package main

import (
	"log"

	"github.com/nicolomaioli/clipd/server"
	"github.com/rs/zerolog"
)

func main() {
	config := server.Config{
		Addr:     ":8080",
		Develop:  true,
		LogLevel: zerolog.DebugLevel,
	}

	s := server.NewServer(config)

	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
