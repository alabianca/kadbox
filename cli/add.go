package cli

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use: "add [file or directory]",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(runAdd(args))
	},
}

func runAdd(args []string) int {
	repo, err := getClosestConfigFileRelativeToWd()
	if err != nil {
		return printError(err)
	}
	// find out in which repo we need to store the file relative to wd
	repoDir, err := getClosestKadboxRepoRelativeToWd()
	if err != nil {
		return printError(err)
	}

	// 1. create the writers
	buf := new(bytes.Buffer) // everyting is written into this buffer befor sent to the daemon
	bodyWriter := multipart.NewWriter(buf)
	fileWriter, err := bodyWriter.CreateFormFile("upload", path.Join(repoDir, storeDirectoryName))
	if err != nil {
		return printError(err)
	}


	gzw := gzip.NewWriter(fileWriter)
	if _, err := writeTarball(gzw, args[0]); err != nil {
		return printError(err)
	}

	gzw.Close()

	// upload the file
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	url := fmt.Sprintf("http://%s:%s/storage", repo.API.Address, strconv.Itoa(int(repo.API.Port)))
	fmt.Printf("Upload size %d\n", buf.Len())
	res, err := http.Post(url, contentType, buf)
	if err != nil {
		return printError(err)
	}

	if _, err := io.Copy(os.Stdout, res.Body); err != nil {
		return printError(err)
	}



	fmt.Println()

	return 0
}