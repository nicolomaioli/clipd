package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var yankReg string

func init() {
	rootCmd.AddCommand(yankCmd)

	// Address
	yankCmd.Flags().StringP("address", "a", ":8080", "address of the clipd server")
	viper.BindPFlag("client.address", yankCmd.Flags().Lookup("address"))
	yankCmd.Flags().StringVarP(&yankReg, "register", "r", "", "named register")
}

var yankCmd = &cobra.Command{
	Use:   "yank",
	Short: "Yank content to a register",
	Long: `
Yank content to a register. If a named register is not provided, the default
register ("default") is used. It is intended to work with pipes and i/o
redirection.
	`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		info, err := os.Stdin.Stat()
		if err != nil {
			logger.Fatal(err)
		}

		if (info.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
			logger.Fatal("yank is intended to work with pipes")
		}

		if info.Size() == 0 {
			os.Exit(0)
		}

		body := bufio.NewReader(os.Stdin)
		clipdAddr := viper.GetString("client.address")
		url := fmt.Sprintf("http://%s/clipd", clipdAddr)

		if yankReg != "" {
			url += fmt.Sprintf("/%s", yankReg)
		}

		res, err := http.Post(
			url,
			"text/plain",
			body,
		)

		if res != nil {
			defer res.Body.Close()
		}

		if err != nil {
			logger.Fatal(err)
		}

		if res.StatusCode != http.StatusOK {
			logger.Fatalf("http error %d", res.StatusCode)
		}
	},
}
