package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAddChain = "add_chain"

var _ sdk.Msg = &MsgAddChain{}

func NewMsgAddChain(creator string, chainId string, connectionId string) *MsgAddChain {
	return &MsgAddChain{
		Creator:      creator,
		ChainId:      chainId,
		ConnectionId: connectionId,
	}
}

func (msg *MsgAddChain) Route() string {
	return RouterKey
}

func (msg *MsgAddChain) Type() string {
	return TypeMsgAddChain
}

func (msg *MsgAddChain) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddChain) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddChain) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
