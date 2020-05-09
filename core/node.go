package core

import (
	"context"
	"github.com/libp2p/go-libp2p-core/peer"
)

type Node interface {
	Bootstrap(ctx context.Context) error
}

type NodeClient interface {
	LocalPeerID() peer.ID
	Advertise(key string)
	FindPeers(ctx context.Context, key string) ([]peer.AddrInfo, error)
	ConnectionManager() ConnectionManager
}
