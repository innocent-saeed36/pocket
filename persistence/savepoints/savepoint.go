package savepoints

import (
	"github.com/pokt-network/pocket/logger"
	"github.com/pokt-network/pocket/persistence/kvstore"
	coreTypes "github.com/pokt-network/pocket/shared/core/types"
	"github.com/pokt-network/pocket/shared/modules"
	moduleTypes "github.com/pokt-network/pocket/shared/modules/types"
)

var (
	_   modules.PersistenceReadContext = &savepoint{}
	log                                = logger.Global.CreateLoggerForModule("savepoint")
)

type savepoint struct {
	height int64

	nodeStoresPath  string
	valueStoresPath string
}

func (*savepoint) GetKVStores() (nodeStores map[int]kvstore.BackupableKVStore, valueStores map[int]kvstore.BackupableKVStore) {
	// no-op
	return nil, nil
}

func (s *savepoint) Close() error {
	// no-op
	return nil
}

func (s *savepoint) GetAccountAmount(address []byte, height int64) (string, error) {
	panic("unimplemented")
}

func (s *savepoint) GetAllAccounts(height int64) (accs []*coreTypes.Account, err error) {
	panic("unimplemented")
}

// GetAllApps implements modules.PersistenceReadContext
func (*savepoint) GetAllApps(height int64) ([]*coreTypes.Actor, error) {
	panic("unimplemented")
}

// GetAllFishermen implements modules.PersistenceReadContext
func (*savepoint) GetAllFishermen(height int64) ([]*coreTypes.Actor, error) {
	panic("unimplemented")
}

// GetAllPools implements modules.PersistenceReadContext
func (*savepoint) GetAllPools(height int64) ([]*coreTypes.Account, error) {
	panic("unimplemented")
}

// GetAllServicers implements modules.PersistenceReadContext
func (*savepoint) GetAllServicers(height int64) ([]*coreTypes.Actor, error) {
	panic("unimplemented")
}

// GetAllStakedActors implements modules.PersistenceReadContext
func (*savepoint) GetAllStakedActors(height int64) ([]*coreTypes.Actor, error) {
	panic("unimplemented")
}

// GetAllValidators implements modules.PersistenceReadContext
func (*savepoint) GetAllValidators(height int64) ([]*coreTypes.Actor, error) {
	panic("unimplemented")
}

// GetAppExists implements modules.PersistenceReadContext
func (*savepoint) GetAppExists(address []byte, height int64) (exists bool, err error) {
	panic("unimplemented")
}

// GetAppOutputAddress implements modules.PersistenceReadContext
func (*savepoint) GetAppOutputAddress(operator []byte, height int64) (output []byte, err error) {
	panic("unimplemented")
}

// GetAppPauseHeightIfExists implements modules.PersistenceReadContext
func (*savepoint) GetAppPauseHeightIfExists(address []byte, height int64) (int64, error) {
	panic("unimplemented")
}

// GetAppStakeAmount implements modules.PersistenceReadContext
func (*savepoint) GetAppStakeAmount(height int64, address []byte) (string, error) {
	panic("unimplemented")
}

// GetAppStatus implements modules.PersistenceReadContext
func (*savepoint) GetAppStatus(address []byte, height int64) (status int32, err error) {
	panic("unimplemented")
}

// GetAppsReadyToUnstake implements modules.PersistenceReadContext
func (*savepoint) GetAppsReadyToUnstake(height int64, status int32) (apps []*moduleTypes.UnstakingActor, err error) {
	panic("unimplemented")
}

// GetBlockHash implements modules.PersistenceReadContext
func (*savepoint) GetBlockHash(height int64) (string, error) {
	panic("unimplemented")
}

// GetBytesFlag implements modules.PersistenceReadContext
func (*savepoint) GetBytesFlag(paramName string, height int64) ([]byte, bool, error) {
	panic("unimplemented")
}

// GetBytesParam implements modules.PersistenceReadContext
func (*savepoint) GetBytesParam(paramName string, height int64) ([]byte, error) {
	panic("unimplemented")
}

// GetFishermanExists implements modules.PersistenceReadContext
func (*savepoint) GetFishermanExists(address []byte, height int64) (exists bool, err error) {
	panic("unimplemented")
}

// GetFishermanOutputAddress implements modules.PersistenceReadContext
func (*savepoint) GetFishermanOutputAddress(operator []byte, height int64) (output []byte, err error) {
	panic("unimplemented")
}

// GetFishermanPauseHeightIfExists implements modules.PersistenceReadContext
func (*savepoint) GetFishermanPauseHeightIfExists(address []byte, height int64) (int64, error) {
	panic("unimplemented")
}

// GetFishermanStakeAmount implements modules.PersistenceReadContext
func (*savepoint) GetFishermanStakeAmount(height int64, address []byte) (string, error) {
	panic("unimplemented")
}

// GetFishermanStatus implements modules.PersistenceReadContext
func (*savepoint) GetFishermanStatus(address []byte, height int64) (status int32, err error) {
	panic("unimplemented")
}

// GetFishermenReadyToUnstake implements modules.PersistenceReadContext
func (*savepoint) GetFishermenReadyToUnstake(height int64, status int32) (fishermen []*moduleTypes.UnstakingActor, err error) {
	panic("unimplemented")
}

// GetHeight implements modules.PersistenceReadContext
func (*savepoint) GetHeight() (int64, error) {
	panic("unimplemented")
}

// GetIntFlag implements modules.PersistenceReadContext
func (*savepoint) GetIntFlag(paramName string, height int64) (int, bool, error) {
	panic("unimplemented")
}

// GetIntParam implements modules.PersistenceReadContext
func (*savepoint) GetIntParam(paramName string, height int64) (int, error) {
	panic("unimplemented")
}

// GetMaximumBlockHeight implements modules.PersistenceReadContext
func (*savepoint) GetMaximumBlockHeight() (uint64, error) {
	panic("unimplemented")
}

// GetMinimumBlockHeight implements modules.PersistenceReadContext
func (*savepoint) GetMinimumBlockHeight() (uint64, error) {
	panic("unimplemented")
}

// GetParameter implements modules.PersistenceReadContext
func (*savepoint) GetParameter(paramName string, height int64) (any, error) {
	panic("unimplemented")
}

// GetPoolAmount implements modules.PersistenceReadContext
func (*savepoint) GetPoolAmount(address []byte, height int64) (amount string, err error) {
	panic("unimplemented")
}

// GetServicerCount implements modules.PersistenceReadContext
func (*savepoint) GetServicerCount(chain string, height int64) (int, error) {
	panic("unimplemented")
}

// GetServicerExists implements modules.PersistenceReadContext
func (*savepoint) GetServicerExists(address []byte, height int64) (exists bool, err error) {
	panic("unimplemented")
}

// GetServicerOutputAddress implements modules.PersistenceReadContext
func (*savepoint) GetServicerOutputAddress(operator []byte, height int64) (output []byte, err error) {
	panic("unimplemented")
}

// GetServicerPauseHeightIfExists implements modules.PersistenceReadContext
func (*savepoint) GetServicerPauseHeightIfExists(address []byte, height int64) (int64, error) {
	panic("unimplemented")
}

// GetServicerStakeAmount implements modules.PersistenceReadContext
func (*savepoint) GetServicerStakeAmount(height int64, address []byte) (string, error) {
	panic("unimplemented")
}

// GetServicerStatus implements modules.PersistenceReadContext
func (*savepoint) GetServicerStatus(address []byte, height int64) (status int32, err error) {
	panic("unimplemented")
}

// GetServicersReadyToUnstake implements modules.PersistenceReadContext
func (*savepoint) GetServicersReadyToUnstake(height int64, status int32) (servicers []*moduleTypes.UnstakingActor, err error) {
	panic("unimplemented")
}

// GetStringFlag implements modules.PersistenceReadContext
func (*savepoint) GetStringFlag(paramName string, height int64) (string, bool, error) {
	panic("unimplemented")
}

// GetStringParam implements modules.PersistenceReadContext
func (*savepoint) GetStringParam(paramName string, height int64) (string, error) {
	panic("unimplemented")
}

// GetValidatorExists implements modules.PersistenceReadContext
func (*savepoint) GetValidatorExists(address []byte, height int64) (exists bool, err error) {
	panic("unimplemented")
}

// GetValidatorMissedBlocks implements modules.PersistenceReadContext
func (*savepoint) GetValidatorMissedBlocks(address []byte, height int64) (int, error) {
	panic("unimplemented")
}

// GetValidatorOutputAddress implements modules.PersistenceReadContext
func (*savepoint) GetValidatorOutputAddress(operator []byte, height int64) (output []byte, err error) {
	panic("unimplemented")
}

// GetValidatorPauseHeightIfExists implements modules.PersistenceReadContext
func (*savepoint) GetValidatorPauseHeightIfExists(address []byte, height int64) (int64, error) {
	panic("unimplemented")
}

// GetValidatorStakeAmount implements modules.PersistenceReadContext
func (*savepoint) GetValidatorStakeAmount(height int64, address []byte) (string, error) {
	panic("unimplemented")
}

// GetValidatorStatus implements modules.PersistenceReadContext
func (*savepoint) GetValidatorStatus(address []byte, height int64) (status int32, err error) {
	panic("unimplemented")
}

// GetValidatorsReadyToUnstake implements modules.PersistenceReadContext
func (*savepoint) GetValidatorsReadyToUnstake(height int64, status int32) (validators []*moduleTypes.UnstakingActor, err error) {
	panic("unimplemented")
}

func (*savepoint) Release() {
	// no-op
}
