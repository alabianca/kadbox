package cli

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"strconv"
)

func init() {
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use: "get [hash]",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(runGet(args))
	},
}

func runGet(args []string) int {
	fileHash := args[0]

	repo, err := getClosestConfigFileRelativeToWd()
	if err != nil {
		return printError(err)
	}

	url := fmt.Sprintf("http://%s:%s/?key=%s", repo.API.Address, strconv.Itoa(int(repo.API.Port)), fileHash)

	res, err := http.Get(url)
	if err != nil {
		return printError(err)
	}

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, res.Body); err != nil {
		return printError(err)
	}

	if res.StatusCode != http.StatusOK {
		io.Copy(os.Stderr, buf)
	} else {
		gzr, err := gzip.NewReader(buf)
		if err != nil {
			return printError(err)
		}

		if err := readTarball(gzr, "."); err != nil {
			return printError(err)
		}
	}

	return 0

}
