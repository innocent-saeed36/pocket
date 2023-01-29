package e2e_tests

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/pokt-network/pocket/consensus"
	typesCons "github.com/pokt-network/pocket/consensus/types"
	coreTypes "github.com/pokt-network/pocket/shared/core/types"
	"github.com/pokt-network/pocket/shared/modules"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestHotstuff4Nodes1BlockHappyPath(t *testing.T) {
	// Test preparation
	clockMock := clock.NewMock()
	timeReminder(t, clockMock, time.Second)

	// Test configs
	runtimeMgrs := GenerateNodeRuntimeMgrs(t, numValidators, clockMock)
	buses := GenerateBuses(t, runtimeMgrs)

	// Create & start test pocket nodes
	eventsChannel := make(modules.EventsChannel, 100)
	pocketNodes := CreateTestConsensusPocketNodes(t, buses, eventsChannel)
	StartAllTestPocketNodes(t, pocketNodes)

	// Debug message to start consensus by triggering first view change
	for _, pocketNode := range pocketNodes {
		TriggerNextView(t, pocketNode)
	}
	advanceTime(t, clockMock, 10*time.Millisecond)

	// 1. NewRound
	newRoundMessages, err := WaitForNetworkConsensusEvents(t, clockMock, eventsChannel, consensus.NewRound, consensus.Propose, numValidators*numValidators, 250, true)
	require.NoError(t, err)
	for nodeId, pocketNode := range pocketNodes {
		nodeState := GetConsensusNodeState(pocketNode)
		assertNodeConsensusView(t, nodeId,
			typesCons.ConsensusNodeState{
				Height: 1,
				Step:   uint8(consensus.NewRound),
				Round:  0,
			},
			nodeState)
		require.Equal(t, false, nodeState.IsLeader)
		require.Equal(t, typesCons.NodeId(0), nodeState.LeaderId)
	}

	for _, message := range newRoundMessages {
		P2PBroadcast(t, pocketNodes, message)
	}
	advanceTime(t, clockMock, 10*time.Millisecond)

	// IMPROVE: Use seeding for deterministic leader election in unit tests.
	// Leader election is deterministic for now, so we know its NodeId
	leaderId := typesCons.NodeId(2)
	leader := pocketNodes[leaderId]

	// 2. Prepare
	prepareProposal, err := WaitForNetworkConsensusEvents(t, clockMock, eventsChannel, consensus.Prepare, consensus.Propose, numValidators, 250, true)
	require.NoError(t, err)
	for nodeId, pocketNode := range pocketNodes {
		nodeState := GetConsensusNodeState(pocketNode)
		assertNodeConsensusView(t, nodeId,
			typesCons.ConsensusNodeState{
				Height: 1,
				Step:   uint8(consensus.Prepare),
				Round:  0,
			},
			nodeState)
		require.Equal(t, leaderId, nodeState.LeaderId, fmt.Sprintf("%d should be the current leader", leaderId))
	}

	for _, message := range prepareProposal {
		P2PBroadcast(t, pocketNodes, message)
	}
	advanceTime(t, clockMock, 10*time.Millisecond)

	// 3. PreCommit
	prepareVotes, err := WaitForNetworkConsensusEvents(t, clockMock, eventsChannel, consensus.Prepare, consensus.Vote, numValidators, 250, true)
	require.NoError(t, err)

	for _, vote := range prepareVotes {
		P2PSend(t, leader, vote)
	}
	advanceTime(t, clockMock, 10*time.Millisecond)

	preCommitProposal, err := WaitForNetworkConsensusEvents(t, clockMock, eventsChannel, consensus.PreCommit, consensus.Propose, numValidators, 250, true)
	require.NoError(t, err)
	for nodeId, pocketNode := range pocketNodes {
		nodeState := GetConsensusNodeState(pocketNode)
		assertNodeConsensusView(t, nodeId,
			typesCons.ConsensusNodeState{
				Height: 1,
				Step:   uint8(consensus.PreCommit),
				Round:  0,
			},
			nodeState)
		require.Equal(t, leaderId, nodeState.LeaderId, fmt.Sprintf("%d should be the current leader", leaderId))
	}

	for _, message := range preCommitProposal {
		P2PBroadcast(t, pocketNodes, message)
	}
	advanceTime(t, clockMock, 10*time.Millisecond)

	// 4. Commit
	preCommitVotes, err := WaitForNetworkConsensusEvents(t, clockMock, eventsChannel, consensus.PreCommit, consensus.Vote, numValidators, 250, true)
	require.NoError(t, err)

	for _, vote := range preCommitVotes {
		P2PSend(t, leader, vote)
	}
	advanceTime(t, clockMock, 10*time.Millisecond)

	commitProposal, err := WaitForNetworkConsensusEvents(t, clockMock, eventsChannel, consensus.Commit, consensus.Propose, numValidators, 250, true)
	require.NoError(t, err)
	for nodeId, pocketNode := range pocketNodes {
		nodeState := GetConsensusNodeState(pocketNode)
		assertNodeConsensusView(t, nodeId,
			typesCons.ConsensusNodeState{
				Height: 1,
				Step:   uint8(consensus.Commit),
				Round:  0,
			},
			nodeState)
		require.Equal(t, leaderId, nodeState.LeaderId, fmt.Sprintf("%d should be the current leader", leaderId))
	}

	for _, message := range commitProposal {
		P2PBroadcast(t, pocketNodes, message)
	}
	advanceTime(t, clockMock, 10*time.Millisecond)

	// 5. Decide
	commitVotes, err := WaitForNetworkConsensusEvents(t, clockMock, eventsChannel, consensus.Commit, consensus.Vote, numValidators, 250, true)
	require.NoError(t, err)

	for _, vote := range commitVotes {
		P2PSend(t, leader, vote)
	}
	advanceTime(t, clockMock, 10*time.Millisecond)

	decideProposal, err := WaitForNetworkConsensusEvents(t, clockMock, eventsChannel, consensus.Decide, consensus.Propose, numValidators, 250, true)
	require.NoError(t, err)
	for pocketId, pocketNode := range pocketNodes {
		nodeState := GetConsensusNodeState(pocketNode)
		// Leader has already committed the block and hence moved to the next height.
		if pocketId == leaderId {
			assertNodeConsensusView(t, pocketId,
				typesCons.ConsensusNodeState{
					Height: 2,
					Step:   uint8(consensus.NewRound),
					Round:  0,
				},
				nodeState)
			require.Equal(t, nodeState.LeaderId, typesCons.NodeId(0), "Leader should be empty")
			continue
		}
		assertNodeConsensusView(t, pocketId,
			typesCons.ConsensusNodeState{
				Height: 1,
				Step:   uint8(consensus.Decide),
				Round:  0,
			},
			nodeState)
		require.Equal(t, leaderId, nodeState.LeaderId, fmt.Sprintf("%d should be the current leader", leaderId))
	}

	for _, message := range decideProposal {
		P2PBroadcast(t, pocketNodes, message)
	}
	advanceTime(t, clockMock, 10*time.Millisecond)

	// 1. NewRound - begin again
	_, err = WaitForNetworkConsensusEvents(t, clockMock, eventsChannel, consensus.NewRound, consensus.Propose, numValidators*numValidators, 250, true)
	require.NoError(t, err)
	for pocketId, pocketNode := range pocketNodes {
		nodeState := GetConsensusNodeState(pocketNode)
		assertNodeConsensusView(t, pocketId,
			typesCons.ConsensusNodeState{
				Height: 2,
				Step:   uint8(consensus.NewRound),
				Round:  0,
			},
			nodeState)
		require.Equal(t, nodeState.LeaderId, typesCons.NodeId(0), "Leader should be empty")
	}
}

func TestHotstuff4NodesByzantineLeaderProposalRejected(t *testing.T) {
	fmt.Println("\n STARTING BYZANTINE NODE TEST")
	// Test preparation
	clockMock := clock.NewMock()
	timeReminder(t, clockMock, time.Second)

	// Test configs
	runtimeMgrs := GenerateNodeRuntimeMgrs(t, numValidators, clockMock)
	buses := GenerateBuses(t, runtimeMgrs)

	// Create & start test pocket nodes
	eventsChannel := make(modules.EventsChannel, 100)
	pocketNodes := CreateTestConsensusPocketNodes(t, buses, eventsChannel)
	StartAllTestPocketNodes(t, pocketNodes)

	testHeight := uint64(3)
	testStep := uint8(consensus.NewRound)
	testRound := uint64(0)

	leaderId := typesCons.NodeId(3)
	leader := pocketNodes[leaderId]
	consensusPK, err := leader.GetBus().GetConsensusModule().GetPrivateKey()
	require.NoError(t, err)

	leaderByzantineHeight := testHeight + 10

	// Placeholder block
	blockHeader := &coreTypes.BlockHeader{
		Height:            leaderByzantineHeight,
		StateHash:         stateHash,
		PrevStateHash:     "",
		NumTxs:            0,
		ProposerAddress:   consensusPK.Address(),
		QuorumCertificate: nil,
	}
	block := &coreTypes.Block{
		BlockHeader:  blockHeader,
		Transactions: make([][]byte, 0),
	}

	leaderConsensusModImpl := GetConsensusModImpl(leader)
	leaderConsensusModImpl.MethodByName("SetBlock").Call([]reflect.Value{reflect.ValueOf(block)})

	for _, pocketNode := range pocketNodes {
		// Update height, step, leaderId, and utility context via setters exposed with the debug interface
		consensusModImpl := GetConsensusModImpl(pocketNode)
		consensusModImpl.MethodByName("SetHeight").Call([]reflect.Value{reflect.ValueOf(testHeight)})
		consensusModImpl.MethodByName("SetStep").Call([]reflect.Value{reflect.ValueOf(testStep)})
		consensusModImpl.MethodByName("SetRound").Call([]reflect.Value{reflect.ValueOf(testRound)})
		consensusModImpl.MethodByName("SetLeaderId").Call([]reflect.Value{reflect.Zero(reflect.TypeOf(&leaderId))})

		// utilityContext is only set on new rounds, which is skipped in this test
		utilityContext, err := pocketNode.GetBus().GetUtilityModule().NewContext(int64(testHeight))
		require.NoError(t, err)
		consensusModImpl.MethodByName("SetUtilityContext").Call([]reflect.Value{reflect.ValueOf(utilityContext)})
	}

	// // Debug message to start consensus by triggering view change
	// for _, pocketNode := range pocketNodes {
	// 	TriggerNextView(t, pocketNode)
	// }
	// advanceTime(t, clockMock, 10*time.Millisecond)

	// 1. NewRound
	// newRoundMessages, err := WaitForNetworkConsensusEvents(t, clockMock, eventsChannel, consensus.NewRound, consensus.Propose, numValidators*numValidators, 250, true)
	// require.NoError(t, err)
	// for pocketId, pocketNode := range pocketNodes {
	// 	nodeState := GetConsensusNodeState(pocketNode)
	// 	assertNodeConsensusView(t, pocketId,
	// 		typesCons.ConsensusNodeState{
	// 			Height: 4,
	// 			Step:   uint8(consensus.NewRound),
	// 			Round:  1,
	// 		},
	// 		nodeState)
	// 	require.Equal(t, false, nodeState.IsLeader)
	// 	require.Equal(t, nodeState.LeaderId, typesCons.NodeId(0), "Leader should be empty")
	// }

	// for _, message := range newRoundMessages {
	// 	P2PBroadcast(t, pocketNodes, message)
	// }
	// advanceTime(t, clockMock, 10*time.Millisecond)

	prepareProposal := &typesCons.HotstuffMessage{
		Type:          consensus.Propose,
		Height:        leaderByzantineHeight,
		Step:          consensus.Prepare, //typesCons.HotstuffStep(testStep),
		Round:         testRound,
		Block:         block,
		Justification: nil,
	}
	anyMsg, err := anypb.New(prepareProposal)
	require.NoError(t, err)

	fmt.Println("\n LEADER IS SENDING WRONG PROPOSAL")
	P2PBroadcast(t, pocketNodes, anyMsg)

	numExpectedMsgs := 0
	_, err = WaitForNetworkConsensusEvents(t, clockMock, eventsChannel, consensus.Prepare, consensus.Vote, numExpectedMsgs, 250, true)
	require.NoError(t, err)

	// for nodeId, pocketNode := range pocketNodes {
	// 	nodeState := GetConsensusNodeState(pocketNode)
	// 	if nodeId == leaderId {
	// 		require.Equal(t, consensus.Prepare.String(), typesCons.HotstuffStep(nodeState.Step).String())
	// 	} else {
	// 		require.Equal(t, consensus.PreCommit.String(), typesCons.HotstuffStep(nodeState.Step).String())
	// 	}
	// 	require.Equal(t, testHeight, nodeState.Height)
	// 	require.Equal(t, uint8(0), nodeState.Round)
	// 	require.Equal(t, leaderId, nodeState.LeaderId)
	// 	//require.Equal(t, leaderId, nodeState.LeaderId, fmt.Sprintf("%d should be the current leader", leaderId))
	// }

	for pocketId, pocketNode := range pocketNodes {
		nodeState := GetConsensusNodeState(pocketNode)
		assertNodeConsensusView(t, pocketId,
			typesCons.ConsensusNodeState{
				Height: testHeight,
				Step:   uint8(consensus.NewRound),
				Round:  uint8(testRound),
			},
			nodeState)
		//require.Equal(t, false, nodeState.IsLeader)
		require.Equal(t, typesCons.NodeId(0), nodeState.LeaderId, "Leader should be empty")
	}

	// Debug message to start consensus by triggering next view
	for _, pocketNode := range pocketNodes {
		TriggerNextView(t, pocketNode)
	}
	advanceTime(t, clockMock, 10*time.Millisecond)

	leaderId = 3
	leader = pocketNodes[leaderId]
	for _, pocketNode := range pocketNodes {
		// Update height, step, leaderId, and utility context via setters exposed with the debug interface
		consensusModImpl := GetConsensusModImpl(pocketNode)
		consensusModImpl.MethodByName("SetLeaderId").Call([]reflect.Value{reflect.Zero(reflect.TypeOf(&leaderId))})
	}

	_, err = WaitForNetworkConsensusEvents(t, clockMock, eventsChannel, consensus.NewRound, consensus.Propose, numValidators*numValidators, 250, true)
	require.NoError(t, err)

	for pocketId, pocketNode := range pocketNodes {
		nodeState := GetConsensusNodeState(pocketNode)
		assertNodeConsensusView(t, pocketId,
			typesCons.ConsensusNodeState{
				Height: testHeight,
				Step:   uint8(consensus.NewRound),
				Round:  uint8(testRound + 1),
			},
			nodeState)
		//require.Equal(t, false, nodeState.IsLeader)
		require.Equal(t, typesCons.NodeId(0), nodeState.LeaderId, "Leader should be empty")
	}

}

// TODO: Implement these tests and use them as a starting point for new ones. Consider using ChatGPT to help you out :)

func TestHotstuff4Nodes1Byzantine1Block(t *testing.T) {
	t.Skip()
}

func TestHotstuff4Nodes2Byzantine1Block(t *testing.T) {
	t.Skip()
}

func TestHotstuff4Nodes1BlockNetworkPartition(t *testing.T) {
	t.Skip()
}

func TestHotstuff4Nodes1Block4Rounds(t *testing.T) {
	t.Skip()
}
func TestHotstuff4Nodes2Blocks(t *testing.T) {
	t.Skip()
}

func TestHotstuff4Nodes2NewNodes1Block(t *testing.T) {
	t.Skip()
}

func TestHotstuff4Nodes2DroppedNodes1Block(t *testing.T) {
	t.Skip()
}

func TestHotstuff4NodesFailOnPrepare(t *testing.T) {
	t.Skip()
}

func TestHotstuff4NodesFailOnPrecommit(t *testing.T) {
	t.Skip()
}

func TestHotstuff4NodesFailOnCommit(t *testing.T) {
	t.Skip()
}

func TestHotstuff4NodesFailOnDecide(t *testing.T) {
	t.Skip()
}

func TestHotstuffValidatorWithLockedQC(t *testing.T) {
	t.Skip()
}

func TestHotstuffValidatorWithLockedQCMissingNewRoundMsg(t *testing.T) {
	t.Skip()
}
