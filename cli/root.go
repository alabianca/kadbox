package cli

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alabianca/kadbox/core"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"os"
	"path"
	"strings"
)

const (
	rootDirectoryName  = ".kadbox"
	configName         = "kadconfig"
	storeDirectoryName = "store"
	configNotFoundErr  = "config not found"
	rootNotFoundErr    = "root not found"
)

var rootCmd = &cobra.Command{
	Use:   "kadbox",
	Short: "kadbox share and store files in a decentralized manner",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

//// rootDirectory creates the root directory .kadbox.
//// relative to the given root and returns the full path
//// if root is an empty string we use the users home directory
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
	p, err := mkdirIfNotExist(path.Join(home, rootDirectoryName))
	if err != nil {
		panic(err)
	}

	return p
}

func storeDirectory(root string) string {
	if p, err := mkdirIfNotExist(path.Join(root, rootDirectoryName, storeDirectoryName)); err != nil {
		panic(err)
	} else {
		return p
	}
}

//rootDirExists checks if the root directory .kadbox exists in dir
func rootDirExists(dir string) (string, bool) {
	p := path.Join(dir, rootDirectoryName)
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

func createNewConfigFromDefault() core.Repo {
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

	repo := core.Repo{
		Gateways:    []string{"/ip4/134.209.171.195/tcp/5000/p2p/QmaVXTwU1j2PLquaRkHJUrG5A9HojTmwcWrdYfinTcswN7"},
		ListenAddrs: []string{"/ip4/0.0.0.0/tcp/0"},
		Identity: core.Identity{
			PrivateKey: hex.EncodeToString(marshalled),
			ID:         id,
		},
		API: core.API{
			Address: "127.0.0.1",
			Port:    3000,
		},
	}

	return repo
}

func configFileExists(dir string) bool {
	if _, err := os.Stat(path.Join(dir, configName)); os.IsNotExist(err) {
		return false
	}

	return true
}

// initializeConfig create the config file at root
func initializeConfig(root string) core.Repo {
	home := rootDirectory(root)

	p := path.Join(home, configName)
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
	var repo core.Repo

	if err := decoder.Decode(&repo); err != nil {
		panic(err)
	}

	return repo

}

// getConfigFile starts at dir and moves up to home
// until it finds a config file. If not config file is found configNotFoundErr is returned
func getConfigFile(dir string, home string) (core.Repo, error) {
	if configFileExists(path.Join(dir, rootDirectoryName)) {
		return initializeConfig(dir), nil
	}

	if dir == home {
		return core.Repo{}, errors.New(configNotFoundErr)
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
		return "", errors.New(rootNotFoundErr)
	}

	parts := strings.Split(dir, string(os.PathSeparator))
	dir = strings.Join(
		parts[:len(parts)-1],
		string(os.PathSeparator),
	)

	return getKadboxRepoDirectory(dir, home)
}

func getClosestKadboxRepoRelativeToWd() (string, error) {
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

func getClosestConfigFileRelativeToWd() (core.Repo, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return core.Repo{}, err
	}

	wd, err := os.Getwd()
	if err != nil {
		return core.Repo{}, err
	}

	return getConfigFile(wd, home)
}

func printError(err error) int {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	return 1
}
