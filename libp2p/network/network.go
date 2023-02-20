package network

import (
	"context"
	"fmt"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	libp2pHost "github.com/libp2p/go-libp2p/core/host"

	"github.com/pokt-network/pocket/libp2p/identity"
	"github.com/pokt-network/pocket/libp2p/protocol"
	"github.com/pokt-network/pocket/p2p/providers"
	"github.com/pokt-network/pocket/p2p/types"
	"github.com/pokt-network/pocket/shared/crypto"
	"github.com/pokt-network/pocket/shared/modules"
)

type libp2pNetwork struct {
	logger      *modules.Logger
	bus         modules.Bus
	host        libp2pHost.Host
	topic       *pubsub.Topic
	addrBookMap types.AddrBookMap
}

var (
	ErrNetwork = types.NewErrFactory("libp2p network error")
	Year       = time.Hour * 24 * 365
	// TECHDEBT: consider more carefully and parameterize.
	DefaultPeerTTL = 2 * Year
)

// TECHDEBT: factor out args which are common to network
// implementations to an options or config struct.
func NewLibp2pNetwork(
	bus modules.Bus,
	addrBookProvider providers.AddrBookProvider,
	currentHeightProvider providers.CurrentHeightProvider,
	logger *modules.Logger,
	host libp2pHost.Host,
	topic *pubsub.Topic,
) (types.Network, error) {
	addrBook, err := addrBookProvider.GetStakedAddrBookAtHeight(currentHeightProvider.CurrentHeight())
	if err != nil {
		logger.Fatal().Err(err).Msg("getting staked address book")
	}

	addrBookMap := make(types.AddrBookMap)
	for _, peer := range addrBook {
		addrBookMap[peer.Address.String()] = peer
		pubKey, err := identity.Libp2pPublicKeyFromPeer(peer)
		if err != nil {
			return nil, ErrNetwork(fmt.Sprintf(
				"converting peer public key, pokt address: %s", peer.Address,
			), err)
		}
		libp2pPeer, err := identity.Libp2pAddrInfoFromPeer(peer)
		if err != nil {
			return nil, ErrNetwork(fmt.Sprintf(
				"converting peer info, pokt address: %s", peer.Address,
			), err)
		}

		host.Peerstore().AddAddrs(libp2pPeer.ID, libp2pPeer.Addrs, DefaultPeerTTL)
		if err := host.Peerstore().AddPubKey(libp2pPeer.ID, pubKey); err != nil {
			return nil, ErrNetwork(fmt.Sprintf(
				"adding peer public key, pokt address: %s", peer.Address,
			), err)
		}
	}

	return &libp2pNetwork{
		logger: logger,
		// TODO: is it unconventional to set bus here?
		bus:         bus,
		host:        host,
		topic:       topic,
		addrBookMap: addrBookMap,
	}, nil
}

// NetworkBroadcast uses the configured pubsub router to broadcast data to peers.
func (p2pNet *libp2pNetwork) NetworkBroadcast(data []byte) error {
	// IMPROVE: receive context in interface methods?
	ctx := context.Background()

	// NB: Routed send using pubsub
	if err := p2pNet.topic.Publish(ctx, data); err != nil {
		return ErrNetwork("unable to publish to topic", err)
	}
	return nil
}

// NetworkSend connects sends data directly to the specified peer.
func (p2pNet *libp2pNetwork) NetworkSend(data []byte, poktAddr crypto.Address) error {
	// IMPROVE: add context to interface methods.
	ctx := context.Background()

	selfPoktAddr, err := p2pNet.GetBus().GetP2PModule().GetAddress()
	if err != nil {
		return ErrNetwork(fmt.Sprintf(
			"sending to poktAddr: %s", poktAddr,
		), err)
	}

	// NB: don't send to self.
	if selfPoktAddr.Equals(poktAddr) {
		return nil
	}

	peer, ok := p2pNet.addrBookMap[poktAddr.String()]
	if !ok {
		// NB: this should not happen.
		return ErrNetwork("", fmt.Errorf(
			"peer not found in address book, pokt address: %s", poktAddr,
		))
	}

	peerAddrInfo, err := identity.Libp2pAddrInfoFromPeer(peer)
	if err != nil {
		return ErrNetwork("parsing peer multiaddr", err)
	}

	stream, err := p2pNet.host.NewStream(ctx, peerAddrInfo.ID, protocol.PoktProtocolID)
	if err != nil {
		return ErrNetwork(fmt.Sprintf(
			"opening peer stream, pokt address: %s", poktAddr,
		), err)
	}

	if _, err := stream.Write(data); err != nil {
		return ErrNetwork(fmt.Sprintf(
			"writing to stream (peer address: %s)", poktAddr,
		), err)
	}
	defer func() {
		// NB: close the stream so that peer receives EOF.
		if err := stream.Close(); err != nil {
			p2pNet.logger.Error().Err(err).Msg(fmt.Sprintf(
				"closing peer stream, pokt address: %s", poktAddr,
			))
		}
	}()
	return nil
}

// This function was added to specifically support the RainTree implementation.
// Handles the raw data received from the network and returns the data to be processed
// by the application layer.
func (p2pNet *libp2pNetwork) HandleNetworkData(data []byte) ([]byte, error) {
	return data, nil
}

func (p2pNet *libp2pNetwork) GetAddrBook() types.AddrBook {
	addrBook := make(types.AddrBook, 0)
	for _, peer := range p2pNet.addrBookMap {
		addrBook = append(addrBook, peer)
	}
	return addrBook
}

func (p2pNet *libp2pNetwork) AddPeerToAddrBook(peer *types.NetworkPeer) error {
	p2pNet.addrBookMap[peer.Address.String()] = peer

	pubKey, err := identity.Libp2pPublicKeyFromPeer(peer)
	if err != nil {
		return ErrNetwork(fmt.Sprintf(
			"converting peer public key, pokt address: %s", peer.Address,
		), err)
	}
	libp2pPeer, err := identity.Libp2pAddrInfoFromPeer(peer)
	if err != nil {
		return ErrNetwork(fmt.Sprintf(
			"converting peer info, pokt address: %s", peer.Address,
		), err)
	}

	p2pNet.host.Peerstore().AddAddrs(libp2pPeer.ID, libp2pPeer.Addrs, DefaultPeerTTL)
	if err := p2pNet.host.Peerstore().AddPubKey(libp2pPeer.ID, pubKey); err != nil {
		return ErrNetwork(fmt.Sprintf(
			"adding peer public key, pokt address: %s", peer.Address,
		), err)
	}
	return nil
}

func (p2pNet *libp2pNetwork) RemovePeerFromAddrBook(peer *types.NetworkPeer) error {
	delete(p2pNet.addrBookMap, peer.Address.String())

	libp2pPeer, err := identity.Libp2pAddrInfoFromPeer(peer)
	if err != nil {
		return ErrNetwork(fmt.Sprintf(
			"converting peer info, pokt address: %s", peer.Address,
		), err)
	}

	p2pNet.host.Peerstore().RemovePeer(libp2pPeer.ID)
	return nil
}

func (p2pNet *libp2pNetwork) GetBus() modules.Bus {
	return p2pNet.bus
}

func (p2pNet *libp2pNetwork) SetBus(bus modules.Bus) {
	p2pNet.bus = bus
}

func (p2pNet *libp2pNetwork) Close() error {
	return p2pNet.host.Close()
}
