package cmd

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var reg string

func init() {
	rootCmd.AddCommand(yankCmd)

	// Address
	yankCmd.Flags().StringP("address", "a", ":8080", "address of the clipd server")
	viper.BindPFlag("client.address", yankCmd.Flags().Lookup("address"))
	yankCmd.Flags().StringVarP(&reg, "register", "r", "", "named register")
}

var yankCmd = &cobra.Command{
	Use:   "yank",
	Short: "Yank content to a register",
	Long: `
Yank content to a register. If a named register is not provided, the default
register ("default") is used.
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		content := args[0]

		body := bytes.NewBuffer([]byte(content))
		clipdAddr := viper.GetString("client.address")
		url := fmt.Sprintf("http://%s/clipd", clipdAddr)

		if reg != "" {
			url += fmt.Sprintf("/%s", reg)
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
			log.Fatal(err)
		}

		if res.StatusCode != http.StatusOK {
			log.Fatalf("http error %d", res.StatusCode)
		}
	},
}
