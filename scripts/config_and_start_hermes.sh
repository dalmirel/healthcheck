#!/bin/bash
set -eu
CHAIN1_GRPC_ADDR="localhost:9091"
CHAIN1_CHAIN_ID="healthcheck"
CHAIN1_RPC_LADDR="localhost:26658"

CHAIN2_GRPC_ADDR="localhost:9101"
CHAIN2_CHAIN_ID="monitored"
CHAIN2_RPC_LADDR="localhost:26668"

MNEMONIC=./mnemonic.txt

# Setup Hermes in packet relayer mode
killall hermes 2> /dev/null || true

#write over existing hermes config file:
tee ~/.hermes/config.toml<<EOF
[global]
log_level = "trace"

[mode]

[mode.clients]
enabled = true
refresh = true
misbehaviour = true

[mode.connections]
enabled = true

[mode.channels]
enabled = true

[mode.packets]
enabled = true

[[chains]]
account_prefix = "cosmos"
clock_drift = "5s"
gas_multiplier = 10
grpc_addr = "tcp://${CHAIN1_GRPC_ADDR}"
id = "$CHAIN1_CHAIN_ID"
key_name = "relayer"
max_gas = 2000000
rpc_addr = "http://${CHAIN1_RPC_LADDR}"
rpc_timeout = "10s"
store_prefix = "ibc"
trusting_period = "599s"
event_source = { mode = 'push', url = 'ws://${CHAIN1_RPC_LADDR}/websocket', batch_delay = '500ms' }

[chains.gas_price]
       denom = "stake"
       price = 1

[chains.trust_threshold]
       denominator = "3"
       numerator = "1"

[[chains]]
account_prefix = "cosmos"
clock_drift = "5s"
gas_multiplier = 10
grpc_addr = "tcp://${CHAIN2_GRPC_ADDR}"
id = "$CHAIN2_CHAIN_ID"
key_name = "relayer"
max_gas = 2000000
rpc_addr = "http://${CHAIN2_RPC_LADDR}"
rpc_timeout = "10s"
store_prefix = "ibc"
trusting_period = "599s"
event_source = { mode = 'push', url = 'ws://${CHAIN2_RPC_LADDR}/websocket', batch_delay = '500ms' }

[chains.gas_price]
       denom = "stake"
       price = 1

[chains.trust_threshold]
       denominator = "3"
       numerator = "1"
EOF

# Delete all previous keys in relayer
hermes keys delete --chain $CHAIN1_CHAIN_ID --all
hermes keys delete --chain $CHAIN2_CHAIN_ID --all

# Restore keys to hermes relayer
hermes keys add --chain $CHAIN1_CHAIN_ID --mnemonic-file $MNEMONIC 
hermes keys add --chain $CHAIN2_CHAIN_ID --mnemonic-file $MNEMONIC 

sleep 5
hermes create connection --a-chain $CHAIN1_CHAIN_ID  --b-chain $CHAIN2_CHAIN_ID
sleep 5

hermes start &> ~/.hermes/logs &