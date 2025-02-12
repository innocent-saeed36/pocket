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

func FuzzFisherman(f *testing.F) {
	fuzzSingleProtocolActor(f,
		newTestGenericActor(types.FishermanActor, newTestFisherman),
		getGenericActor(types.FishermanActor, getTestFisherman),
		types.FishermanActor)
}

func TestGetSetFishermanStakeAmount(t *testing.T) {
	db := NewTestPostgresContext(t, 1)
	getTestGetSetStakeAmountTest(t, db, createAndInsertDefaultTestFisherman, db.GetFishermanStakeAmount, db.SetFishermanStakeAmount, 1)
}

func TestGetFishermanUpdatedAtHeight(t *testing.T) {
	getFishermanUpdatedFunc := func(db *persistence.PostgresContext, height int64) ([]*coreTypes.Actor, error) {
		return db.GetActorsUpdated(types.FishermanActor, height)
	}
	getAllActorsUpdatedAtHeightTest(t, createAndInsertDefaultTestFisherman, getFishermanUpdatedFunc, 1)
}

func TestInsertFishermanAndExists(t *testing.T) {
	db := NewTestPostgresContext(t, 0)

	fisherman, err := createAndInsertDefaultTestFisherman(db)
	require.NoError(t, err)

	db.Height = 1

	fisherman2, err := createAndInsertDefaultTestFisherman(db)
	require.NoError(t, err)

	addrBz, err := hex.DecodeString(fisherman.Address)
	require.NoError(t, err)
	addrBz2, err := hex.DecodeString(fisherman2.Address)
	require.NoError(t, err)

	exists, err := db.GetFishermanExists(addrBz, 0)
	require.NoError(t, err)
	require.True(t, exists, "actor that should exist at previous height does not")
	exists, err = db.GetFishermanExists(addrBz, 1)
	require.NoError(t, err)
	require.True(t, exists, "actor that should exist at current height does not")

	exists, err = db.GetFishermanExists(addrBz2, 0)
	require.NoError(t, err)
	require.False(t, exists, "actor that should not exist at previous height fishermanears to")
	exists, err = db.GetFishermanExists(addrBz2, 1)
	require.NoError(t, err)
	require.True(t, exists, "actor that should exist at current height does not")
}

func TestUpdateFisherman(t *testing.T) {
	db := NewTestPostgresContext(t, 0)

	fisherman, err := createAndInsertDefaultTestFisherman(db)
	require.NoError(t, err)

	addrBz, err := hex.DecodeString(fisherman.Address)
	require.NoError(t, err)

	fisher, err := db.GetFisherman(addrBz, 0)
	require.NoError(t, err)
	require.NotNil(t, fisher)
	require.Equal(t, DefaultChains, fisher.Chains, "default chains incorrect for current height")
	require.Equal(t, DefaultStake, fisher.StakedAmount, "default stake incorrect for current height")

	db.Height = 1

	require.NotEqual(t, DefaultStake, StakeToUpdate)   // sanity check to make sure the tests are correct
	require.NotEqual(t, DefaultChains, ChainsToUpdate) // sanity check to make sure the tests are correct
	err = db.UpdateFisherman(addrBz, fisherman.ServiceUrl, StakeToUpdate, ChainsToUpdate)
	require.NoError(t, err)

	fisher, err = db.GetFisherman(addrBz, 0)
	require.NoError(t, err)
	require.NotNil(t, fisher)
	require.Equal(t, DefaultChains, fisher.Chains, "default chains incorrect for current height")
	require.Equal(t, DefaultStake, fisher.StakedAmount, "default stake incorrect for current height")

	fisher, err = db.GetFisherman(addrBz, 1)
	require.NoError(t, err)
	require.NotNil(t, fisher)
	require.Equal(t, ChainsToUpdate, fisher.Chains, "chains not updated for current height")
	require.Equal(t, StakeToUpdate, fisher.StakedAmount, "stake not updated for current height")
}

func TestGetFishermenReadyToUnstake(t *testing.T) {
	db := NewTestPostgresContext(t, 0)

	fisherman, err := createAndInsertDefaultTestFisherman(db)
	require.NoError(t, err)

	fisherman2, err := createAndInsertDefaultTestFisherman(db)
	require.NoError(t, err)

	fisherman3, err := createAndInsertDefaultTestFisherman(db)
	require.NoError(t, err)

	addrBz, err := hex.DecodeString(fisherman.Address)
	require.NoError(t, err)
	addrBz2, err := hex.DecodeString(fisherman2.Address)
	require.NoError(t, err)
	addrBz3, err := hex.DecodeString(fisherman3.Address)
	require.NoError(t, err)

	// Unstake fisherman at height 0
	err = db.SetFishermanUnstakingHeightAndStatus(addrBz, 0, int32(coreTypes.StakeStatus_Unstaking))
	require.NoError(t, err)

	// Unstake fisherman2 and fisherman3 at height 1
	err = db.SetFishermanUnstakingHeightAndStatus(addrBz2, 1, int32(coreTypes.StakeStatus_Unstaking))
	require.NoError(t, err)
	err = db.SetFishermanUnstakingHeightAndStatus(addrBz3, 1, int32(coreTypes.StakeStatus_Unstaking))
	require.NoError(t, err)

	// Check unstaking fishermans at height 0
	unstakingFishermen, err := db.GetFishermenReadyToUnstake(0, int32(coreTypes.StakeStatus_Unstaking))
	require.NoError(t, err)
	require.Equal(t, 1, len(unstakingFishermen), "wrong number of actors ready to unstake at height 0")
	require.Equal(t, fisherman.Address, unstakingFishermen[0].GetAddress(), "unexpected fishermanlication actor returned")

	// Check unstaking fishermans at height 1
	unstakingFishermen, err = db.GetFishermenReadyToUnstake(1, int32(coreTypes.StakeStatus_Unstaking))
	require.NoError(t, err)
	require.Equal(t, 2, len(unstakingFishermen), "wrong number of actors ready to unstake at height 1")
	require.ElementsMatch(t, []string{fisherman2.Address, fisherman3.Address}, []string{unstakingFishermen[0].Address, unstakingFishermen[1].Address})
}

func TestGetFishermanStatus(t *testing.T) {
	db := NewTestPostgresContext(t, 1)

	fisherman, err := createAndInsertDefaultTestFisherman(db)
	require.NoError(t, err)

	addrBz, err := hex.DecodeString(fisherman.Address)
	require.NoError(t, err)

	// Check status before the fisherman exists
	status, err := db.GetFishermanStatus(addrBz, 0)
	require.Error(t, err)
	require.Equal(t, int32(coreTypes.StakeStatus_UnknownStatus), status, "unexpected status")

	// Check status after the fisherman exists
	status, err = db.GetFishermanStatus(addrBz, 1)
	require.NoError(t, err)
	require.Equal(t, DefaultStakeStatus, status, "unexpected status")
}

func TestGetFishermanPauseHeightIfExists(t *testing.T) {
	db := NewTestPostgresContext(t, 1)

	fisherman, err := createAndInsertDefaultTestFisherman(db)
	require.NoError(t, err)

	addrBz, err := hex.DecodeString(fisherman.Address)
	require.NoError(t, err)

	// Check pause height when fisherman does not exist
	pauseHeight, err := db.GetFishermanPauseHeightIfExists(addrBz, 0)
	require.Error(t, err)
	require.Equal(t, DefaultPauseHeight, pauseHeight, "unexpected pause height")

	// Check pause height when fisherman does not exist
	pauseHeight, err = db.GetFishermanPauseHeightIfExists(addrBz, 1)
	require.NoError(t, err)
	require.Equal(t, DefaultPauseHeight, pauseHeight, "unexpected pause height")
}

func TestSetFishermanPauseHeightAndUnstakeLater(t *testing.T) {
	db := NewTestPostgresContext(t, 0)

	fisherman, err := createAndInsertDefaultTestFisherman(db)
	require.NoError(t, err)

	pauseHeight := int64(1)
	unstakingHeight := pauseHeight + 10

	addrBz, err := hex.DecodeString(fisherman.Address)
	require.NoError(t, err)

	err = db.SetFishermanPauseHeight(addrBz, pauseHeight)
	require.NoError(t, err)

	fisher, err := db.GetFisherman(addrBz, 0)
	require.NoError(t, err)
	require.NotNil(t, fisher)
	require.Equal(t, pauseHeight, fisher.PausedHeight, "pause height not updated")

	err = db.SetFishermanStatusAndUnstakingHeightIfPausedBefore(pauseHeight+1, unstakingHeight, -1 /*unused*/)
	require.NoError(t, err)

	fisher, err = db.GetFisherman(addrBz, 0)
	require.NoError(t, err)
	require.NotNil(t, fisher)
	require.Equal(t, unstakingHeight, fisher.UnstakingHeight, "unstaking height was not set correctly")
}

func TestGetFishermanOutputAddress(t *testing.T) {
	db := NewTestPostgresContext(t, 0)

	fisherman, err := createAndInsertDefaultTestFisherman(db)
	require.NoError(t, err)

	addrBz, err := hex.DecodeString(fisherman.Address)
	require.NoError(t, err)

	output, err := db.GetFishermanOutputAddress(addrBz, 0)
	require.NoError(t, err)
	require.Equal(t, fisherman.Output, hex.EncodeToString(output), "unexpected output address")
}

func newTestFisherman() (*coreTypes.Actor, error) {
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
		Chains:          DefaultChains,
		ServiceUrl:      DefaultServiceURL,
		StakedAmount:    DefaultStake,
		PausedHeight:    DefaultPauseHeight,
		UnstakingHeight: DefaultUnstakingHeight,
		Output:          hex.EncodeToString(outputAddr),
	}, nil
}

func createAndInsertDefaultTestFisherman(db *persistence.PostgresContext) (*coreTypes.Actor, error) {
	fisherman, err := newTestFisherman()
	if err != nil {
		return nil, err
	}
	addrBz, err := hex.DecodeString(fisherman.Address)
	if err != nil {
		log.Fatalf("an error occurred converting address to bytes %s", fisherman.Address)
	}
	pubKeyBz, err := hex.DecodeString(fisherman.PublicKey)
	if err != nil {
		log.Fatalf("an error occurred converting pubKey to bytes %s", fisherman.PublicKey)
	}
	outputBz, err := hex.DecodeString(fisherman.Output)
	if err != nil {
		log.Fatalf("an error occurred converting output to bytes %s", fisherman.Output)
	}
	return fisherman, db.InsertFisherman(
		addrBz,
		pubKeyBz,
		outputBz,
		false,
		DefaultStakeStatus,
		DefaultServiceURL,
		DefaultStake,
		DefaultChains,
		DefaultPauseHeight,
		DefaultUnstakingHeight)
}

func getTestFisherman(db *persistence.PostgresContext, address []byte) (*coreTypes.Actor, error) {
	return db.GetFisherman(address, db.Height)
}
