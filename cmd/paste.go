package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pasteReg string

func init() {
	rootCmd.AddCommand(pasteCmd)

	// Address
	pasteCmd.Flags().StringP("address", "a", ":8080", "address of the clipd server")
	viper.BindPFlag("client.address", pasteCmd.Flags().Lookup("address"))
	pasteCmd.Flags().StringVarP(&pasteReg, "register", "r", "", "named register")
}

var pasteCmd = &cobra.Command{
	Use:   "paste",
	Short: "Paste content from a register",
	Long: `
Paste content from a register. If a named register is not provided, the default
register ("default") is used.
	`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		clipdAddr := viper.GetString("client.address")
		url := fmt.Sprintf("http://%s/clipd", clipdAddr)

		if pasteReg != "" {
			url += fmt.Sprintf("/%s", pasteReg)
		}

		res, err := http.Get(
			url,
		)

		if err != nil {
			logger.Fatal(err)
		}

		if res != nil {
			defer res.Body.Close()
		}

		if res.StatusCode != http.StatusOK {
			logger.Printf("http error %d", res.StatusCode)
		}

		content, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logger.Fatal(err)
		}

		fmt.Fprint(os.Stdout, string(content))
	},
}
