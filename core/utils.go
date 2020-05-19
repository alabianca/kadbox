package core

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/mitchellh/go-homedir"
	"os"
	"path"
	"strings"
)

const (
	RootDirectoryName  = ".kadbox"
	ConfigName         = "kadconfig"
	StoreDirectoryName = "store"
	ConfigNotFoundErr  = "config not found"
	RootNotFoundErr    = "root not found"
)

// rootDirectory creates the root directory .kadbox.
// relative to the given root and returns the full path
// if root is an empty string we use the users home directory
func rootDirectory(root string) string {
	var home string
	var err error
	if root == "" {
		home, err = homedir.Dir()
		if err != nil {
			panic(err)
		}

	} else {
		home = root
	}

	// check if the directory exists
	p, err := mkdirIfNotExist(path.Join(home, RootDirectoryName))
	if err != nil {
		panic(err)
	}

	return p
}

func storeDirectory(root string) string {
	if p, err := mkdirIfNotExist(path.Join(root, RootDirectoryName, StoreDirectoryName)); err != nil {
		panic(err)
	} else {
		return p
	}
}

// rootDirExists checks if the root directory .kadbox exists in dir
func rootDirExists(dir string) (string, bool) {
	p := path.Join(dir, RootDirectoryName)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return "", false
	}

	return p, true
}

func mkdirIfNotExist(name string) (string, error) {
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		return name, nil
	}

	if err := os.Mkdir(name, 0755); err != nil {
		return "", err
	}

	return name, nil
}

func configFileExists(dir string) bool {
	if _, err := os.Stat(path.Join(dir, ConfigName)); os.IsNotExist(err) {
		return false
	}

	return true
}

func createNewConfigFromDefault() Repo {
	fmt.Println("Generating Public/Private Key Pair ...")
	priv, pub, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	if err != nil {
		panic(err)
	}

	id, err := peer.IDFromPublicKey(pub)
	if err != nil {
		panic(err)
	}

	marshalled, err := priv.Bytes()
	if err != nil {
		panic(err)
	}

	repo := Repo{
		Gateways:    []string{"/ip4/127.0.0.1/tcp/4000/p2p/QmQnAZsyiJSovuqg8zjP3nKdm6Pwb75Mpn8HnGyD5WYZ15"},
		ListenAddrs: []string{"/ip4/0.0.0.0/tcp/0"},
		Identity: Identity{
			PrivateKey: hex.EncodeToString(marshalled),
			ID:         id,
		},
		API: API{
			Address: "127.0.0.1",
			Port:    3000,
		},
	}

	return repo
}

// initializeConfig create the config file at root
func initializeConfig(root string) Repo {
	home := rootDirectory(root)

	p := path.Join(home, ConfigName)
	if !configFileExists(home) {
		fmt.Println("Config file does not yet exist so we create it at ", p)
		defConfig := createNewConfigFromDefault()
		file, err := os.Create(p)
		if err != nil {
			panic(err)
		}

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "	")
		if err := encoder.Encode(&defConfig); err != nil {
			panic(err)
		}

		return defConfig
	}

	// we have a config already so load it
	fmt.Println("Opening up ", p)
	file, err := os.Open(p)
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(file)
	var repo Repo

	if err := decoder.Decode(&repo); err != nil {
		panic(err)
	}

	return repo

}

// getConfigFile starts at dir and moves up to home
// until it finds a config file. If not config file is found ConfigNotFoundErr is returned
func getConfigFile(dir string, home string) (Repo, error) {
	if configFileExists(path.Join(dir, RootDirectoryName)) {
		return initializeConfig(dir), nil
	}

	if dir == home {
		return Repo{}, errors.New(ConfigNotFoundErr)
	}

	parts := strings.Split(dir, string(os.PathSeparator))
	dir = strings.Join(
		parts[:len(parts)-1],
		string(os.PathSeparator),
	)

	return getConfigFile(dir, home)

}

// getKadboxRepo starts at dir and moves up to home
/// unitl it finds a .kadbox repo directory. If found the path is returned
// otherwise repo not found is returned
func getKadboxRepoDirectory(dir string, home string) (string, error) {
	if p, ok := rootDirExists(dir); ok {
		return p, nil
	}

	if dir == home {
		return "", errors.New(RootNotFoundErr)
	}

	parts := strings.Split(dir, string(os.PathSeparator))
	dir = strings.Join(
		parts[:len(parts)-1],
		string(os.PathSeparator),
	)

	return getKadboxRepoDirectory(dir, home)
}

func GetClosestKadboxRepoRelativeToWd() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return getKadboxRepoDirectory(wd, home)
}

func getClosestConfigFileRelativeToWd() (Repo, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Repo{}, err
	}

	wd, err := os.Getwd()
	if err != nil {
		return Repo{}, err
	}

	return getConfigFile(wd, home)
}
