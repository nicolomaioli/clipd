package cmd

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool
var logger *log.Logger

var rootCmd = &cobra.Command{
	Use:   "clipd",
	Short: "A simple clipboard with support for multiple registries ",
	Long: `
A simple clipboard with support for multiple registries. It provides an http
server, as well as paste and yank commands. Configuration can be passed with
flags, or using a config file (refer to the README for an example config file).
	`,
	Args: cobra.NoArgs,
}

// Execute is the main entry point of cmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.clipd.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "write cli logs to stdout")
}

func initConfig() {
	// Initialise client logger
	logger = log.New(ioutil.Discard, "", 0)

	if verbose {
		logger.SetOutput(os.Stdout)
	}

	// Read config file
	// If an error occurs, proceed without config
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := homedir.Dir()
		if err != nil {
			logger.Printf("homedir error: %s", err)
			return
		}

		// Search config in home directory with name ".clipd" (without extension)
		viper.AddConfigPath(home)
		viper.SetConfigName(".clipd")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err != nil {
		logger.Printf("read config error: %s", err)
	}
}
