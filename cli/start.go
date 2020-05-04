package cli

import (
	"context"
	"fmt"
	"github.com/alabianca/kadbox/core"
	"github.com/alabianca/kadbox/core/http"
	"github.com/alabianca/kadbox/core/kadprotocol"
	"github.com/alabianca/kadbox/core/node"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/spf13/cobra"
	"net"
	"os"
	"strconv"
)

func init() {
	serverCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use: "start",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Start command")
		os.Exit(runStart())
	},
}

func runStart() int {
	home, err := os.UserHomeDir()
	if err != nil {
		return printError(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return printError(err)
	}

	repo, err := getConfigFile(wd, home)
	if err != nil {
		return printError(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nde, err := node.New(ctx, repo)
	if err != nil {
		return printError(err)
	}

	kadpService := kadprotocol.New()


	listenAddress := core.ListenAddress(net.JoinHostPort(repo.API.Address, strconv.Itoa(int(repo.API.Port))))

	storage := http.StorageService{
		Node: nde,
		Protocol: kadpService,
	}
	server := core.NewServer(core.AppContext{Node: nde}, &storage, listenAddress)

	fmt.Printf("Daemon is listening at %s\n", server.Addr())

	go func() {
		nde.SetStreamHandler(func(stream network.Stream) {
			kadpService.HandleStream(stream)

		})
	}()

	if err := server.ListenAndServe(); err != nil {
		return printError(err)
	}

	return 0

}


