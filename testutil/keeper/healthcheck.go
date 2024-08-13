package keeper

import (
	"testing"

	"healthcheck/x/healthcheck/keeper"
	"healthcheck/x/healthcheck/types"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	tendermintClient "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	"github.com/stretchr/testify/require"
)

// healthcheckChannelKeeper is a stub of cosmosibckeeper.ChannelKeeper.
type healthcheckChannelKeeper struct{}

func (healthcheckChannelKeeper) GetChannel(ctx sdk.Context, portID, channelID string) (channeltypes.Channel, bool) {
	return channeltypes.Channel{}, false
}

func (healthcheckChannelKeeper) GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool) {
	return 0, false
}

func (healthcheckChannelKeeper) SendPacket(
	ctx sdk.Context,
	channelCap *capabilitytypes.Capability,
	sourcePort string,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	data []byte,
) (uint64, error) {
	return 0, nil
}

func (healthcheckChannelKeeper) ChanCloseInit(ctx sdk.Context, portID, channelID string, chanCap *capabilitytypes.Capability) error {
	return nil
}

func (healthcheckChannelKeeper) GetChannelClientState(ctx sdk.Context, portID, channelID string) (string, ibcexported.ClientState, error) {
	return "", nil, nil
}

// healthcheckportKeeper is a stub of cosmosibckeeper.PortKeeper
type healthcheckPortKeeper struct{}

func (healthcheckPortKeeper) BindPort(ctx sdk.Context, portID string) *capabilitytypes.Capability {
	return &capabilitytypes.Capability{}
}

type healthcheckClientKeeper struct{}

func (healthcheckClientKeeper) GetClientState(ctx sdk.Context, clientID string) (ibcexported.ClientState, bool) {
	return &tendermintClient.ClientState{}, true
}

type healthcheckConnectionKeeper struct{}

func (healthcheckConnectionKeeper) GetConnection(ctx sdk.Context, connectionID string) (connectiontypes.ConnectionEnd, bool) {
	return connectiontypes.ConnectionEnd{}, true
}

func HealthcheckKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	logger := log.NewNopLogger()

	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	appCodec := codec.NewProtoCodec(registry)
	capabilityKeeper := capabilitykeeper.NewKeeper(appCodec, storeKey, memStoreKey)

	paramsSubspace := typesparams.NewSubspace(appCodec,
		types.Amino,
		storeKey,
		memStoreKey,
		"HealthcheckParams",
	)
	k := keeper.NewKeeper(
		appCodec,
		storeKey,
		memStoreKey,
		paramsSubspace,
		healthcheckChannelKeeper{},
		healthcheckPortKeeper{},
		healthcheckConnectionKeeper{},
		healthcheckClientKeeper{},
		capabilityKeeper.ScopeToModule("HealthcheckScopedKeeper"),
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, logger)

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}
