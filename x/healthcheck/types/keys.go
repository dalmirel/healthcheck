package types

const (
	// ModuleName defines the module name
	ModuleName = "healthcheck"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_healthcheck"

	// TODO Mirel: Moved to the commonTypes - version 1, portID: healthcheck;monitored:

	// Version defines the current version the IBC module supports
	// Version = "healthcheck-1"

	// PortID is the default port id that module binds to
	// PortID = "healthcheck"

	DefaultUpdateInterval  = 10
	DefaultTimeoutInterval = 20
)

var (
	// PortKey defines the key to store the port ID in store
	PortKey = KeyPrefix("healthcheck-port-")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

type chainActivityStatus = uint64

const (
	Inactive chainActivityStatus = iota
	Active
)
