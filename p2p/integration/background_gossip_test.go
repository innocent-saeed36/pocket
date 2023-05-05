package integration

import (
	"testing"

	"github.com/cucumber/godog"

	"github.com/pokt-network/pocket/internal/testutil"
)

const backgroundGossipFeaturePath = "background_gossip.feature"

func TestBackgroundGossipIntegration(t *testing.T) {
	t.Parallel()

	testutil.RunGherkinFeature(t, backgroundGossipFeaturePath, initBackgroundGossipScenarios)
}
func aFaultyNetworkOfPeers(arg1 int) error {
	return godog.ErrPending
}

func aFullyConnectedNetworkOfPeers(arg1 int) error {
	return godog.ErrPending
}

func aNodeBroadcastsATestMessageViaItsBackgroundRouter() error {
	return godog.ErrPending
}

func numberOfFaultyPeers(arg1 int) error {
	return godog.ErrPending
}

func numberOfNodesShouldReceiveTheTestMessage(arg1 int) error {
	return godog.ErrPending
}

func initBackgroundGossipScenarios(ctx *godog.ScenarioContext) {
	ctx.Step(`^a faulty network of (\d+) peers$`, aFaultyNetworkOfPeers)
	ctx.Step(`^a fully connected network of (\d+) peers$`, aFullyConnectedNetworkOfPeers)
	ctx.Step(`^a node broadcasts a test message via its background router$`, aNodeBroadcastsATestMessageViaItsBackgroundRouter)
	ctx.Step(`^(\d+) number of faulty peers$`, numberOfFaultyPeers)
	ctx.Step(`^(\d+) number of nodes join the network$`, numberOfNodesJoinTheNetwork)
	ctx.Step(`^(\d+) number of nodes leave the network$`, numberOfNodesLeaveTheNetwork)
	ctx.Step(`^(\d+) number of nodes should receive the test message$`, numberOfNodesShouldReceiveTheTestMessage)
}
