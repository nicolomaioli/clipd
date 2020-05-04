package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nicolomaioli/clipd/server"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the clipd server",
	Long:  `Start the clipd server`,
	Run: func(cmd *cobra.Command, args []string) {
		config := server.Config{
			Addr:     ":8080",
			Develop:  true,
			LogLevel: zerolog.DebugLevel,
		}

		s := server.NewClipdServer(config)

		// Handle graceful shutdown
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			err := s.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				log.Fatalf("ListenAndServe error: %s", err)
			}
		}()

		<-c
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer func() {
			s.CleanAfterShutdown()
			cancel()
		}()

		if err := s.Server.Shutdown(ctx); err != nil {
			log.Fatalf("Server Shutdown Failed:%+v", err)
		}
		log.Print("Server Exited Properly")
	},
}
