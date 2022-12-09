package raintree

import (
	"log"
	"math"
	"strings"
)

// Refer to the P2P specification for a formal description and proof of how the constants are selected
const (
	firstMsgTargetPercentage  = float64(1) / float64(3)
	secondMsgTargetPercentage = float64(2) / float64(3)
	shrinkagePercentage       = float64(2) / float64(3)
	maxLevelsLogBase          = float64(3)
	floatPrecision            = float64(0.0000001)
)

func (n *rainTreeNetwork) getAddrBookLength(level uint32, height uint64) int {
	peersManagerStateView := n.peersManager.getNetworkView()

	// if we are propagating a message from a previous height, we need to instantiate an ephemeral peersManager (without add/remove)
	if height < n.GetBus().GetConsensusModule().CurrentHeight() {
		peersManagerWithAddrBookProvider, err := NewPeersManagerWithAddrBookProvider(n.selfAddr, n.addrBookProvider, height)
		if err != nil {
			log.Fatalf("[ERROR] Error initializing rainTreeNetwork peersManagerWithAddrBookProvider: %v", err)
		}
		peersManagerStateView = peersManagerWithAddrBookProvider.getNetworkView()
	}

	shrinkageCoefficient := math.Pow(shrinkagePercentage, float64(peersManagerStateView.maxNumLevels-level))
	return int(float64(len(peersManagerStateView.addrList)) * (shrinkageCoefficient))
}

// getTargetsAtLevel returns the targets for a given level
func (n *rainTreeNetwork) getTargetsAtLevel(level uint32) []target {
	height := n.GetBus().GetConsensusModule().CurrentHeight()
	addrBookLengthAtHeight := n.getAddrBookLength(level, height)
	firstTarget := n.getTarget(firstMsgTargetPercentage, addrBookLengthAtHeight, level)
	secondTarget := n.getTarget(secondMsgTargetPercentage, addrBookLengthAtHeight, level)

	log.Printf("[DEBUG] Targets at height (%d): %s", level, n.debugMsgTargetString(firstTarget, secondTarget))

	return []target{firstTarget, secondTarget}
}

func (n *rainTreeNetwork) getTarget(targetPercentage float64, addrBookLen int, level uint32) target {
	i := int(targetPercentage * float64(addrBookLen))

	peersManagerStateView := n.peersManager.getNetworkView()

	target := target{
		serviceUrl:             peersManagerStateView.addrBookMap[peersManagerStateView.addrList[i]].ServiceUrl,
		percentage:             targetPercentage,
		level:                  level,
		addrBookLengthAtHeight: addrBookLen,
		index:                  i,
	}

	// If the target is 0, it is a reference to self, which is a `Demote` in RainTree terms.
	// This is handled separately.
	if i == 0 {
		target.isSelf = true
		return target
	}

	addrStr := peersManagerStateView.addrList[i]
	if addr, ok := peersManagerStateView.addrBookMap[addrStr]; ok {
		target.address = addr.Address
		return target
	}
	log.Printf("[DEBUG] addrStr %s not found in addrBookMap", addrStr)
	return target
}

// Only used for debug logging to understand what RainTree is doing under the hood
func (n *rainTreeNetwork) debugMsgTargetString(target1, target2 target) string {
	s := strings.Builder{}
	s.WriteString(target1.DebugString(n))
	s.WriteString(" --|-- ")
	s.WriteString(target2.DebugString(n))
	return s.String()
}
