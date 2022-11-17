package keeper_test

import (
	"math/rand"
	"testing"

	"github.com/babylonchain/babylon/testutil/datagen"
	zctypes "github.com/babylonchain/babylon/x/zoneconcierge/types"
	"github.com/stretchr/testify/require"
)

func FuzzFinalizedChainInfo(f *testing.F) {
	datagen.AddRandomSeedsToFuzzer(f, 10)

	f.Fuzz(func(t *testing.T, seed int64) {
		rand.Seed(seed)

		_, babylonChain, czChain, zcKeeper := SetupTest(t)

		ctx := babylonChain.GetContext()
		hooks := zcKeeper.Hooks()

		// invoke the hook a random number of times to simulate a random number of blocks
		numHeaders := datagen.RandomInt(100) + 1
		numForkHeaders := datagen.RandomInt(10) + 1
		SimulateHeadersAndForksViaHook(ctx, hooks, czChain.ChainID, numHeaders, numForkHeaders)

		// simulate the scenario that a random epoch has ended and finalised
		epochNum := datagen.RandomInt(10)
		hooks.AfterEpochEnds(ctx, epochNum)
		hooks.AfterRawCheckpointFinalized(ctx, epochNum)

		// check if the chain info of this epoch is recorded or not
		resp, err := zcKeeper.FinalizedChainInfo(ctx, &zctypes.QueryFinalizedChainInfoRequest{ChainId: czChain.ChainID})
		require.NoError(t, err)
		chainInfo := resp.FinalizedChainInfo
		require.Equal(t, numHeaders-1, chainInfo.LatestHeader.Height)
		require.Equal(t, numForkHeaders, uint64(len(chainInfo.LatestForks.Headers)))
	})
}