package core

import (
	"encoding/hex"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
)

type Identity struct {
	ID         peer.ID `json:"id"`
	PrivateKey string  `json:"privateKey"`
}

func (i Identity) GetPrivateKey() (crypto.PrivKey, error) {
	var key crypto.PrivKey
	bts, err := hex.DecodeString(i.PrivateKey)
	if err != nil {
		return key, err
	}

	return crypto.UnmarshalPrivateKey(bts)
}

type API struct {
	Port    int16  `json:"port"`
	Address string `json:"address"`
}

type Repo struct {
	Gateways    []string `json:"gateways"`
	Identity    Identity `json:"identity"`
	ListenAddrs []string `json:"listenAddresses"`
	API         API      `json:"api"`
}
