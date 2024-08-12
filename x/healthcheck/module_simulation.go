package healthcheck

import (
	"math/rand"

	"healthcheck/testutil/sample"
	healthchecksimulation "healthcheck/x/healthcheck/simulation"
	"healthcheck/x/healthcheck/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	commonTypes "healthcheck/x/common"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = healthchecksimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
	_ = rand.Rand{}
)

const (
	opWeightMsgAddChain = "op_weight_msg_add_chain"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddChain int = 100

	opWeightMsgDeleteChain = "op_weight_msg_delete_chain"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteChain int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	healthcheckGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: commonTypes.HealthcheckPortID,
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&healthcheckGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgAddChain int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgAddChain, &weightMsgAddChain, nil,
		func(_ *rand.Rand) {
			weightMsgAddChain = defaultWeightMsgAddChain
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddChain,
		healthchecksimulation.SimulateMsgAddChain(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteChain int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDeleteChain, &weightMsgDeleteChain, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteChain = defaultWeightMsgDeleteChain
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteChain,
		healthchecksimulation.SimulateMsgDeleteChain(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgAddChain,
			defaultWeightMsgAddChain,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				healthchecksimulation.SimulateMsgAddChain(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgDeleteChain,
			defaultWeightMsgDeleteChain,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				healthchecksimulation.SimulateMsgDeleteChain(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
