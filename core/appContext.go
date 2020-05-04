package core

import "github.com/libp2p/go-libp2p-core/peer"

const (
	Protocol = "kadbox"
	Scheme   = "kadbox://"
)

type AppContext struct {
	Node Node
}

func ProtocolKey(key string) string {
	return "/" + Protocol + "/" + key
}

func PeerIDFromBytes(b []byte) (peer.ID, error) {
	return peer.IDFromBytes(b)
}
