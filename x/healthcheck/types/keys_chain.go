package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// ChainKeyPrefix is the prefix to retrieve all Chain
	ChainKeyPrefix = "Chain/value/"
)

// ChainKey returns the store key to retrieve a Chain from the index fields
func ChainKey(
	chainId string,
) []byte {
	var key []byte

	chainIdBytes := []byte(chainId)
	key = append(key, chainIdBytes...)
	key = append(key, []byte("/")...)

	return key
}
