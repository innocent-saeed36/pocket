package test

import (
	"encoding/hex"
	"log"
	"testing"

	"github.com/pokt-network/pocket/persistence/types"

	"github.com/pokt-network/pocket/persistence"
	coreTypes "github.com/pokt-network/pocket/shared/core/types"
	"github.com/pokt-network/pocket/shared/crypto"
	"github.com/stretchr/testify/require"
)

func FuzzValidator(f *testing.F) {
	fuzzSingleProtocolActor(f,
		newTestGenericActor(types.ValidatorActor, newTestValidator),
		getGenericActor(types.ValidatorActor, getTestValidator),
		types.ValidatorActor)
}

func TestGetSetValidatorStakeAmount(t *testing.T) {
	db := NewTestPostgresContext(t, 1)
	getTestGetSetStakeAmountTest(t, db, createAndInsertDefaultTestValidator, db.GetValidatorStakeAmount, db.SetValidatorStakeAmount, 1)
}

func TestGetValidatorUpdatedAtHeight(t *testing.T) {
	getValidatorsUpdatedFunc := func(db *persistence.PostgresContext, height int64) ([]*coreTypes.Actor, error) {
		return db.GetActorsUpdated(types.ValidatorActor, height)
	}
	getAllActorsUpdatedAtHeightTest(t, createAndInsertDefaultTestValidator, getValidatorsUpdatedFunc, 5)
}

func TestInsertValidatorAndExists(t *testing.T) {
	db := NewTestPostgresContext(t, 0)

	validator, err := createAndInsertDefaultTestValidator(db)
	require.NoError(t, err)

	db.Height = 1

	validator2, err := createAndInsertDefaultTestValidator(db)
	require.NoError(t, err)

	addrBz, err := hex.DecodeString(validator.Address)
	require.NoError(t, err)

	addrBz2, err := hex.DecodeString(validator2.Address)
	require.NoError(t, err)

	exists, err := db.GetValidatorExists(addrBz, 0)
	require.NoError(t, err)
	require.True(t, exists, "actor that should exist at previous height does not")
	exists, err = db.GetValidatorExists(addrBz, 1)
	require.NoError(t, err)
	require.True(t, exists, "actor that should exist at current height does not")

	exists, err = db.GetValidatorExists(addrBz2, 0)
	require.NoError(t, err)
	require.False(t, exists, "actor that should not exist at previous height does")
	exists, err = db.GetValidatorExists(addrBz2, 1)
	require.NoError(t, err)
	require.True(t, exists, "actor that should exist at current height does not")
}

func TestUpdateValidator(t *testing.T) {
	db := NewTestPostgresContext(t, 0)

	validator, err := createAndInsertDefaultTestValidator(db)
	require.NoError(t, err)

	addrBz, err := hex.DecodeString(validator.Address)
	require.NoError(t, err)

	val, err := db.GetValidator(addrBz, 0)
	require.NoError(t, err)
	require.NotNil(t, val)
	require.Equal(t, DefaultStake, val.StakedAmount, "default stake incorrect for current height")

	db.Height = 1

	require.NotEqual(t, DefaultStake, StakeToUpdate) // sanity check to make sure the tests are correct
	err = db.UpdateValidator(addrBz, validator.ServiceUrl, StakeToUpdate)
	require.NoError(t, err)

	val, err = db.GetValidator(addrBz, 0)
	require.NoError(t, err)
	require.NotNil(t, val)
	require.Equal(t, DefaultStake, val.StakedAmount, "default stake incorrect for current height")

	val, err = db.GetValidator(addrBz, 1)
	require.NoError(t, err)
	require.NotNil(t, val)
	require.Equal(t, StakeToUpdate, val.StakedAmount, "stake not updated for current height")
}

func TestGetValidatorsReadyToUnstake(t *testing.T) {
	db := NewTestPostgresContext(t, 0)

	validator, err := createAndInsertDefaultTestValidator(db)
	require.NoError(t, err)

	validator2, err := createAndInsertDefaultTestValidator(db)
	require.NoError(t, err)

	validator3, err := createAndInsertDefaultTestValidator(db)
	require.NoError(t, err)

	addrBz, err := hex.DecodeString(validator.Address)
	require.NoError(t, err)

	addrBz2, err := hex.DecodeString(validator2.Address)
	require.NoError(t, err)

	addrBz3, err := hex.DecodeString(validator3.Address)
	require.NoError(t, err)

	// Unstake validator at height 0
	err = db.SetValidatorUnstakingHeightAndStatus(addrBz, 0, int32(coreTypes.StakeStatus_Unstaking))
	require.NoError(t, err)

	// Unstake validator2 and validator3 at height 1
	err = db.SetValidatorUnstakingHeightAndStatus(addrBz2, 1, int32(coreTypes.StakeStatus_Unstaking))
	require.NoError(t, err)
	err = db.SetValidatorUnstakingHeightAndStatus(addrBz3, 1, int32(coreTypes.StakeStatus_Unstaking))
	require.NoError(t, err)

	// Check unstaking validators at height 0
	unstakingValidators, err := db.GetValidatorsReadyToUnstake(0, int32(coreTypes.StakeStatus_Unstaking))
	require.NoError(t, err)
	require.Equal(t, 1, len(unstakingValidators), "wrong number of actors ready to unstake at height 0")

	require.Equal(t, validator.Address, unstakingValidators[0].GetAddress(), "unexpected unstaking validator returned")

	// Check unstaking validators at height 1
	unstakingValidators, err = db.GetValidatorsReadyToUnstake(1, int32(coreTypes.StakeStatus_Unstaking))
	require.NoError(t, err)
	require.Equal(t, 2, len(unstakingValidators), "wrong number of actors ready to unstake at height 1")
	require.ElementsMatch(t, []string{validator2.Address, validator3.Address}, []string{unstakingValidators[0].Address, unstakingValidators[1].Address})
}

func TestGetValidatorStatus(t *testing.T) {
	db := NewTestPostgresContext(t, 1)

	validator, err := createAndInsertDefaultTestValidator(db)
	require.NoError(t, err)

	addrBz, err := hex.DecodeString(validator.Address)
	require.NoError(t, err)

	// Check status before the validator exists
	status, err := db.GetValidatorStatus(addrBz, 0)
	require.Error(t, err)
	require.Equal(t, int32(coreTypes.StakeStatus_UnknownStatus), status, "unexpected status")

	// Check status after the validator exists
	status, err = db.GetValidatorStatus(addrBz, 1)
	require.NoError(t, err)
	require.Equal(t, DefaultStakeStatus, status, "unexpected status")
}

func TestGetValidatorPauseHeightIfExists(t *testing.T) {
	db := NewTestPostgresContext(t, 1)

	validator, err := createAndInsertDefaultTestValidator(db)
	require.NoError(t, err)

	addrBz, err := hex.DecodeString(validator.Address)
	require.NoError(t, err)

	// Check pause height when validator does not exist
	pauseHeight, err := db.GetValidatorPauseHeightIfExists(addrBz, 0)
	require.Error(t, err)
	require.Equal(t, DefaultPauseHeight, pauseHeight, "unexpected pause height")

	// Check pause height when validator does not exist
	pauseHeight, err = db.GetValidatorPauseHeightIfExists(addrBz, 1)
	require.NoError(t, err)
	require.Equal(t, DefaultPauseHeight, pauseHeight, "unexpected pause height")
}

func TestSetValidatorPauseHeightAndUnstakeLater(t *testing.T) {
	db := NewTestPostgresContext(t, 0)

	validator, err := createAndInsertDefaultTestValidator(db)
	require.NoError(t, err)

	pauseHeight := int64(1)
	unstakingHeight := pauseHeight + 10

	addrBz, err := hex.DecodeString(validator.Address)
	require.NoError(t, err)

	err = db.SetValidatorPauseHeight(addrBz, pauseHeight)
	require.NoError(t, err)

	val, err := db.GetValidator(addrBz, 0)
	require.NoError(t, err)
	require.NotNil(t, val)
	require.Equal(t, pauseHeight, val.PausedHeight, "pause height not updated")

	err = db.SetValidatorsStatusAndUnstakingHeightIfPausedBefore(pauseHeight+1, unstakingHeight, -1 /*unused*/)
	require.NoError(t, err)

	val, err = db.GetValidator(addrBz, 0)
	require.NoError(t, err)
	require.NotNil(t, val)
	require.Equal(t, unstakingHeight, val.UnstakingHeight, "unstaking height was not set correctly")
}

func TestGetValidatorOutputAddress(t *testing.T) {
	db := NewTestPostgresContext(t, 0)

	validator, err := createAndInsertDefaultTestValidator(db)
	require.NoError(t, err)

	addrBz, err := hex.DecodeString(validator.Address)
	require.NoError(t, err)

	output, err := db.GetValidatorOutputAddress(addrBz, 0)
	require.NoError(t, err)
	require.Equal(t, validator.Output, hex.EncodeToString(output), "unexpected output address")
}

func newTestValidator() (*coreTypes.Actor, error) {
	operatorKey, err := crypto.GeneratePublicKey()
	if err != nil {
		return nil, err
	}

	outputAddr, err := crypto.GenerateAddress()
	if err != nil {
		return nil, err
	}

	return &coreTypes.Actor{
		Address:         hex.EncodeToString(operatorKey.Address()),
		PublicKey:       hex.EncodeToString(operatorKey.Bytes()),
		ServiceUrl:      DefaultServiceURL,
		StakedAmount:    DefaultStake,
		PausedHeight:    DefaultPauseHeight,
		UnstakingHeight: DefaultUnstakingHeight,
		Output:          hex.EncodeToString(outputAddr),
	}, nil
}

func createAndInsertDefaultTestValidator(db *persistence.PostgresContext) (*coreTypes.Actor, error) {
	validator, err := newTestValidator()
	if err != nil {
		return nil, err
	}
	addrBz, err := hex.DecodeString(validator.Address)
	if err != nil {
		log.Fatalf("an error occurred converting address to bytes %s", validator.Address)
	}
	pubKeyBz, err := hex.DecodeString(validator.PublicKey)
	if err != nil {
		log.Fatalf("an error occurred converting pubKey to bytes %s", validator.PublicKey)
	}
	outputBz, err := hex.DecodeString(validator.Output)
	if err != nil {
		log.Fatalf("an error occurred converting output to bytes %s", validator.Output)
	}
	return validator, db.InsertValidator(
		addrBz,
		pubKeyBz,
		outputBz,
		false,
		DefaultStakeStatus,
		DefaultServiceURL,
		DefaultStake,
		DefaultPauseHeight,
		DefaultUnstakingHeight)
}

func getTestValidator(db *persistence.PostgresContext, address []byte) (*coreTypes.Actor, error) {
	return db.GetValidator(address, db.Height)
}
