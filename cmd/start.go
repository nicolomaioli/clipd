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
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringP("address", "a", ":8080", "address of the clipd server")
	startCmd.Flags().BoolP("develop", "d", false, "set developer mode")
	startCmd.Flags().IntP("logLevel", "l", 3, `set log level (0-7)`)

	viper.BindPFlag("server.address", startCmd.Flags().Lookup("address"))
	viper.BindPFlag("server.develop", startCmd.Flags().Lookup("develop"))
	viper.BindPFlag("server.logLevel", startCmd.Flags().Lookup("logLevel"))
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the clipd server",
	Long: `
Start the clipd server. Logs are printed to Stdout and can be redirected.
	`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config := &server.Config{
			Addr:     viper.GetString("server.address"),
			Develop:  viper.GetBool("server.develop"),
			LogLevel: zerolog.Level(viper.GetInt("server.logLevel")),
		}

		s := server.NewClipdServer(config)

		// Add listners for shutdown
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		// Start the server in a separate goroutine
		go func() {
			err := s.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				log.Fatalf("ListenAndServe error: %s", err)
			}
		}()

		// Block until signal, then handle shutdown gracefully
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
