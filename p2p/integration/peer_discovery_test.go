package integration

import (
	"testing"

	"github.com/cucumber/godog"

	"github.com/pokt-network/pocket/internal/testutil"
)

const peerDiscoveryFeaturePath = "peer_discovery.feature"

func TestPeerDiscoveryIntegration(t *testing.T) {
	t.Parallel()

	testutil.RunGherkinFeature(t, peerDiscoveryFeaturePath, initPeerDiscoveryScenarios)
}

func aNode(arg1 string) error {
	return godog.ErrPending
}

func aNodeInPartition(arg1, arg2 string) error {
	return godog.ErrPending
}

func aNodeJoinsPartitionsAnd(arg1, arg2, arg3 string) error {
	return godog.ErrPending
}

func allNodesInPartitionShouldDiscoverAllNodesInPartition(arg1, arg2 string) error {
	return godog.ErrPending
}

func eachNodeShouldHaveNumberOfPeersInTheirRespectivePeerstores(arg1 int) error {
	return godog.ErrPending
}

func eachNodeShouldNotHaveAnyNodesWhichLeftInTheirPeerstores() error {
	return godog.ErrPending
}

func numberOfNodesBootstrapInPartition(arg1 int, arg2 string) error {
	return godog.ErrPending
}

func otherNodesShouldHaveNumberOfPeersInTheirPeerstores(arg1 int) error {
	return godog.ErrPending
}

func otherNodesShouldNotBeIncludedInTheirRespectivePeerstores() error {
	return godog.ErrPending
}

func theNetworkShouldContainNumberOfNodes(arg1 int) error {
	return godog.ErrPending
}

func theNodeLeavesTheNetwork(arg1 string) error {
	return godog.ErrPending
}

func theNodeShouldHaveNumberOfPeersInItsPeerstore(arg1 string, arg2 int) error {
	return godog.ErrPending
}

func theNodeShouldNotBeIncludedInItsOwnPeerstore(arg1 string) error {
	return godog.ErrPending
}

func initPeerDiscoveryScenarios(ctx *godog.ScenarioContext) {
	ctx.Step(`^a "([^"]*)" node$`, aNode)
	ctx.Step(`^a "([^"]*)" node in partition "([^"]*)"$`, aNodeInPartition)
	ctx.Step(`^a "([^"]*)" node joins partitions "([^"]*)" and "([^"]*)"$`, aNodeJoinsPartitionsAnd)
	ctx.Step(`^all nodes in partition "([^"]*)" should discover all nodes in partition "([^"]*)"$`, allNodesInPartitionShouldDiscoverAllNodesInPartition)
	ctx.Step(`^each node should have (\d+) number of peers in their respective peerstores$`, eachNodeShouldHaveNumberOfPeersInTheirRespectivePeerstores)
	ctx.Step(`^each node should not have any nodes which left in their peerstores$`, eachNodeShouldNotHaveAnyNodesWhichLeftInTheirPeerstores)
	ctx.Step(`^(\d+) number of nodes bootstrap in partition "([^"]*)"$`, numberOfNodesBootstrapInPartition)
	ctx.Step(`^(\d+) number of nodes join the network$`, numberOfNodesJoinTheNetwork)
	ctx.Step(`^(\d+) number of nodes leave the network$`, numberOfNodesLeaveTheNetwork)
	ctx.Step(`^other nodes should have (\d+) number of peers in their peerstores$`, otherNodesShouldHaveNumberOfPeersInTheirPeerstores)
	ctx.Step(`^other nodes should not be included in their respective peerstores$`, otherNodesShouldNotBeIncludedInTheirRespectivePeerstores)
	ctx.Step(`^the network should contain (\d+) number of nodes$`, theNetworkShouldContainNumberOfNodes)
	ctx.Step(`^the "([^"]*)" node leaves the network$`, theNodeLeavesTheNetwork)
	ctx.Step(`^the "([^"]*)" node should have (\d+) number of peers in its peerstore$`, theNodeShouldHaveNumberOfPeersInItsPeerstore)
	ctx.Step(`^the "([^"]*)" node should not be included in its own peerstore$`, theNodeShouldNotBeIncludedInItsOwnPeerstore)
}
