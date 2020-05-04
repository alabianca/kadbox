package core

import (
	"context"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
)

type Node interface {
	Bootstrap() error
}

type NodeClient interface {
	LocalPeerID() peer.ID
	PutValue(ctx context.Context, key string, value []byte) error
	GetValue(ctx context.Context, key string) ([]byte, error)
	NewStream(ctx context.Context, id peer.ID, protocols ...protocol.ID) (network.Stream, error)
}
