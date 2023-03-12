package types

type StateMachineState string

const (
	StateMachineState_Stopped StateMachineState = "Stopped"

	StateMachineState_P2P_Bootstrapping StateMachineState = "P2P_Bootstrapping"
	StateMachineState_P2P_Bootstrapped  StateMachineState = "P2P_Bootstrapped"

	StateMachineState_Consensus_Unsynched StateMachineState = "Consensus_Unsynched"
	StateMachineState_Consensus_SyncMode  StateMachineState = "Consensus_SyncMode"
	StateMachineState_Consensus_Synched   StateMachineState = "Consensus_Synched"

	// Used by synched validators to participate in the HotPOKT lifecycle
	StateMachineState_Consensus_Pacemaker StateMachineState = "Consensus_Pacemaker"
)
