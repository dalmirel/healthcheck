#!/bin/bash
set -eu

# this can be a parameter
HOME="/home/mirel"

# helper functions 
set_chain1_vars(){
    CHAIN_HOME="$HOME"/.healthcheck
    BINARY="healthcheckd"
    CHAIN_ID="healthcheck"
    RPC_LADDR_PORT="26658"
    GRPC_ADDR_PORT="9091"
    ADDR_PORT="26655"
    P2P_LADDR_PORT="26656"
    TX_FLAGS="--gas-adjustment 15 --gas-prices 0.00001stake --gas 300000 --chain-id $CHAIN_ID --node tcp://localhost:$RPC_LADDR_PORT --from wallet1 --home $CHAIN_HOME --keyring-backend test -y"
}

set_chain2_vars(){
    CHAIN_HOME="$HOME"/.monitored
    BINARY="healthcheckd"
    CHAIN_ID="monitored"
    RPC_LADDR_PORT="26668"
    GRPC_ADDR_PORT="9101"
    ADDR_PORT="26665"
    P2P_LADDR_PORT="26666"
    TX_FLAGS="--gas 10000000 --gas-adjustment 100 --gas-prices 0.00001stake --chain-id $CHAIN_ID --node tcp://localhost:$RPC_LADDR_PORT --from wallet1 --home $CHAIN_HOME --keyring-backend test -y"
}

start_chain(){
    if ! bash config_and_start_chain.sh $CHAIN_HOME $BINARY $CHAIN_ID $RPC_LADDR_PORT $GRPC_ADDR_PORT $ADDR_PORT $P2P_LADDR_PORT; 
    then
	    echo "Error starting $CHAIN_ID chain."
	    exit 1
    fi
}

# start two chains (healthcheck and monitored) from the same binary
set_chain1_vars
start_chain
set_chain2_vars
start_chain