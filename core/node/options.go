package node

import (
	"context"
	"github.com/alabianca/kadbox/core"
	"github.com/libp2p/go-libp2p"
	relay "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/routing"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-kad-dht/dual"
)

func noopConfig(cfg *libp2p.Config) error {
	return nil
}

type Option func(n *Node) libp2p.Option

func Routing(ctx context.Context) Option {
	return func(n *Node) libp2p.Option {
		return libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			var err error
			n.dht, err = dual.New(
				ctx,
				h,
				dht.NamespacedValidator(core.Protocol, &NullValidator{}),
				dht.ProtocolPrefix(core.Protocol),)

			// also set up the node's routing discovery
			n.routingDiscovery = discovery.NewRoutingDiscovery(n.dht)

			return n.dht, err
		})
	}
}

func Security(name string, tpt interface{}) Option {
	return func(n *Node) libp2p.Option {
		return libp2p.Security(name, tpt)
	}
}

func Identity(key crypto.PrivKey) Option {
	return func(n *Node) libp2p.Option {
		return libp2p.Identity(key)
	}
}

func ListenAddr(addrs ...string) Option {
	return func(n *Node) libp2p.Option {
		return libp2p.ListenAddrStrings(addrs...)
	}
}

func Gateways(addrs ...string) Option {
	return func(n *Node) libp2p.Option {
		n.Gateways = addrs
		return noopConfig
	}
}

func DefaultNATManager() Option {
	return func(n *Node) libp2p.Option {
		return libp2p.NATPortMap()
	}
}

func EnableAutoRelay() Option {
	return func(n *Node) libp2p.Option {
		return libp2p.EnableAutoRelay()
	}
}

func ActAsRelay() Option {
	return func(n *Node) libp2p.Option {
		return libp2p.EnableRelay(relay.OptActive, relay.OptHop)
	}
}

