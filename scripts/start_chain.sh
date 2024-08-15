#!/bin/bash
set -eu

# this can be a parameter
HOME="/home/mirel"

CHAIN_HOME="$HOME"/.healthcheck
BINARY="healthcheckd"
CHAIN_ID="healthcheck"
echo $CHAIN_HOME

# clean home folder - delete genesis
rm -rf $CHAIN_HOME

$BINARY init "healthcheck" --chain-id $CHAIN_ID --home $CHAIN_HOME

$BINARY keys add wallet1 --keyring-backend test
$BINARY keys add wallet2 --keyring-backend test
$BINARY add-genesis-account wallet1 1000000000000stake --keyring-backend test
$BINARY add-genesis-account wallet2 1000000000000stake --keyring-backend test
$BINARY gentx wallet1 1000000000stake --keyring-backend test --chain-id $CHAIN_ID --home $CHAIN_HOME
$BINARY collect-gentxs --home $CHAIN_HOME

$BINARY start --minimum-gas-prices 0.0001stake --home $CHAIN_HOME &> "$CHAIN_HOME"/logs &


# stop the node
# pkill healthcheckd