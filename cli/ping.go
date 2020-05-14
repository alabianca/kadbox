package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

func init() {
	rootCmd.AddCommand(pingCmd)
}

var pingCmd = &cobra.Command{
	Use: "ping",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(runPing())
	},
}

func runPing() int {
	repo, err := getClosestConfigFileRelativeToWd()
	if err != nil {
		return printError(err)
	}

	url := fmt.Sprintf("http://%s:%d/ping", repo.API.Address, repo.API.Port)

	res, err := http.Get(url)
	if err != nil {
		return printError(err)
	}

	if res.StatusCode == http.StatusOK {
		fmt.Println("Server is up")
	} else {
		fmt.Println("Looks like we can't reach the server")
	}

	return 0
}
