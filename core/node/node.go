package node

import (
	"context"
	"errors"
	"fmt"
	"github.com/alabianca/kadbox/core"
	"github.com/libp2p/go-libp2p"
	autonat "github.com/libp2p/go-libp2p-autonat-svc"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/libp2p/go-libp2p-kad-dht/dual"
	"github.com/multiformats/go-multiaddr"
	"sync"
)

type Node struct {
	Gateways []string
	host     host.Host
	dht      *dual.DHT
}

func New(ctx context.Context, opts ...Option) (*Node, error) {
	var n Node

	var p2pOpts []libp2p.Option
	for _, opt := range opts {
		p2pOpts = append(p2pOpts, opt(&n))
	}

	var err error
	n.host, err = libp2p.New(
		ctx,
		p2pOpts...,
	)

	if err != nil {
		return nil, err
	}

	return &n, err
}

func (n *Node) SetStreamHandler(handler network.StreamHandler) {
	n.host.SetStreamHandler(core.Protocol, handler)
}

func (n *Node) Bootstrap(ctx context.Context) error {
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

func (n *Node) PutValue(ctx context.Context, key string, value []byte) error {
	return n.dht.PutValue(ctx, key, value)
}

func (n *Node) GetValue(ctx context.Context, key string) ([]byte, error) {
	return n.dht.GetValue(ctx, key)
}

func (n *Node) NewStream(ctx context.Context, peerID peer.ID, protocols ...protocol.ID) (network.Stream, error) {
	return n.host.NewStream(ctx, peerID, protocols...)
}

func (n *Node) LocalPeerID() peer.ID {
	return n.host.ID()
}

func (n *Node) EnableAutoNATService(ctx context.Context, opts ...libp2p.Option) error {
	_, err := autonat.NewAutoNATService(ctx, n.host, opts...)
	return err
}

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

