package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

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
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
)

type IPFSStorage struct {
	ipfs *icore.CoreAPI
	node *core.IpfsNode

	ctx           context.Context
	directoryHash string
	gateway       string
}

func (s *IPFSStorage) Setup(config Config) {
	log.Println("Setting up IPFS for external storage of video. Please wait..")

	s.gateway = config.IPFS.Gateway

	s.ctx = context.Background()

	ipfsInstance, node, _ := s.createIPFSInstance()
	s.ipfs = ipfsInstance
	s.node = node

	s.createIPFSDirectory("./hls")
}

func (s *IPFSStorage) Save(filePath string, retryCount int) string {
	someFile, err := getUnixfsNode(filePath)
	defer someFile.Close()

	if err != nil {
		log.Panicln(fmt.Errorf("Could not get File: %s", err))
	}

	opts := []options.UnixfsAddOption{
		options.Unixfs.Pin(false),
		// options.Unixfs.CidVersion(1),
		// options.Unixfs.RawLeaves(false),
		// options.Unixfs.Nocopy(false),
	}

	cidFile, err := (*s.ipfs).Unixfs().Add(s.ctx, someFile, opts...)

	if err != nil {
		log.Panicln(fmt.Errorf("Could not add File: %s", err))
	}

	// fmt.Printf("Added file to IPFS with CID %s\n", cidFile.String())

	newHash := s.addFileToDirectory(cidFile, filepath.Base(filePath))

	return s.gateway + newHash
}

func (s *IPFSStorage) GenerateRemotePlaylist(playlist string, variant Variant) string {
	var newPlaylist = ""

	scanner := bufio.NewScanner(strings.NewReader(playlist))
	for scanner.Scan() {
		line := scanner.Text()
		if line[0:1] != "#" {
			fullRemotePath := variant.getSegmentForFilename(line)
			if fullRemotePath != nil {
				line = fullRemotePath.RemoteID
			} else {
				line = ""
			}
		}

		newPlaylist = newPlaylist + line + "\n"
	}

	return newPlaylist
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
	repoPath, err := ioutil.TempDir("", "ipfs-shell")
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

func (s *IPFSStorage) addFileToDirectory(originalFileHashToModifyPath path.Path, filename string) string {
	// fmt.Println("directoryToAddTo: "+s.directoryHash, "filename: "+filename, "originalFileHashToModifyPath: "+originalFileHashToModifyPath.String())
	directoryToAddToPath := path.New(s.directoryHash)
	newDirectoryHash, err := (*s.ipfs).Object().AddLink(s.ctx, directoryToAddToPath, filename, originalFileHashToModifyPath)

	verifyError(err)
	return newDirectoryHash.String() + "/" + filename
}

func (s *IPFSStorage) createIPFSInstance() (*icore.CoreAPI, *core.IpfsNode, error) {
	// Spawn a node using a temporary path, creating a temporary repo for the run
	api, node, error := spawnEphemeral(s.ctx)
	// api, node, error := spawnDefault(ctx)
	return &api, node, error
}

func (s *IPFSStorage) startIPFSNode() { //} icore.CoreAPI {
	defer log.Println("IPFS node exited")

	log.Println("IPFS node is running")

	bootstrapNodes := []string{
		// IPFS Bootstrapper nodes.
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",

		// IPFS Cluster Pinning nodes
		"/ip4/138.201.67.219/tcp/4001/p2p/QmUd6zHcbkbcs7SMxwLs48qZVX3vpcM8errYS7xEczwRMA",

		"/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",      // mars.i.ipfs.io
		"/ip4/104.131.131.82/udp/4001/quic/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ", // mars.i.ipfs.io

		// You can add more nodes here, for example, another IPFS node you might have running locally, mine was:
		// "/ip4/127.0.0.1/tcp/4010/p2p/QmZp2fhDLxjYue2RiUvLwT9MWdnbDxam32qYFnGmxZDh5L",
		// "/ip4/127.0.0.1/udp/4010/quic/p2p/QmZp2fhDLxjYue2RiUvLwT9MWdnbDxam32qYFnGmxZDh5L",
	}

	go connectToPeers(s.ctx, *s.ipfs, bootstrapNodes)

	addr := "/ip4/127.0.0.1/tcp/5001"
	var opts = []corehttp.ServeOption{
		corehttp.GatewayOption(true, "/ipfs", "/ipns"),
	}

	if err := corehttp.ListenAndServe(s.node, addr, opts...); err != nil {
		return
	}
}

func (s *IPFSStorage) createIPFSDirectory(directoryName string) {
	directory, err := getUnixfsNode(directoryName)
	verifyError(err)
	defer directory.Close()

	newlyCreatedDirectoryHash, err := (*s.ipfs).Unixfs().Add(s.ctx, directory)
	verifyError(err)
	s.directoryHash = newlyCreatedDirectoryHash.String()
}
