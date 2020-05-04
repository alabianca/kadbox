package core

import (
	"github.com/libp2p/go-libp2p-core/network"
	"io"
)

type KadProtocolService interface {
	HandleStream(stream network.Stream) KadProtocol
}

type KadProtocol interface {
	Want(key string) (io.Reader, chan error)
}
