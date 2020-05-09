package core

import (
	"context"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
)

type ConnectionManager interface {
	NewStream(ctx context.Context, info peer.AddrInfo, id protocol.ID) (network.Stream, error)
}
