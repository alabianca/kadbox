package cli

import (
	"context"
	"fmt"
	"github.com/alabianca/kadbox/core"
	"github.com/alabianca/kadbox/core/http"
	"github.com/alabianca/kadbox/core/kadprotocol"
	"github.com/alabianca/kadbox/core/node"
	"github.com/alabianca/kadbox/log"
	"github.com/libp2p/go-libp2p-core/network"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	secio "github.com/libp2p/go-libp2p-secio"
	"github.com/spf13/cobra"
	"net"
	"os"
	"strconv"
)

var (
	isGateway *bool
)

func init() {
	serverCmd.AddCommand(startCmd)
	isGateway = startCmd.Flags().BoolP("gateway", "g", false, "specify true if this node should be a gateway node")
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

	key, err := repo.Identity.GetPrivateKey()
	if err != nil {
		return printError(err)
	}

	// node options
	options := []node.Option{
			node.Gateways(repo.Gateways...),
			node.Identity(key),
			node.ListenAddr(repo.ListenAddrs...),
			node.Security(secio.ID, secio.New),
	}

	// if we are a gateway we also act as a relay
	// and allow peers to relay their connections through me
	if *isGateway {
		options = append(options, node.ActAsRelay(), node.Routing(ctx, dht.Mode(dht.ModeServer)))
	} else {
		options = append(options, node.EnableRelay(), node.StaticRelays(repo.Gateways...), node.Routing(ctx)) // use the gateways as relays
	}


	nde, err := node.New(ctx, options...)
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

	log.Infof("Daemon is listening at %s\n", server.Addr())

	go func() {
		nde.SetStreamHandler(func(stream network.Stream) {
			kadpService.HandleStream(stream)

		})
	}()

	if err := server.ListenAndServe(ctx); err != nil {
		return printError(err)
	}

	return 0

}
