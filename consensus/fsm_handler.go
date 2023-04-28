package consensus

import (
	"fmt"

	typesCons "github.com/pokt-network/pocket/consensus/types"
	"github.com/pokt-network/pocket/shared/codec"
	coreTypes "github.com/pokt-network/pocket/shared/core/types"
	"github.com/pokt-network/pocket/shared/messaging"
	"google.golang.org/protobuf/types/known/anypb"
)

// HandleEvent handles FSM state transition events.
func (m *consensusModule) HandleEvent(transitionMessageAny *anypb.Any) error {
	m.m.Lock()
	defer m.m.Unlock()

	switch transitionMessageAny.MessageName() {
	case messaging.StateMachineTransitionEventType:
		msg, err := codec.GetCodec().FromAny(transitionMessageAny)
		if err != nil {
			return err
		}

		stateTransitionMessage, ok := msg.(*messaging.StateMachineTransitionEvent)
		if !ok {
			return fmt.Errorf("failed to cast message to StateSyncMessage")
		}
		return m.handleStateTransitionEvent(stateTransitionMessage)
	default:
		return typesCons.ErrUnknownStateSyncMessageType(transitionMessageAny.MessageName())
	}
}

func (m *consensusModule) handleStateTransitionEvent(msg *messaging.StateMachineTransitionEvent) error {
	m.logger.Debug().Fields(messaging.TransitionEventToMap(msg)).Msg("Received state machine transition msg")
	fsm_state := msg.NewState

	switch coreTypes.StateMachineState(fsm_state) {
	case coreTypes.StateMachineState_P2P_Bootstrapped:
		return m.HandleBootstrapped(msg)

	case coreTypes.StateMachineState_Consensus_Unsynced:
		return m.HandleUnsynced(msg)

	case coreTypes.StateMachineState_Consensus_SyncMode:
		return m.HandleSyncMode(msg)

	case coreTypes.StateMachineState_Consensus_Synced:
		return m.HandleSynced(msg)

	case coreTypes.StateMachineState_Consensus_Pacemaker:
		return m.HandlePacemaker(msg)

	default:
		m.logger.Warn().Msgf("Consensus module not handling this event: %s", msg.Event)

	}

	return nil
}

// HandleBootstrapped handles FSM event P2P_IsBootstrapped, and P2P_Bootstrapped is the destination state.
// Bootrstapped mode is when the node (validator or non-validator) is first coming online.
// This is a transition mode from node bootstrapping to a node being out-of-sync.
func (m *consensusModule) HandleBootstrapped(msg *messaging.StateMachineTransitionEvent) error {
	m.logger.Debug().Msg("Node is in bootstrapped state")
	return nil
}

// HandleUnsynced handles FSM event Consensus_IsUnsynced, and Unsynced is the destination state.
// In Unsynced mode node (validator or non-validator) is out of sync with the rest of the network.
// This mode is a transition mode from the node being up-to-date (i.e. Pacemaker mode, Synced mode) with the latest network height to being out-of-sync.
// As soon as node transitions to this mode, it will transition to the sync mode.
func (m *consensusModule) HandleUnsynced(msg *messaging.StateMachineTransitionEvent) error {
	m.logger.Debug().Msg("Node is in Unsyched state, as node is out of sync sending syncmode event to start syncing")

	return m.GetBus().GetStateMachineModule().SendEvent(coreTypes.StateMachineEvent_Consensus_IsSyncing)
}

// HandleSyncMode handles FSM event Consensus_IsSyncing, and SyncMode is the destination state.
// In Sync mode node (validator or non-validator) starts syncing with the rest of the network.
func (m *consensusModule) HandleSyncMode(msg *messaging.StateMachineTransitionEvent) error {
	m.logger.Debug().Msg("Node is in Sync Mode, starting to sync...")

	isValidator, err := m.GetBus().GetPersistenceModule().IsValidator(int64(m.height), m.GetNodeAddress())
	if err != nil {
		m.logger.Err(err).Msg("failed to check if the node is a validator")
		return err
	}

	// if the node is validator, start active state sync with CatchToHeight
	// else, node will always stay in sync mode, and it will catch up to the final state with passive state sync
	if isValidator {
		m.logger.Debug().Msg("Validator node is starting active state sync")
		//aggregatedMetadata := m.getAggregatedStateSyncMetadata()
		//m.stateSync.SetAggregatedMetadata(&aggregatedMetadata)
		go m.stateSync.CatchToHeight()
	}

	return nil
}

// HandleSynced handles FSM event IsSyncedNonValidator for Non-Validators, and Synced is the destination state.
// Currently, FSM never transition to this state and a non-validator node always stays in syncmode.
// CONSIDER: when a non-validator sync is implemented, maybe there is a case that requires transitioning to this state.
func (m *consensusModule) HandleSynced(msg *messaging.StateMachineTransitionEvent) error {
	m.logger.Debug().Msg("Non-validator node is in Synced mode")
	return nil
}

// HandlePacemaker handles FSM event IsSyncedValidator, and Pacemaker is the destination state.
// Execution of this state means the validator node is synced.
func (m *consensusModule) HandlePacemaker(msg *messaging.StateMachineTransitionEvent) error {
	m.logger.Debug().Msg("Validator node is Synced and in Pacemaker mode. It will stay in this mode until it receives a new block proposal that has a higher height than the current block height")
	return nil
}
