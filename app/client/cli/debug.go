package cli

import (
	"os"

	"github.com/manifoldco/promptui"
	"github.com/pokt-network/pocket/logger"
	"github.com/pokt-network/pocket/p2p"
	"github.com/pokt-network/pocket/p2p/providers/addrbook_provider"
	rpcABP "github.com/pokt-network/pocket/p2p/providers/addrbook_provider/rpc"
	"github.com/pokt-network/pocket/p2p/providers/current_height_provider"
	rpcCHP "github.com/pokt-network/pocket/p2p/providers/current_height_provider/rpc"
	"github.com/pokt-network/pocket/p2p/types"
	"github.com/pokt-network/pocket/runtime"
	"github.com/pokt-network/pocket/shared/messaging"
	"github.com/pokt-network/pocket/shared/modules"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/anypb"
)

// TECHDEBT: Lowercase variables / constants that do not need to be exported.
const (
	PromptResetToGenesis      string = "ResetToGenesis"
	PromptPrintNodeState      string = "PrintNodeState"
	PromptTriggerNextView     string = "TriggerNextView"
	PromptTogglePacemakerMode string = "TogglePacemakerMode"

	PromptShowLatestBlockInStore string = "ShowLatestBlockInStore"

	PromptSendMetadataRequest string = "MetadataRequest"
	PromptSendBlockRequest    string = "BlockRequest"
)

var (
	// A P2P module is initialized in order to broadcast a message to the local network
	p2pMod modules.P2PModule

	items = []string{
		PromptResetToGenesis,
		PromptPrintNodeState,
		PromptTriggerNextView,
		PromptTogglePacemakerMode,
		PromptShowLatestBlockInStore,
		PromptSendMetadataRequest,
		PromptSendBlockRequest,
	}

	configPath  string = getEnv("CONFIG_PATH", "build/config/config1.json")
	genesisPath string = getEnv("GENESIS_PATH", "build/config/genesis.json")
)

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func init() {
	debugCmd := NewDebugCommand()
	rootCmd.AddCommand(debugCmd)
}

func NewDebugCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "debug",
		Short: "Debug utility for rapid development",
		Args:  cobra.ExactArgs(0),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			runtimeMgr := runtime.NewManagerFromFiles(
				configPath, genesisPath,
				runtime.WithClientDebugMode(),
				runtime.WithRandomPK(),
			)

			modulesRegistry := runtimeMgr.GetBus().GetModulesRegistry()
			addressBookProvider := rpcABP.NewRPCAddrBookProvider(
				rpcABP.WithP2PConfig(
					runtimeMgr.GetConfig().P2P,
				),
			)
			modulesRegistry.RegisterModule(addressBookProvider)
			currentHeightProvider := rpcCHP.NewRPCCurrentHeightProvider()
			modulesRegistry.RegisterModule(currentHeightProvider)

			p2pM, err := p2p.Create(runtimeMgr.GetBus())
			if err != nil {
				logger.Global.Fatal().Err(err).Msg("Failed to create p2p module")
			}
			p2pMod = p2pM.(modules.P2PModule)

			if err := p2pMod.Start(); err != nil {
				logger.Global.Fatal().Err(err).Msg("Failed to start p2p module")
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// if we don't inject the bus into the command context, now here we'd have to do hacky things like:
			// - initializing everything every time a command is ran (not ideal, but it would work I suppose)
			// - having a global reference to the bus that we would pass here or... reference globally directly
			// - any other ideas?
			return runDebug(cmd, modules.Bus(nil), args)
		},
	}
}

func runDebug(cmd *cobra.Command, bus modules.Bus, args []string) (err error) {
	for {
		if selection, err := promptGetInput(); err == nil {
			handleSelect(cmd, bus, selection)
		} else {
			return err
		}
	}
}

func promptGetInput() (string, error) {
	prompt := promptui.Select{
		Label: "Select an action",
		Items: items,
		Size:  len(items),
	}

	_, result, err := prompt.Run()

	if err == promptui.ErrInterrupt {
		os.Exit(0)
	}

	if err != nil {
		logger.Global.Error().Err(err).Msg("Prompt failed")
		return "", err
	}

	return result, nil
}

func handleSelect(cmd *cobra.Command, bus modules.Bus, selection string) {
	switch selection {
	case PromptResetToGenesis:
		m := &messaging.DebugMessage{
			Action:  messaging.DebugMessageAction_DEBUG_CONSENSUS_RESET_TO_GENESIS,
			Message: nil,
		}
		broadcastDebugMessage(cmd, bus, m)
	case PromptPrintNodeState:
		m := &messaging.DebugMessage{
			Action:  messaging.DebugMessageAction_DEBUG_CONSENSUS_PRINT_NODE_STATE,
			Message: nil,
		}
		broadcastDebugMessage(cmd, bus, m)
	case PromptTriggerNextView:
		m := &messaging.DebugMessage{
			Action:  messaging.DebugMessageAction_DEBUG_CONSENSUS_TRIGGER_NEXT_VIEW,
			Message: nil,
		}
		broadcastDebugMessage(cmd, bus, m)
	case PromptTogglePacemakerMode:
		m := &messaging.DebugMessage{
			Action:  messaging.DebugMessageAction_DEBUG_CONSENSUS_TOGGLE_PACE_MAKER_MODE,
			Message: nil,
		}
		broadcastDebugMessage(cmd, bus, m)
	case PromptShowLatestBlockInStore:
		m := &messaging.DebugMessage{
			Action:  messaging.DebugMessageAction_DEBUG_SHOW_LATEST_BLOCK_IN_STORE,
			Message: nil,
		}
		sendDebugMessage(cmd, bus, m)
	case PromptSendMetadataRequest:
		m := &messaging.DebugMessage{
			Action:  messaging.DebugMessageAction_DEBUG_CONSENSUS_SEND_METADATA_REQ,
			Message: nil,
		}
		broadcastDebugMessage(cmd, bus, m)
	case PromptSendBlockRequest:
		m := &messaging.DebugMessage{
			Action:  messaging.DebugMessageAction_DEBUG_CONSENSUS_SEND_BLOCK_REQ,
			Message: nil,
		}
		broadcastDebugMessage(cmd, bus, m)
	default:
		logger.Global.Error().Msg("Selection not yet implemented...")
	}
}

// Broadcast to the entire validator set
func broadcastDebugMessage(cmd *cobra.Command, bus modules.Bus, debugMsg *messaging.DebugMessage) {
	anyProto, err := anypb.New(debugMsg)
	if err != nil {
		logger.Global.Fatal().Err(err).Msg("Failed to create Any proto")
	}

	// TODO(olshansky): Once we implement the cleanup layer in RainTree, we'll be able to use
	// broadcast. The reason it cannot be done right now is because this client is not in the
	// address book of the actual validator nodes, so `node1.consensus` never receives the message.
	// p2pMod.Broadcast(anyProto)

	addrBook, err := fetchAddressBook(cmd, bus)
	for _, val := range addrBook {
		addr := val.Address
		if err != nil {
			logger.Global.Fatal().Err(err).Msg("Failed to convert validator address into pocketCrypto.Address")
		}
		if err := p2pMod.Send(addr, anyProto); err != nil {
			logger.Global.Fatal().Err(err).Msg("Failed to send debug message")
		}
	}

}

// Send to just a single (i.e. first) validator in the set
func sendDebugMessage(cmd *cobra.Command, bus modules.Bus, debugMsg *messaging.DebugMessage) {
	anyProto, err := anypb.New(debugMsg)
	if err != nil {
		logger.Global.Error().Err(err).Msg("Failed to create Any proto")
	}

	addrBook, err := fetchAddressBook(cmd, bus)
	if err != nil {
		logger.Global.Fatal().Msg("Unable to retrieve the addrBook")
	}

	var validatorAddress []byte
	if len(addrBook) == 0 {
		logger.Global.Fatal().Msg("No validators found")
	}

	// if the message needs to be broadcast, it'll be handled by the business logic of the message handler
	validatorAddress = addrBook[0].Address
	if err != nil {
		logger.Global.Fatal().Err(err).Msg("Failed to convert validator address into pocketCrypto.Address")
	}

	if err := p2pMod.Send(validatorAddress, anyProto); err != nil {
		logger.Global.Fatal().Err(err).Msg("Failed to send debug message")
	}
}

// fetchAddressBook retrieves the providers from the CLI context and uses them to retrieve the address book for the current height
func fetchAddressBook(cmd *cobra.Command, bus modules.Bus) (types.AddrBook, error) {
	modulesRegistry := bus.GetModulesRegistry()
	addrBookProvider, err := modulesRegistry.GetModule(addrbook_provider.ModuleName)
	if err != nil {
		logger.Global.Fatal().Msg("Unable to retrieve the addrBookProvider")
	}
	currentHeightProvider, err := modulesRegistry.GetModule(current_height_provider.ModuleName)
	if err != nil {
		logger.Global.Fatal().Msg("Unable to retrieve the currentHeightProvider")
	}

	height := currentHeightProvider.(current_height_provider.CurrentHeightProvider).CurrentHeight()
	addrBook, err := addrBookProvider.(addrbook_provider.AddrBookProvider).GetStakedAddrBookAtHeight(height)
	if err != nil {
		logger.Global.Fatal().Msg("Unable to retrieve the addrBook")
	}
	return addrBook, err
}
