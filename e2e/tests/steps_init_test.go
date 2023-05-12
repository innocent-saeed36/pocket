//go:build e2e

package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/regen-network/gocuke"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	pocketLogger "github.com/pokt-network/pocket/logger"
	"github.com/pokt-network/pocket/runtime/defaults"
	cryptoPocket "github.com/pokt-network/pocket/shared/crypto"
	pocketk8s "github.com/pokt-network/pocket/shared/k8s"
)

var (
	logger = pocketLogger.Global.CreateLoggerForModule("e2e")

	// validatorKeys is hydrated by the clientset with credentials for all validators.
	// validatorKeys maps validator IDs to their private key as a hex string.
	validatorKeys map[string]string
	// clientset is the kubernetes API we acquire from the user's $HOME/.kube/config
	clientset *kubernetes.Clientset
	// validator holds command results between runs and reports errors to the test suite
	validator = &validatorPod{}
	// validatorA maps to suffix ID 001 of the kube pod that we use as our control agent
)

const (
	// defines the host & port scheme that LocalNet uses for naming validators.
	// e.g. validator-001 thru validator-999
	validatorServiceURLTmpl = "validator-%s-pocket:%d"
	// validatorA maps to suffix ID 001 and is also used by the cluster-manager
	// though it has no special permissions.
	validatorA = "001"
	// validatorB maps to suffix ID 002 and receives in the Send test.
	validatorB = "002"
	chainId    = "0001"
)

type rootSuite struct {
	gocuke.TestingT
}

func init() {
	cs, err := getClientset()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to get clientset")
	}
	clientset = cs
	vkmap, err := pocketk8s.FetchValidatorPrivateKeys(clientset)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to get validator key map")
	}
	validatorKeys = vkmap
}

// TestFeatures runs the e2e tests specified in any .features files in this directory
// * This test suite assumes that a LocalNet is running that can be accessed by `kubectl`
func TestFeatures(t *testing.T) {
	gocuke.NewRunner(t, &rootSuite{}).Path("*.feature").Run()
}

// InitializeScenario registers step regexes to function handlers

func (s *rootSuite) TheUserHasAValidator() {
	res, err := s.validator.RunCommand("help")
	require.NoError(s, err)
	s.validator.result = res
}

func (s *rootSuite) TheValidatorShouldHaveExitedWithoutError() {
	require.NoError(s, validator.result.Err)
}

func (s *rootSuite) TheUserRunsTheCommand(cmd string) {
	cmds := strings.Split(cmd, " ")
	res, err := validator.RunCommand(cmds...)
	require.NoError(s, err)
	validator.result = res
}

func (s *rootSuite) TheUserShouldBeAbleToSeeStandardOutputContaining(arg1 string) {
	require.Contains(s, validator.result.Stdout, arg1)
}

func (s *rootSuite) TheUserStakesTheirValidatorWithAmountUpokt(amount int64) {
	require.NoError(s, stakeValidator(fmt.Sprintf("%d", amount)))
}

func (s *rootSuite) TheUserShouldBeAbleToUnstakeTheirValidator() {
	require.NoError(s, unstakeValidator())
}

// sends amount from validator-001 to validator-002
func (s *rootSuite) TheUserSendsUpoktToAnotherAddress(amount int64) {
	privateKey := getPrivateKey(validatorKeys, validatorA)
	valB := getPrivateKey(validatorKeys, validatorB)
	args := []string{
		"--non_interactive=true",
		"--remote_cli_url=" + rpcURL,
		"Account",
		"Send",
		privateKey.Address().String(),
		valB.Address().String(),
		fmt.Sprintf("%d", amount),
	}
	res, err := validator.RunCommand(args...)
	require.NoError(s, err)
	validator.result = res
}

// stakeValidator runs Validator stake command with the address, amount, chains..., and serviceURL provided
func stakeValidator(amount string) error {
	privateKey := getPrivateKey(validatorKeys, validatorA)
	validatorServiceUrl := fmt.Sprintf(validatorServiceURLTmpl, validatorA, defaults.DefaultP2PPort)
	args := []string{
		"--non_interactive=true",
		"--remote_cli_url=" + rpcURL,
		"Validator",
		"Stake",
		privateKey.Address().String(),
		amount,
		chainId,
		validatorServiceUrl,
	}
	res, err := validator.RunCommand(args...)
	validator.result = res
	if err != nil {
		return err
	}
	return nil
}

// unstakeValidator unstakes the Validator at the same address that stakeValidator uses
func unstakeValidator() error {
	privKey := getPrivateKey(validatorKeys, validatorA)
	args := []string{
		"--non_interactive=true",
		"--remote_cli_url=" + rpcURL,
		"Validator",
		"Unstake",
		privKey.Address().String(),
	}
	res, err := validator.RunCommand(args...)
	validator.result = res
	if err != nil {
		return err
	}
	return nil
}

// getPrivateKey generates a new keypair from the private hex key that we get from the clientset
func getPrivateKey(keyMap map[string]string, validatorId string) cryptoPocket.PrivateKey {
	privHexString := keyMap[validatorId]
	keyPair, err := cryptoPocket.CreateNewKeyFromString(privHexString, "", "")
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to extract keypair")
	}
	privateKey, err := keyPair.Unarmour("")
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to extract privkey")
	}
	return privateKey
}

// getClientset uses the default path `$HOME/.kube/config` to build a kubeconfig
// and then connects to that cluster and returns a *Clientset or an error
func getClientset() (*kubernetes.Clientset, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home dir: %w", err)
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		logger.Info().Msgf("no default kubeconfig at %s; attempting to load InClusterConfig", kubeConfigPath)
		config := inClusterConfig()
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return nil, fmt.Errorf("failed to get clientset from config: %w", err)
		}
		return clientset, nil
	}
	logger.Info().Msgf("e2e tests loaded default kubeconfig located at %s", kubeConfigPath)
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get clientset from config: %w", err)
	}
	return clientset, nil
}

func inClusterConfig() *rest.Config {
	config, err := rest.InClusterConfig()
	if err != nil {
		logger.Fatal().AnErr("inClusterConfig", err)
	}
	return config
}
