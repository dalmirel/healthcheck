package types

const (
	// ModuleName defines the module name
	ModuleName = "monitored"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_monitored"

	// TODO Mirel: Moved to the commonTypes - version 1, portID: healthcheck;monitored:

	// Version defines the current version the IBC module supports
	// Version = "monitored-1"

	// PortID is the default port id that module binds to
	// PortID = "monitored"
)

var (
	// PortKey defines the key to store the port ID in store
	PortKey = KeyPrefix("monitored-port-")

	// key to store the channel ID
	// that is used to send healthcheck updates to registry chain
	RegistryChainChannelIDKey = KeyPrefix("RegistryChainChannelID")

	// key to store the last block height
	// for which the healthcheck status was sent to registry chain
	LastHealthcheckUpdateHeightKey = KeyPrefix("LastHealthcheckUpdateHeight")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
