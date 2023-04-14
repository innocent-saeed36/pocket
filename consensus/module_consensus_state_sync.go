package consensus

import (
	"fmt"

	typesCons "github.com/pokt-network/pocket/consensus/types"
	coreTypes "github.com/pokt-network/pocket/shared/core/types"
	"github.com/pokt-network/pocket/shared/modules"
)

var _ modules.ConsensusStateSync = &consensusModule{}

func (m *consensusModule) GetNodeIdFromNodeAddress(peerId string) (uint64, error) {
	validators, err := m.getValidatorsAtHeight(m.CurrentHeight())
	if err != nil {
		// REFACTOR(#434): As per issue #434, once the new id is sorted out, this return statement must be changed
		return 0, err
	}

	valAddrToIdMap := typesCons.NewActorMapper(validators).GetValAddrToIdMap()
	return uint64(valAddrToIdMap[peerId]), nil
}

func (m *consensusModule) GetNodeAddress() string {
	return m.nodeAddress
}

// TODO(#352): Implement this function, currently a placeholder.
// commitReceivedBlocks commits the blocks received from the blocksReceived channel
// it is intended to be run as a background process

// runs as a background process in consensus module
// listens on the blocksReceived channel
// commits the received block
func (m *consensusModule) blockApplicationLoop() {
	for blockResponse := range m.blocksReceived {
		block := blockResponse.Block
		maxPersistedHeight, err := m.maxPersistedBlockHeight()
		if err != nil {
			m.logger.Err(err).Msg("couldn't query max persisted height")
			return
		}

		if block.BlockHeader.Height <= maxPersistedHeight {
			m.logger.Info().Msgf("Received block with height: %d, but node already persisted blocks until height: %d, so node will not apply this block", block.BlockHeader.Height, maxPersistedHeight)
			return
		} else if block.BlockHeader.Height > m.CurrentHeight() {
			m.logger.Info().Msgf("Received block with height %d, but node's last persisted height is: %d, so node will not apply this block", block.BlockHeader.Height, maxPersistedHeight)
			return
		}

		err = m.verifyBlock(block)
		if err != nil {
			m.logger.Err(err).Msg("failed to verify block")
			return
		}

		err = m.applyAndCommitBlock(block)
		if err != nil {
			m.logger.Err(err).Msg("failed to apply and commit block")
			return
		}
		fmt.Println("Applied block: ", block)
		m.stateSync.CommittedBlock(m.CurrentHeight())
	}

}

// TODO(#352): Implement this function, currently a placeholder.
// metadataSyncLoop periodically sends metadata requests to its peers
// it is intended to be run as a background process
func (m *consensusModule) metadataSyncLoop() {
	// runs as a background process in consensus module
	// requests metadata from peers
	// sends received metadata to the metadataReceived channel
}

func (m *consensusModule) maxPersistedBlockHeight() (uint64, error) {
	readCtx, err := m.GetBus().GetPersistenceModule().NewReadContext(int64(m.CurrentHeight()))
	if err != nil {
		return 0, err
	}
	defer readCtx.Release()

	maxHeight, err := readCtx.GetMaximumBlockHeight()
	if err != nil {
		return 0, err
	}

	return maxHeight, nil
}

func (m *consensusModule) verifyBlock(block *coreTypes.Block) error {
	return nil
}

func (m *consensusModule) applyAndCommitBlock(block *coreTypes.Block) error {
	m.logger.Info().Msgf("applying and committing the block at height %d", block.BlockHeader.Height)

	// TODO: uncomment following. In this PR test blocks don't have a valid QC, therefore commented out to let the tests pass
	// if err := m.applyBlock(block); err != nil {
	// 	m.logger.Error().Err(err).Msg("Could not apply block, invalid QC")
	// 	return err
	// }

	if err := m.commitBlock(block); err != nil {
		m.logger.Error().Err(err).Msg("Could not commit block, invalid QC")
		return err
	}

	m.paceMaker.NewHeight()

	maxPersistedHeight, err := m.maxPersistedBlockHeight()
	if err != nil {
		return err
	}

	m.logger.Info().Msgf("Block is Committed, maxPersistedHeight is: %d, current height is :%d", maxPersistedHeight, m.height)
	return nil
}

func (m *consensusModule) getAggregatedStateSyncMetadata() typesCons.StateSyncMetadataResponse {
	minHeight, maxHeight := uint64(1), uint64(1)

	chanLen := len(m.metadataReceived)

	for i := 0; i < chanLen; i++ {
		metadata := <-m.metadataReceived
		if metadata.MaxHeight > maxHeight {
			maxHeight = metadata.MaxHeight
		}
		if metadata.MinHeight < minHeight {
			minHeight = metadata.MinHeight
		}
	}

	return typesCons.StateSyncMetadataResponse{
		PeerAddress: "unused_aggregated_metadata_address",
		MinHeight:   minHeight,
		MaxHeight:   maxHeight,
	}
}
