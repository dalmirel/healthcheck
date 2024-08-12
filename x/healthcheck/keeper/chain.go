package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"healthcheck/x/healthcheck/types"
)

// GetChainCount get the total number of chain
func (k Keeper) GetChainCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.ChainCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetChainCount set the total number of chain
func (k Keeper) SetChainCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.ChainCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendChain appends a chain in the store with a new id and update the count
func (k Keeper) AppendChain(
	ctx sdk.Context,
	chain types.Chain,
) uint64 {
	// Create the chain
	count := k.GetChainCount(ctx)

	// Set the ID of the appended value
	chain.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ChainKey))
	appendedValue := k.cdc.MustMarshal(&chain)
	store.Set(GetChainIDBytes(chain.Id), appendedValue)

	// Update chain count
	k.SetChainCount(ctx, count+1)

	return count
}

// SetChain set a specific chain in the store
func (k Keeper) SetChain(ctx sdk.Context, chain types.Chain) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ChainKey))
	b := k.cdc.MustMarshal(&chain)
	store.Set(GetChainIDBytes(chain.Id), b)
}

// GetChain returns a chain from its id
func (k Keeper) GetChain(ctx sdk.Context, id uint64) (val types.Chain, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ChainKey))
	b := store.Get(GetChainIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveChain removes a chain from the store
func (k Keeper) RemoveChain(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ChainKey))
	store.Delete(GetChainIDBytes(id))
}

// GetAllChain returns all chain
func (k Keeper) GetAllChain(ctx sdk.Context) (list []types.Chain) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ChainKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Chain
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetChainIDBytes returns the byte representation of the ID
func GetChainIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetChainIDFromBytes returns ID in uint64 format from a byte array
func GetChainIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
