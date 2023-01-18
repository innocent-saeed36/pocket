package state_sync

import (
	"fmt"

	typesCons "github.com/pokt-network/pocket/consensus/types"
	"github.com/pokt-network/pocket/shared/codec"
	"google.golang.org/protobuf/types/known/anypb"
)

// This module is responsible for handling requests and business logic that advertises and shares
// local state metadata with other peers synching to the latest block.
type StateSyncServerModule interface {
	//modules.Module

	// Advertise (send) the local state sync metadata to the requesting peer
	HandleStateSyncMetadataRequest(*typesCons.StateSyncMetadataRequest) error

	// Send the block being requested by the peer
	HandleGetBlockRequest(*typesCons.GetBlockRequest) error
}

func (m *stateSync) HandleStateSyncMetadataRequest(metadataReq *typesCons.StateSyncMetadataRequest) error {

	peerId, err := m.GetBus().GetConsensusModule().GetCurrentNodeAddressFromNodeId()
	if err != nil {
		return err
	}
	minHeight, maxHeight := m.aggregateMetaResults()

	m.nodeLog(fmt.Sprintf("%s RECEIVED STATE SYNC METADATA REQUEST FROM: %s", peerId, metadataReq.PeerId))

	metadataRes := typesCons.StateSyncMetadataResponse{
		PeerId:    peerId,
		MinHeight: minHeight,
		MaxHeight: maxHeight,
	}

	anyStateSyncMessage, err := anypb.New(&metadataRes)
	if err != nil {
		return err
	}

	m.nodeLog(fmt.Sprintf("REPLYING TO STATE SYNC METADATA REQUEST: %s", metadataRes.String()))

	return m.sendToPeer(anyStateSyncMessage, metadataReq.PeerId)
	//return nil
}

func (m *stateSync) HandleGetBlockRequest(blockReq *typesCons.GetBlockRequest) error {

	peerId, err := m.GetBus().GetConsensusModule().GetCurrentNodeAddressFromNodeId()
	if err != nil {
		return err
	}

	block, err := m.getBlockAtHeight(blockReq.Height)
	if err != nil {
		return err
	}

	metadataRes := typesCons.GetBlockResponse{
		PeerId: peerId,
		Block:  block,
	}

	anyStateSyncMessage, err := anypb.New(&metadataRes)
	if err != nil {
		return err
	}

	return m.sendToPeer(anyStateSyncMessage, blockReq.PeerId)
}

// TODO! Placeholder function for metadata aggregation of data received from different peers
func (m *stateSync) aggregateMetaResults() (uint64, uint64) {
	minHeight := m.GetBus().GetConsensusModule().CurrentHeight()
	maxHeight := m.GetBus().GetConsensusModule().CurrentHeight()
	return minHeight, maxHeight
}

func (m *stateSync) getBlockAtHeight(blockHeight uint64) (*typesCons.Block, error) {

	blockStore := m.GetBus().GetPersistenceModule().GetBlockStore()
	heightBytes := heightToBytes(int64(blockHeight))
	blockBytes, err := blockStore.Get(heightBytes)

	if err != nil {
		return &typesCons.Block{}, err
	}

	var block typesCons.Block
	err = codec.GetCodec().Unmarshal(blockBytes, &block)
	if err != nil {
		return &typesCons.Block{}, err
	}

	return &block, nil

}
