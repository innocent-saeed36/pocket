package identity

import (
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"

	"github.com/pokt-network/pocket/p2p/common"
	"github.com/pokt-network/pocket/p2p/transport"
	"github.com/pokt-network/pocket/p2p/types"
	poktCrypto "github.com/pokt-network/pocket/shared/crypto"
)

var (
	ErrIdentity = common.NewErrFactory("")
)

// PoktPeerFromStream returns a pokt peer, **without `ServiceURL**, from a libp2p stream.
func PoktPeerFromStream(stream network.Stream) (*types.NetworkPeer, error) {
	remotePubKeyBytes, err := stream.Conn().RemotePublicKey().Raw()
	// NB: abort handling this stream.
	if err != nil {
		return nil, err
	}
	poktPubKey, err := poktCrypto.NewPublicKeyFromBytes(remotePubKeyBytes)
	if err != nil {
		return nil, err
	}

	return &types.NetworkPeer{
		Dialer:    transport.NewLibP2PTransport(stream),
		PublicKey: poktPubKey,
		// NB: pokt analogue of libp2p peer.ID
		Address: poktPubKey.Address(),
	}, nil
}

func PubKeyFromPoktPeer(poktPeer *types.NetworkPeer) (crypto.PubKey, error) {
	pubKey, err := crypto.UnmarshalEd25519PublicKey(poktPeer.PublicKey.Bytes())
	if err != nil {
		return nil, ErrIdentity("unable to unmarshal peer ed25519 public key", err)
	}

	return pubKey, nil
}

func PeerAddrInfoFromPoktPeer(poktPeer *types.NetworkPeer) (peer.AddrInfo, error) {
	pubKey, err := PubKeyFromPoktPeer(poktPeer)
	if err != nil {
		return peer.AddrInfo{}, err
	}

	peerID, err := peer.IDFromPublicKey(pubKey)
	if err != nil {
		return peer.AddrInfo{}, ErrIdentity("unable to retrieve ID from peer public key", err)
	}

	// NB: test if service URL hostname is an IP address or an FQDN
	peerUrl, err := url.Parse("ip://" + poktPeer.ServiceUrl)
	if err != nil {
		return peer.AddrInfo{}, ErrIdentity(fmt.Sprintf(
			"unable to parse peer service URL: %s", poktPeer.ServiceUrl,
		), err)
	}

	// TODO: parameterize transport.
	peerTransportStr := "tcp"
	peerIPVersionStr := "ip4"
	// TODO: there's probably a more conventional way to do this.
	// NB: check if we're dealing with IPv4 or IPv6
	if strings.Count(peerUrl.Hostname(), ":") > 0 {
		peerIPVersionStr = "ip6"
	}

	// TODO: consider using a `/dns<4 or 6>/<hostname>`
	// multiaddr instead of resolving with stdlib here.
	// > The address parameter can use a host name, but this is not recommended,
	// > because it will return at most one of the host name's IP addresses.
	// (see: https://pkg.go.dev/net#ResolveIPAddr)
	peerIP, err := net.ResolveIPAddr(peerIPVersionStr, peerUrl.Hostname())
	if err != nil {
		return peer.AddrInfo{}, ErrIdentity(fmt.Sprintf(
			"unable to resolve peer IP for hostname: %s", peerUrl.Hostname(),
		), err)
	}

	peerMultiAddrStr := fmt.Sprintf("/%s/%s/%s/%s", peerIPVersionStr, peerIP, peerTransportStr, peerUrl.Port())
	peerMultiaddr, err := multiaddr.NewMultiaddr(peerMultiAddrStr)
	if err != nil {
		return peer.AddrInfo{}, err
	}

	return peer.AddrInfo{
		ID: peerID,
		Addrs: []multiaddr.Multiaddr{
			peerMultiaddr,
		},
	}, nil
}
