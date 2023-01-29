package utility

import (
	"log"
	"math/big"

	"github.com/pokt-network/pocket/shared/converters"
	coreTypes "github.com/pokt-network/pocket/shared/core/types"
	typesUtil "github.com/pokt-network/pocket/utility/types"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (u *utilityContext) updateParam(paramName string, value any) typesUtil.Error {
	store := u.Store()
	switch t := value.(type) {
	case *wrapperspb.Int32Value:
		if err := store.SetParam(paramName, (int(t.Value))); err != nil {
			return typesUtil.ErrUpdateParam(err)
		}
		return nil
	case *wrapperspb.StringValue:
		if err := store.SetParam(paramName, t.Value); err != nil {
			return typesUtil.ErrUpdateParam(err)
		}
		return nil
	case *wrapperspb.BytesValue:
		if err := store.SetParam(paramName, t.Value); err != nil {
			return typesUtil.ErrUpdateParam(err)
		}
		return nil
	default:
		break
	}
	log.Fatalf("unhandled value type %T for %v", value, value)
	return typesUtil.ErrUnknownParam(paramName)
}

func (u *utilityContext) GetParameter(paramName string, height int64) (any, error) {
	return u.Store().GetParameter(paramName, height)
}

func (u *utilityContext) GetAppMinimumStake() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.AppMinimumStakeParamName)
}

func (u *utilityContext) GetAppMaxChains() (int, typesUtil.Error) {
	return u.getIntParam(typesUtil.AppMaxChainsParamName)
}

func (u *utilityContext) GetBaselineAppStakeRate() (int, typesUtil.Error) {
	return u.getIntParam(typesUtil.AppBaselineStakeRateParamName)
}

func (u *utilityContext) GetStabilityAdjustment() (int, typesUtil.Error) {
	return u.getIntParam(typesUtil.AppStakingAdjustmentParamName)
}

func (u *utilityContext) GetAppUnstakingBlocks() (int64, typesUtil.Error) {
	return u.getInt64Param(typesUtil.AppUnstakingBlocksParamName)
}

func (u *utilityContext) GetAppMinimumPauseBlocks() (int, typesUtil.Error) {
	return u.getIntParam(typesUtil.AppMinimumPauseBlocksParamName)
}

func (u *utilityContext) GetAppMaxPausedBlocks() (maxPausedBlocks int, err typesUtil.Error) {
	return u.getIntParam(typesUtil.AppMaxPauseBlocksParamName)
}

func (u *utilityContext) GetServiceNodeMinimumStake() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.ServiceNodeMinimumStakeParamName)
}

func (u *utilityContext) GetServiceNodeMaxChains() (int, typesUtil.Error) {
	return u.getIntParam(typesUtil.ServiceNodeMaxChainsParamName)
}

func (u *utilityContext) GetServiceNodeUnstakingBlocks() (int64, typesUtil.Error) {
	return u.getInt64Param(typesUtil.ServiceNodeUnstakingBlocksParamName)
}

func (u *utilityContext) GetServiceNodeMinimumPauseBlocks() (int, typesUtil.Error) {
	return u.getIntParam(typesUtil.ServiceNodeMinimumPauseBlocksParamName)
}

func (u *utilityContext) GetServiceNodeMaxPausedBlocks() (maxPausedBlocks int, err typesUtil.Error) {
	return u.getIntParam(typesUtil.ServiceNodeMaxPauseBlocksParamName)
}

func (u *utilityContext) getValidatorMinimumStake() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.ValidatorMinimumStakeParamName)
}

func (u *utilityContext) GetValidatorUnstakingBlocks() (int64, typesUtil.Error) {
	return u.getInt64Param(typesUtil.ValidatorUnstakingBlocksParamName)
}

func (u *utilityContext) GetValidatorMinimumPauseBlocks() (int, typesUtil.Error) {
	return u.getIntParam(typesUtil.ValidatorMinimumPauseBlocksParamName)
}

func (u *utilityContext) GetValidatorMaxPausedBlocks() (maxPausedBlocks int, err typesUtil.Error) {
	return u.getIntParam(typesUtil.ValidatorMaxPausedBlocksParamName)
}

func (u *utilityContext) GetProposerPercentageOfFees() (proposerPercentage int, err typesUtil.Error) {
	return u.getIntParam(typesUtil.ProposerPercentageOfFeesParamName)
}

func (u *utilityContext) GetValidatorMaxMissedBlocks() (maxMissedBlocks int, err typesUtil.Error) {
	return u.getIntParam(typesUtil.ValidatorMaximumMissedBlocksParamName)
}

func (u *utilityContext) GetMaxEvidenceAgeInBlocks() (maxMissedBlocks int, err typesUtil.Error) {
	return u.getIntParam(typesUtil.ValidatorMaxEvidenceAgeInBlocksParamName)
}

func (u *utilityContext) GetDoubleSignBurnPercentage() (burnPercentage int, err typesUtil.Error) {
	return u.getIntParam(typesUtil.DoubleSignBurnPercentageParamName)
}

func (u *utilityContext) GetMissedBlocksBurnPercentage() (burnPercentage int, err typesUtil.Error) {
	return u.getIntParam(typesUtil.MissedBlocksBurnPercentageParamName)
}

func (u *utilityContext) GetFishermanMinimumStake() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.FishermanMinimumStakeParamName)
}

func (u *utilityContext) GetFishermanMaxChains() (int, typesUtil.Error) {
	return u.getIntParam(typesUtil.FishermanMaxChainsParamName)
}

func (u *utilityContext) GetFishermanUnstakingBlocks() (int64, typesUtil.Error) {
	return u.getInt64Param(typesUtil.FishermanUnstakingBlocksParamName)
}

func (u *utilityContext) GetFishermanMinimumPauseBlocks() (int, typesUtil.Error) {
	return u.getIntParam(typesUtil.FishermanMinimumPauseBlocksParamName)
}

func (u *utilityContext) GetFishermanMaxPausedBlocks() (maxPausedBlocks int, err typesUtil.Error) {
	return u.getIntParam(typesUtil.FishermanMaxPauseBlocksParamName)
}

func (u *utilityContext) GetMessageDoubleSignFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessageDoubleSignFee)
}

func (u *utilityContext) GetMessageSendFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessageSendFee)
}

func (u *utilityContext) GetMessageStakeFishermanFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessageStakeFishermanFee)
}

func (u *utilityContext) GetMessageEditStakeFishermanFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessageEditStakeFishermanFee)
}

func (u *utilityContext) GetMessageUnstakeFishermanFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessageUnstakeFishermanFee)
}

func (u *utilityContext) GetMessagePauseFishermanFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessagePauseFishermanFee)
}

func (u *utilityContext) GetMessageUnpauseFishermanFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessageUnpauseFishermanFee)
}

func (u *utilityContext) GetMessageFishermanPauseServiceNodeFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessageFishermanPauseServiceNodeFee)
}

func (u *utilityContext) GetMessageTestScoreFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessageTestScoreFee)
}

func (u *utilityContext) GetMessageProveTestScoreFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessageProveTestScoreFee)
}

func (u *utilityContext) GetMessageStakeAppFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessageStakeAppFee)
}

func (u *utilityContext) GetMessageEditStakeAppFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessageEditStakeAppFee)
}

func (u *utilityContext) GetMessageUnstakeAppFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessageUnstakeAppFee)
}

func (u *utilityContext) GetMessagePauseAppFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessagePauseAppFee)
}

func (u *utilityContext) GetMessageUnpauseAppFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessageUnpauseAppFee)
}

func (u *utilityContext) GetMessageStakeValidatorFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessageStakeValidatorFee)
}

func (u *utilityContext) GetMessageEditStakeValidatorFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessageEditStakeValidatorFee)
}

func (u *utilityContext) GetMessageUnstakeValidatorFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessageUnstakeValidatorFee)
}

func (u *utilityContext) GetMessagePauseValidatorFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessagePauseValidatorFee)
}

func (u *utilityContext) GetMessageUnpauseValidatorFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessageUnpauseValidatorFee)
}

func (u *utilityContext) GetMessageStakeServiceNodeFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessageStakeServiceNodeFee)
}

func (u *utilityContext) GetMessageEditStakeServiceNodeFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessageEditStakeServiceNodeFee)
}

func (u *utilityContext) GetMessageUnstakeServiceNodeFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessageUnstakeServiceNodeFee)
}

func (u *utilityContext) GetMessagePauseServiceNodeFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessagePauseServiceNodeFee)
}

func (u *utilityContext) GetMessageUnpauseServiceNodeFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessageUnpauseServiceNodeFee)
}

func (u *utilityContext) GetMessageChangeParameterFee() (*big.Int, typesUtil.Error) {
	return u.getBigIntParam(typesUtil.MessageChangeParameterFee)
}

func (u *utilityContext) GetDoubleSignFeeOwner() (owner []byte, err typesUtil.Error) {
	return u.getByteArrayParam(typesUtil.MessageDoubleSignFeeOwner)
}

func (u *utilityContext) GetParamOwner(paramName string) ([]byte, error) {
	// DISCUSS (@deblasis): here we could potentially leverage the struct tags in gov.proto by specifying an `owner` key
	// eg: `app_minimum_stake` could have `pokt:"owner=app_minimum_stake_owner"`
	// in here we would use that map to point to the owner, removing this switch, centralizing the logic and making it declarative
	store, height, er := u.getStoreAndHeight()
	if er != nil {
		return nil, er
	}
	switch paramName {
	case typesUtil.AclOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.BlocksPerSessionParamName:
		return store.GetBytesParam(typesUtil.BlocksPerSessionOwner, height)
	case typesUtil.AppMaxChainsParamName:
		return store.GetBytesParam(typesUtil.AppMaxChainsOwner, height)
	case typesUtil.AppMinimumStakeParamName:
		return store.GetBytesParam(typesUtil.AppMinimumStakeOwner, height)
	case typesUtil.AppBaselineStakeRateParamName:
		return store.GetBytesParam(typesUtil.AppBaselineStakeRateOwner, height)
	case typesUtil.AppStakingAdjustmentParamName:
		return store.GetBytesParam(typesUtil.AppStakingAdjustmentOwner, height)
	case typesUtil.AppUnstakingBlocksParamName:
		return store.GetBytesParam(typesUtil.AppUnstakingBlocksOwner, height)
	case typesUtil.AppMinimumPauseBlocksParamName:
		return store.GetBytesParam(typesUtil.AppMinimumPauseBlocksOwner, height)
	case typesUtil.AppMaxPauseBlocksParamName:
		return store.GetBytesParam(typesUtil.AppMaxPausedBlocksOwner, height)
	case typesUtil.ServiceNodesPerSessionParamName:
		return store.GetBytesParam(typesUtil.ServiceNodesPerSessionOwner, height)
	case typesUtil.ServiceNodeMinimumStakeParamName:
		return store.GetBytesParam(typesUtil.ServiceNodeMinimumStakeOwner, height)
	case typesUtil.ServiceNodeMaxChainsParamName:
		return store.GetBytesParam(typesUtil.ServiceNodeMaxChainsOwner, height)
	case typesUtil.ServiceNodeUnstakingBlocksParamName:
		return store.GetBytesParam(typesUtil.ServiceNodeUnstakingBlocksOwner, height)
	case typesUtil.ServiceNodeMinimumPauseBlocksParamName:
		return store.GetBytesParam(typesUtil.ServiceNodeMinimumPauseBlocksOwner, height)
	case typesUtil.ServiceNodeMaxPauseBlocksParamName:
		return store.GetBytesParam(typesUtil.ServiceNodeMaxPausedBlocksOwner, height)
	case typesUtil.FishermanMinimumStakeParamName:
		return store.GetBytesParam(typesUtil.FishermanMinimumStakeOwner, height)
	case typesUtil.FishermanMaxChainsParamName:
		return store.GetBytesParam(typesUtil.FishermanMaxChainsOwner, height)
	case typesUtil.FishermanUnstakingBlocksParamName:
		return store.GetBytesParam(typesUtil.FishermanUnstakingBlocksOwner, height)
	case typesUtil.FishermanMinimumPauseBlocksParamName:
		return store.GetBytesParam(typesUtil.FishermanMinimumPauseBlocksOwner, height)
	case typesUtil.FishermanMaxPauseBlocksParamName:
		return store.GetBytesParam(typesUtil.FishermanMaxPausedBlocksOwner, height)
	case typesUtil.ValidatorMinimumStakeParamName:
		return store.GetBytesParam(typesUtil.ValidatorMinimumStakeOwner, height)
	case typesUtil.ValidatorUnstakingBlocksParamName:
		return store.GetBytesParam(typesUtil.ValidatorUnstakingBlocksOwner, height)
	case typesUtil.ValidatorMinimumPauseBlocksParamName:
		return store.GetBytesParam(typesUtil.ValidatorMinimumPauseBlocksOwner, height)
	case typesUtil.ValidatorMaxPausedBlocksParamName:
		return store.GetBytesParam(typesUtil.ValidatorMaxPausedBlocksOwner, height)
	case typesUtil.ValidatorMaximumMissedBlocksParamName:
		return store.GetBytesParam(typesUtil.ValidatorMaximumMissedBlocksOwner, height)
	case typesUtil.ProposerPercentageOfFeesParamName:
		return store.GetBytesParam(typesUtil.ProposerPercentageOfFeesOwner, height)
	case typesUtil.ValidatorMaxEvidenceAgeInBlocksParamName:
		return store.GetBytesParam(typesUtil.ValidatorMaxEvidenceAgeInBlocksOwner, height)
	case typesUtil.MissedBlocksBurnPercentageParamName:
		return store.GetBytesParam(typesUtil.MissedBlocksBurnPercentageOwner, height)
	case typesUtil.DoubleSignBurnPercentageParamName:
		return store.GetBytesParam(typesUtil.DoubleSignBurnPercentageOwner, height)
	case typesUtil.MessageDoubleSignFee:
		return store.GetBytesParam(typesUtil.MessageDoubleSignFeeOwner, height)
	case typesUtil.MessageSendFee:
		return store.GetBytesParam(typesUtil.MessageSendFeeOwner, height)
	case typesUtil.MessageStakeFishermanFee:
		return store.GetBytesParam(typesUtil.MessageStakeFishermanFeeOwner, height)
	case typesUtil.MessageEditStakeFishermanFee:
		return store.GetBytesParam(typesUtil.MessageEditStakeFishermanFeeOwner, height)
	case typesUtil.MessageUnstakeFishermanFee:
		return store.GetBytesParam(typesUtil.MessageUnstakeFishermanFeeOwner, height)
	case typesUtil.MessagePauseFishermanFee:
		return store.GetBytesParam(typesUtil.MessagePauseFishermanFeeOwner, height)
	case typesUtil.MessageUnpauseFishermanFee:
		return store.GetBytesParam(typesUtil.MessageUnpauseFishermanFeeOwner, height)
	case typesUtil.MessageFishermanPauseServiceNodeFee:
		return store.GetBytesParam(typesUtil.MessageFishermanPauseServiceNodeFeeOwner, height)
	case typesUtil.MessageTestScoreFee:
		return store.GetBytesParam(typesUtil.MessageTestScoreFeeOwner, height)
	case typesUtil.MessageProveTestScoreFee:
		return store.GetBytesParam(typesUtil.MessageProveTestScoreFeeOwner, height)
	case typesUtil.MessageStakeAppFee:
		return store.GetBytesParam(typesUtil.MessageStakeAppFeeOwner, height)
	case typesUtil.MessageEditStakeAppFee:
		return store.GetBytesParam(typesUtil.MessageEditStakeAppFeeOwner, height)
	case typesUtil.MessageUnstakeAppFee:
		return store.GetBytesParam(typesUtil.MessageUnstakeAppFeeOwner, height)
	case typesUtil.MessagePauseAppFee:
		return store.GetBytesParam(typesUtil.MessagePauseAppFeeOwner, height)
	case typesUtil.MessageUnpauseAppFee:
		return store.GetBytesParam(typesUtil.MessageUnpauseAppFeeOwner, height)
	case typesUtil.MessageStakeValidatorFee:
		return store.GetBytesParam(typesUtil.MessageStakeValidatorFeeOwner, height)
	case typesUtil.MessageEditStakeValidatorFee:
		return store.GetBytesParam(typesUtil.MessageEditStakeValidatorFeeOwner, height)
	case typesUtil.MessageUnstakeValidatorFee:
		return store.GetBytesParam(typesUtil.MessageUnstakeValidatorFeeOwner, height)
	case typesUtil.MessagePauseValidatorFee:
		return store.GetBytesParam(typesUtil.MessagePauseValidatorFeeOwner, height)
	case typesUtil.MessageUnpauseValidatorFee:
		return store.GetBytesParam(typesUtil.MessageUnpauseValidatorFeeOwner, height)
	case typesUtil.MessageStakeServiceNodeFee:
		return store.GetBytesParam(typesUtil.MessageStakeServiceNodeFeeOwner, height)
	case typesUtil.MessageEditStakeServiceNodeFee:
		return store.GetBytesParam(typesUtil.MessageEditStakeServiceNodeFeeOwner, height)
	case typesUtil.MessageUnstakeServiceNodeFee:
		return store.GetBytesParam(typesUtil.MessageUnstakeServiceNodeFeeOwner, height)
	case typesUtil.MessagePauseServiceNodeFee:
		return store.GetBytesParam(typesUtil.MessagePauseServiceNodeFeeOwner, height)
	case typesUtil.MessageUnpauseServiceNodeFee:
		return store.GetBytesParam(typesUtil.MessageUnpauseServiceNodeFeeOwner, height)
	case typesUtil.MessageChangeParameterFee:
		return store.GetBytesParam(typesUtil.MessageChangeParameterFeeOwner, height)
	case typesUtil.BlocksPerSessionOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.AppMaxChainsOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.AppMinimumStakeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.AppBaselineStakeRateOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.AppStakingAdjustmentOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.AppUnstakingBlocksOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.AppMinimumPauseBlocksOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.AppMaxPausedBlocksOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.ServiceNodeMinimumStakeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.ServiceNodeMaxChainsOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.ServiceNodeUnstakingBlocksOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.ServiceNodeMinimumPauseBlocksOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.ServiceNodeMaxPausedBlocksOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.ServiceNodesPerSessionOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.FishermanMinimumStakeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.FishermanMaxChainsOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.FishermanUnstakingBlocksOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.FishermanMinimumPauseBlocksOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.FishermanMaxPausedBlocksOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.ValidatorMinimumStakeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.ValidatorUnstakingBlocksOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.ValidatorMinimumPauseBlocksOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.ValidatorMaxPausedBlocksOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.ValidatorMaximumMissedBlocksOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.ProposerPercentageOfFeesOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.ValidatorMaxEvidenceAgeInBlocksOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MissedBlocksBurnPercentageOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.DoubleSignBurnPercentageOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MessageSendFeeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MessageStakeFishermanFeeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MessageEditStakeFishermanFeeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MessageUnstakeFishermanFeeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MessagePauseFishermanFeeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MessageUnpauseFishermanFeeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MessageFishermanPauseServiceNodeFeeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MessageTestScoreFeeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MessageProveTestScoreFeeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MessageStakeAppFeeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MessageEditStakeAppFeeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MessageUnstakeAppFeeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MessagePauseAppFeeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MessageUnpauseAppFeeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MessageStakeValidatorFeeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MessageEditStakeValidatorFeeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MessageUnstakeValidatorFeeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MessagePauseValidatorFeeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MessageUnpauseValidatorFeeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MessageStakeServiceNodeFeeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MessageEditStakeServiceNodeFeeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MessageUnstakeServiceNodeFeeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MessagePauseServiceNodeFeeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MessageUnpauseServiceNodeFeeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	case typesUtil.MessageChangeParameterFeeOwner:
		return store.GetBytesParam(typesUtil.AclOwner, height)
	default:
		return nil, typesUtil.ErrUnknownParam(paramName)
	}
}

func (u *utilityContext) GetFee(msg typesUtil.Message, actorType coreTypes.ActorType) (amount *big.Int, err typesUtil.Error) {
	switch x := msg.(type) {
	case *typesUtil.MessageSend:
		return u.GetMessageSendFee()
	case *typesUtil.MessageStake:
		switch actorType {
		case coreTypes.ActorType_ACTOR_TYPE_APP:
			return u.GetMessageStakeAppFee()
		case coreTypes.ActorType_ACTOR_TYPE_FISH:
			return u.GetMessageStakeFishermanFee()
		case coreTypes.ActorType_ACTOR_TYPE_SERVICENODE:
			return u.GetMessageStakeServiceNodeFee()
		case coreTypes.ActorType_ACTOR_TYPE_VAL:
			return u.GetMessageStakeValidatorFee()
		default:
			return nil, typesUtil.ErrUnknownActorType(actorType.String())
		}
	case *typesUtil.MessageEditStake:
		switch actorType {
		case coreTypes.ActorType_ACTOR_TYPE_APP:
			return u.GetMessageEditStakeAppFee()
		case coreTypes.ActorType_ACTOR_TYPE_FISH:
			return u.GetMessageEditStakeFishermanFee()
		case coreTypes.ActorType_ACTOR_TYPE_SERVICENODE:
			return u.GetMessageEditStakeServiceNodeFee()
		case coreTypes.ActorType_ACTOR_TYPE_VAL:
			return u.GetMessageEditStakeValidatorFee()
		default:
			return nil, typesUtil.ErrUnknownActorType(actorType.String())
		}
	case *typesUtil.MessageUnstake:
		switch actorType {
		case coreTypes.ActorType_ACTOR_TYPE_APP:
			return u.GetMessageUnstakeAppFee()
		case coreTypes.ActorType_ACTOR_TYPE_FISH:
			return u.GetMessageUnstakeFishermanFee()
		case coreTypes.ActorType_ACTOR_TYPE_SERVICENODE:
			return u.GetMessageUnstakeServiceNodeFee()
		case coreTypes.ActorType_ACTOR_TYPE_VAL:
			return u.GetMessageUnstakeValidatorFee()
		default:
			return nil, typesUtil.ErrUnknownActorType(actorType.String())
		}
	case *typesUtil.MessageUnpause:
		switch actorType {
		case coreTypes.ActorType_ACTOR_TYPE_APP:
			return u.GetMessageUnpauseAppFee()
		case coreTypes.ActorType_ACTOR_TYPE_FISH:
			return u.GetMessageUnpauseFishermanFee()
		case coreTypes.ActorType_ACTOR_TYPE_SERVICENODE:
			return u.GetMessageUnpauseServiceNodeFee()
		case coreTypes.ActorType_ACTOR_TYPE_VAL:
			return u.GetMessageUnpauseValidatorFee()
		default:
			return nil, typesUtil.ErrUnknownActorType(actorType.String())
		}
	case *typesUtil.MessageChangeParameter:
		return u.GetMessageChangeParameterFee()
	default:
		return nil, typesUtil.ErrUnknownMessage(x)
	}
}

func (u *utilityContext) GetMessageChangeParameterSignerCandidates(msg *typesUtil.MessageChangeParameter) ([][]byte, typesUtil.Error) {
	owner, err := u.GetParamOwner(msg.ParameterKey)
	if err != nil {
		return nil, typesUtil.ErrGetParam(msg.ParameterKey, err)
	}
	return [][]byte{owner}, nil
}

func (u *utilityContext) getBigIntParam(paramName string) (*big.Int, typesUtil.Error) {
	store, height, err := u.getStoreAndHeight()
	if err != nil {
		return nil, typesUtil.ErrGetHeight(err)
	}
	value, err := store.GetStringParam(paramName, height)
	if err != nil {
		log.Printf("err: %v\n", err)
		return nil, typesUtil.ErrGetParam(paramName, err)
	}
	amount, err := converters.StringToBigInt(value)
	if err != nil {
		return nil, typesUtil.ErrStringToBigInt(err)
	}
	return amount, nil
}

func (u *utilityContext) getIntParam(paramName string) (int, typesUtil.Error) {
	store, height, err := u.getStoreAndHeight()
	if err != nil {
		return 0, typesUtil.ErrGetHeight(err)
	}
	value, err := store.GetIntParam(paramName, height)
	if err != nil {
		return typesUtil.ZeroInt, typesUtil.ErrGetParam(paramName, err)
	}
	return value, nil
}

func (u *utilityContext) getInt64Param(paramName string) (int64, typesUtil.Error) {
	store, height, err := u.getStoreAndHeight()
	if err != nil {
		return 0, typesUtil.ErrGetHeight(err)
	}
	value, err := store.GetIntParam(paramName, height)
	if err != nil {
		return typesUtil.ZeroInt, typesUtil.ErrGetParam(paramName, err)
	}
	return int64(value), nil
}

func (u *utilityContext) getByteArrayParam(paramName string) ([]byte, typesUtil.Error) {
	store, height, err := u.getStoreAndHeight()
	if err != nil {
		return nil, typesUtil.ErrGetHeight(err)
	}
	value, err := store.GetBytesParam(paramName, height)
	if err != nil {
		return nil, typesUtil.ErrGetParam(paramName, err)
	}
	return value, nil
}
