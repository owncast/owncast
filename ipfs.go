package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	config "github.com/ipfs/go-ipfs-config"
	files "github.com/ipfs/go-ipfs-files"
	icore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"github.com/ipfs/interface-go-ipfs-core/path"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/core/corehttp"
	"github.com/ipfs/go-ipfs/core/node/libp2p" // This package is needed so that all the preloaded plugins are loaded automatically
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
)

var directory = "hls"
var directoryHash string

var node *core.IpfsNode

//var ctx, _ = context.WithCancel(context.Background())
var ctx = context.Background()

func createIPFSDirectory(ipfs *icore.CoreAPI, directoryName string) {
	directory, err := getUnixfsNode(directoryName)
	verifyError(err)
	defer directory.Close()

	newlyCreatedDirectoryHash, err := (*ipfs).Unixfs().Add(ctx, directory)
	verifyError(err)

	directoryHash = newlyCreatedDirectoryHash.String()
	fmt.Println("Created directory hash " + directoryHash)
}

func save(filePath string, ipfs *icore.CoreAPI) string {
	someFile, err := getUnixfsNode(filePath)

	defer someFile.Close()

	if err != nil {
		panic(fmt.Errorf("Could not get File: %s", err))
	}

	opts := []options.UnixfsAddOption{
		options.Unixfs.Pin(false),
		// options.Unixfs.CidVersion(1),
		// options.Unixfs.RawLeaves(false),
		// options.Unixfs.Nocopy(false),
	}

	cidFile, err := (*ipfs).Unixfs().Add(ctx, someFile, opts...)

	if err != nil {
		panic(fmt.Errorf("Could not add File: %s", err))
	}

	// fmt.Printf("Added file to IPFS with CID %s\n", cidFile.String())

	newHash := addFileToDirectory(ipfs, cidFile, directoryHash, filepath.Base(filePath))

	return newHash
}

func addFileToDirectory(ipfs *icore.CoreAPI, originalFileHashToModifyPath path.Path, directoryToAddTo string, filename string) string {
	directoryToAddToPath := path.New(directoryToAddTo)
	newDirectoryHash, err := (*ipfs).Object().AddLink(ctx, directoryToAddToPath, filename, originalFileHashToModifyPath)

	verifyError(err)

	fmt.Printf("New hash: %s\n", newDirectoryHash.String())

	return newDirectoryHash.String() + "/" + filename
}

func setupPlugins(externalPluginsPath string) error {
	// Load any external plugins if available on externalPluginsPath
	plugins, err := loader.NewPluginLoader(filepath.Join(externalPluginsPath, "plugins"))
	if err != nil {
		return fmt.Errorf("error loading plugins: %s", err)
	}

	// Load preloaded and external plugins
	if err := plugins.Initialize(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	if err := plugins.Inject(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	return nil
}

// Creates an IPFS node and returns its coreAPI
func createNode(ctx context.Context, repoPath string) (icore.CoreAPI, *core.IpfsNode, error) {
	fmt.Println("CreateNode...")

	// Open the repo
	repo, err := fsrepo.Open(repoPath)
	verifyError(err)

	if err != nil {
		return nil, nil, err
	}

	// Construct the node

	nodeOptions := &core.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTOption, // This option sets the node to be a full DHT node (both fetching and storing DHT Records)
		// Routing: libp2p.DHTClientOption, // This option sets the node to be a client DHT node (only fetching records)
		Repo: repo,
	}

	node, err := core.NewNode(ctx, nodeOptions)
	node.IsDaemon = true

	if err != nil {
		return nil, nil, err
	}

	// Attach the Core API to the constructed node
	coreAPI, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, nil, err
	}
	return coreAPI, node, nil
}

func createTempRepo(ctx context.Context) (string, error) {
	fmt.Println("createTempRepo...")

	repoPath, err := ioutil.TempDir("", "ipfs-shell")
	fmt.Println(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to get temp dir: %s", err)
	}

	// Create a config with default options and a 2048 bit key
	cfg, err := config.Init(ioutil.Discard, 2048)

	if err != nil {
		return "", err
	}

	// Create the repo with the config
	err = fsrepo.Init(repoPath, cfg)
	if err != nil {
		return "", fmt.Errorf("failed to init ephemeral node: %s", err)
	}

	return repoPath, nil
}

// Spawns a node to be used just for this run (i.e. creates a tmp repo)
func spawnEphemeral(ctx context.Context) (icore.CoreAPI, *core.IpfsNode, error) {
	if err := setupPlugins(""); err != nil {
		return nil, nil, err
	}

	// Create a Temporary Repo
	repoPath, err := createTempRepo(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create temp repo: %s", err)
	}

	// Spawning an ephemeral IPFS node
	coreAPI, node, err := createNode(ctx, repoPath)
	return coreAPI, node, err
}

// Spawns a node on the default repo location, if the repo exists
func spawnDefault(ctx context.Context) (icore.CoreAPI, *core.IpfsNode, error) {
	defaultPath, err := config.PathRoot()
	fmt.Println(defaultPath)
	if err != nil {
		// shouldn't be possible
		return nil, nil, err
	}

	if err := setupPlugins(defaultPath); err != nil {
		return nil, nil, err
	}

	coreAPI, node, err := createNode(ctx, defaultPath)
	return coreAPI, node, err
}

func connectToPeers(ctx context.Context, ipfs icore.CoreAPI, peers []string) error {
	var wg sync.WaitGroup
	peerInfos := make(map[peer.ID]*peerstore.PeerInfo, len(peers))
	for _, addrStr := range peers {
		addr, err := ma.NewMultiaddr(addrStr)
		if err != nil {
			return err
		}
		pii, err := peerstore.InfoFromP2pAddr(addr)
		if err != nil {
			return err
		}
		pi, ok := peerInfos[pii.ID]
		if !ok {
			pi = &peerstore.PeerInfo{ID: pii.ID}
			peerInfos[pi.ID] = pi
		}
		pi.Addrs = append(pi.Addrs, pii.Addrs...)
	}

	wg.Add(len(peerInfos))
	for _, peerInfo := range peerInfos {
		go func(peerInfo *peerstore.PeerInfo) {
			defer wg.Done()
			err := ipfs.Swarm().Connect(ctx, *peerInfo)
			if err != nil {
				log.Printf("failed to connect to %s: %s", peerInfo.ID, err)
			}
		}(peerInfo)
	}
	wg.Wait()
	return nil
}

func getUnixfsNode(path string) (files.Node, error) {
	st, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	f, err := files.NewSerialFile(path, false, st)

	if err != nil {
		return nil, err
	}

	return f, nil
}

func createIPFSInstance() (*icore.CoreAPI, *core.IpfsNode, error) {
	// Spawn a node using a temporary path, creating a temporary repo for the run
	api, node, error := spawnEphemeral(ctx)
	// api, node, error := spawnDefault(ctx)
	return &api, node, error
}

func startIPFSNode(ipfs icore.CoreAPI, node *core.IpfsNode) { //} icore.CoreAPI {
	defer fmt.Println("---- IPFS node exited!")

	fmt.Println("IPFS node is running")

	// bootstrapNodes := []string{
	// 	// IPFS Bootstrapper nodes.
	// 	"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
	// 	"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
	// 	"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
	// 	"/dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",

	// 	// IPFS Cluster Pinning nodes
	// 	"/ip4/138.201.67.219/tcp/4001/p2p/QmUd6zHcbkbcs7SMxwLs48qZVX3vpcM8errYS7xEczwRMA",

	// 	// You can add more nodes here, for example, another IPFS node you might have running locally, mine was:
	// 	// "/ip4/127.0.0.1/tcp/4010/p2p/QmZp2fhDLxjYue2RiUvLwT9MWdnbDxam32qYFnGmxZDh5L",
	// 	// "/ip4/127.0.0.1/udp/4010/quic/p2p/QmZp2fhDLxjYue2RiUvLwT9MWdnbDxam32qYFnGmxZDh5L",
	// }

	// connectToPeers(ctx, ipfs, bootstrapNodes)

	addr := "/ip4/127.0.0.1/tcp/5001"
	var opts = []corehttp.ServeOption{
		corehttp.GatewayOption(true, "/ipfs", "/ipns"),
	}

	if err := corehttp.ListenAndServe(node, addr, opts...); err != nil {
		return
	}

}
