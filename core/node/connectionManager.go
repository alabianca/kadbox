package node

import (
	"context"
	"fmt"
	"github.com/alabianca/kadbox/log"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/multiformats/go-multiaddr"
)

type connectionManager struct {
	host       host.Host
	relayAddrs []multiaddr.Multiaddr
}

// NewStream attempts to establish a connection to peer (info) without the relay at first.
// If it fails, we backoff and try to connect to the peer over the relay. If success we try to establish a stream
// and return it
func (c *connectionManager) NewStream(ctx context.Context, info peer.AddrInfo, protoID protocol.ID) (network.Stream, error) {
	// attempt to create a regular stream first
	log.Infof("Attempting Stream without relay to %s\n", info.ID.Pretty())
	var stream network.Stream
	var err error
	stream, err = c.host.NewStream(ctx, info.ID, protoID)
	if err == nil {
		// all good. we don't need the relay
		log.Debugf("No relay needed for %s\n", info.ID.Pretty())
		return stream, err
	}

	log.Infof("We need a relay for %s\n", info.ID.Pretty())
	// switch out the address info that did not work with our circuit addresses
	newInfo := c.infoWithRelayCircuitAddresses(info)
	// try to establish a connection to the peer over the relay
	if err := c.host.Connect(ctx, newInfo); err != nil {
		return nil, err
	}

	// connected. we should now be able to create a stream
	return c.host.NewStream(ctx, newInfo.ID, protoID)
}

func (c *connectionManager) infoWithRelayCircuitAddresses(oldInfo peer.AddrInfo) peer.AddrInfo {
	var info peer.AddrInfo
	info.ID = oldInfo.ID

	circuitAddrs := make([]multiaddr.Multiaddr, 0, len(c.relayAddrs)*len(oldInfo.Addrs))
	for _, ra := range c.relayAddrs {
		for _, addr := range oldInfo.Addrs {
			ma, err := multiaddr.NewMultiaddr(fmt.Sprintf("%s/p2p-circuit%s/p2p/%s", ra, addr, info.ID.Pretty()))
			if err == nil {
				circuitAddrs = append(circuitAddrs, ma)
			}

		}
	}

	info.Addrs = circuitAddrs

	return info
}
