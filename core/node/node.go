package node

import (
	"context"
	"errors"
	"fmt"
	"github.com/alabianca/kadbox/core"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	discovery "github.com/libp2p/go-libp2p-discovery"
	"github.com/libp2p/go-libp2p-kad-dht/dual"
	"github.com/multiformats/go-multiaddr"
	"sync"
)

type Node struct {
	// @TODO these gateways should probably be multiaddresses to make life easier
	Gateways         []string
	nodeContext      context.Context
	routingDiscovery *discovery.RoutingDiscovery
	host             host.Host
	dht              *dual.DHT
	// the connection manager handles the connection establishment between peers
	// it uses the Gateways as relays as a fallback in case no direct connection can be established
	connectionManager connectionManager
}

func New(ctx context.Context, opts ...Option) (*Node, error) {
	var n Node

	var p2pOpts []libp2p.Option
	for _, opt := range opts {
		p2pOpts = append(p2pOpts, opt(&n))
	}

	var err error
	n.nodeContext = ctx
	n.host, err = libp2p.New(
		ctx,
		p2pOpts...,
	)

	if err != nil {
		return nil, err
	}

	fmt.Println("My addresses are:")
	for _, addr := range n.host.Addrs() {
		fmt.Println(addr)
	}

	var mas []multiaddr.Multiaddr
	for _, g := range n.Gateways {
		ma, err := multiaddr.NewMultiaddr(g)
		if err == nil {
			mas = append(mas, ma)
		}
	}

	n.connectionManager = connectionManager{
		host: n.host,
		relayAddrs: mas,
	}

	return &n, nil
}

func (n *Node) SetStreamHandler(handler network.StreamHandler) {
	n.host.SetStreamHandler(core.Protocol, handler)
}

func (n *Node) Context() context.Context {
	return n.nodeContext
}

func (n *Node) Bootstrap(ctx context.Context) error {
	if err := n.dht.Bootstrap(ctx); err != nil {
		return err
	}

	if len(n.Gateways) == 0 {
		return nil
	}

	var errcs []chan error
	for _, addr := range n.Gateways {
		errcs = append(errcs, n.bootstrapConnect(ctx, addr))
	}

	merged := make(chan error)
	var wg sync.WaitGroup
	wg.Add(len(errcs))
	go func() {
		for _, c := range errcs {
			go func(errc chan error) {
				defer wg.Done()
				err := <-errc
				merged <- err
			}(c)
		}
	}()

	go func() {
		wg.Wait()
		close(merged)
	}()

	var nerr int
	for err := range merged {
		if err != nil {
			fmt.Printf("Gateway connection error: %s\n", err)
			nerr++
		}
	}

	if nerr == len(n.Gateways) {
		return errors.New("could not connect to any gateways")
	}

	return nil
}

func (n *Node) Advertise(key string) {
	fmt.Printf("Advertisin: %s\n", key)
	discovery.Advertise(n.nodeContext, n.routingDiscovery, key)
}

func (n *Node) FindPeers(ctx context.Context, key string) ([] peer.AddrInfo, error) {
	return discovery.FindPeers(ctx, n.routingDiscovery, key)
}

func (n *Node) NewStream(ctx context.Context, peerID peer.ID, protocols ...protocol.ID) (network.Stream, error) {
	return n.host.NewStream(ctx, peerID, protocols...)
}

func (n *Node) LocalPeerID() peer.ID {
	return n.host.ID()
}

func (n *Node) ConnectionManager() core.ConnectionManager {
	return &n.connectionManager
}

//func (n *Node) EnableAutoNATService(ctx context.Context, opts ...libp2p.Option) error {
//	fmt.Println("Starting AutoNAT")
//	nat, err := autonat.New(ctx, n.host)
//	//_, err := autonat.NewAutoNATService(ctx, n.host, opts...)
//	ma, _:= nat.PublicAddr()
//	fmt.Printf("Public Addr %s\n", ma)
//	fmt.Printf("Status: %d\n", nat.Status())
//	return err
//}

func (n *Node) bootstrapConnect(ctx context.Context, addr string) chan error {
	out := make(chan error)
	go func() {
		ma, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			return
		}
		info, err := peer.AddrInfoFromP2pAddr(ma)
		if err != nil {
			return
		}

		out <- n.host.Connect(ctx, *info)
		fmt.Printf("Connected to bootstrap peer %s\n", ma)
	}()

	return out
}
