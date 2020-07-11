module github.com/gabek/owncast

go 1.14

require (
	github.com/aws/aws-sdk-go v1.32.1
	github.com/ipfs/go-ipfs v0.5.1
	github.com/ipfs/go-ipfs-config v0.5.3
	github.com/ipfs/go-ipfs-files v0.0.8
	github.com/ipfs/interface-go-ipfs-core v0.2.7
	github.com/libp2p/go-libp2p-peer v0.2.0
	github.com/libp2p/go-libp2p-peerstore v0.2.6
	github.com/mssola/user_agent v0.5.2
	github.com/multiformats/go-multiaddr v0.2.2
	github.com/nareix/joy5 v0.0.0-20200710135721-d57196c8d506
	github.com/radovskyb/watcher v1.0.7
	github.com/sirupsen/logrus v1.6.0
	github.com/teris-io/shortid v0.0.0-20171029131806-771a37caa5cf
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/nareix/joy4 v0.0.0 => github.com/Seize/joy4 v0.0.0
