package core

import (
	"context"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
)

type Node interface {
	Bootstrap(ctx context.Context) error
}

type NodeClient interface {
	LocalPeerID() peer.ID
	Advertise(key string)
	FindPeers(key string) (<-chan peer.AddrInfo, error)
	NewStream(ctx context.Context, id peer.ID, protocols ...protocol.ID) (network.Stream, error)
}
