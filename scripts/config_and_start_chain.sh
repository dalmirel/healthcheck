#!/bin/bash
set -eu



HOME_DIR="$1"
BINARY="$2"
CHAIN_ID="$3"
RPC_LADDR_PORT="$4"
GRPC_ADDR_PORT="$5"
ADDR_PORT="$6"
P2P_LADDR_PORT="$7"

CHAIN_ID_FLAG="--chain-id $CHAIN_ID"
KEYRING_TEST_FLAG="--keyring-backend test"
NODE_IP="localhost"
SEED_PHRASE_1="direct island hammer gentle cook hollow obvious promote bracket gravity file alcohol rule frost base hint smart foot soup time purity margin lend pencil"
SEED_PHRASE_2="long physical balcony pool increase detail fire light veteran skull blade physical skirt neglect width matrix dish snake soap amount bottom wash bean life"

clean_up(){
    #pkill $BINARY &> /dev/null || true
    rm -rf $HOME_DIR
}

# Clean start
clean_up

$BINARY init $CHAIN_ID $CHAIN_ID_FLAG --home $HOME_DIR

sleep 1
echo $SEED_PHRASE_1 | $BINARY keys add wallet1 $KEYRING_TEST_FLAG --recover --home $HOME_DIR --output json > $HOME_DIR/validator.json 2>&1
echo $SEED_PHRASE_2 | $BINARY keys add wallet2 $KEYRING_TEST_FLAG --recover --home $HOME_DIR --output json > $HOME_DIR/relayer.json 2>&1
$BINARY  add-genesis-account wallet1 1000000000000stake,1000000000000token $KEYRING_TEST_FLAG --home $HOME_DIR
$BINARY  add-genesis-account wallet2 1000000000000stake,1000000000000token $KEYRING_TEST_FLAG --home $HOME_DIR  
$BINARY  gentx wallet1 1000000000stake $KEYRING_TEST_FLAG $CHAIN_ID_FLAG --home $HOME_DIR
$BINARY  collect-gentxs --home $HOME_DIR

$BINARY start \
	--log_level trace \
    --minimum-gas-prices 0.00001stake \
	--home $HOME_DIR \
	--rpc.laddr tcp://$NODE_IP:$RPC_LADDR_PORT \
	--grpc.address $NODE_IP:$GRPC_ADDR_PORT \
	--address tcp://$NODE_IP:$ADDR_PORT \
	--p2p.laddr tcp://$NODE_IP:$P2P_LADDR_PORT \
	--grpc-web.enable=false &> "$HOME_DIR"/logs &

sleep 10
