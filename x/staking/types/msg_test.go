package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	coinPos  = sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)
	coinZero = sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)
)

func TestMsgDecode(t *testing.T) {
	registry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(registry)
	RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)

	// firstly we start testing the pubkey serialization

	pk1bz, err := codec.MarshalAny(cdc, pk1)
	require.NoError(t, err)
	var pkUnmarshaled cryptotypes.PubKey
	err = codec.UnmarshalAny(cdc, &pkUnmarshaled, pk1bz)
	require.NoError(t, err)
	require.True(t, pk1.Equals(pkUnmarshaled.(*ed25519.PubKey)))

	// now let's try to serialize the whole message

	commission1 := NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
	msg, err := NewMsgCreateValidator(valAddr1, pk1, coinPos, Description{}, commission1, sdk.OneInt())
	require.NoError(t, err)
	msgSerialized, err := codec.MarshalAny(cdc, msg)
	require.NoError(t, err)

	var msgUnmarshaled sdk.Msg
	err = codec.UnmarshalAny(cdc, &msgUnmarshaled, msgSerialized)
	require.NoError(t, err)
	msg2, ok := msgUnmarshaled.(*MsgCreateValidator)
	require.True(t, ok)
	require.True(t, msg.Value.IsEqual(msg2.Value))
	require.True(t, msg.Pubkey.Equal(msg2.Pubkey))
}

// test ValidateBasic for MsgCreateValidator
func TestMsgCreateValidator(t *testing.T) {
	commission1 := NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
	commission2 := NewCommissionRates(sdk.NewDec(5), sdk.NewDec(5), sdk.NewDec(5))

	tests := []struct {
		name, moniker, identity, website, securityContact, details string
		CommissionRates                                            CommissionRates
		minSelfDelegation                                          sdk.Int
		validatorAddr                                              sdk.ValAddress
		pubkey                                                     cryptotypes.PubKey
		bond                                                       sdk.Coin
		expectPass                                                 bool
	}{
		{"basic good", "a", "b", "c", "d", "e", commission1, sdk.OneInt(), valAddr1, pk1, coinPos, true},
		{"partial description", "", "", "c", "", "", commission1, sdk.OneInt(), valAddr1, pk1, coinPos, true},
		{"empty description", "", "", "", "", "", commission2, sdk.OneInt(), valAddr1, pk1, coinPos, false},
		{"empty address", "a", "b", "c", "d", "e", commission2, sdk.OneInt(), emptyAddr, pk1, coinPos, false},
		{"empty pubkey", "a", "b", "c", "d", "e", commission1, sdk.OneInt(), valAddr1, emptyPubkey, coinPos, false},
		{"empty bond", "a", "b", "c", "d", "e", commission2, sdk.OneInt(), valAddr1, pk1, coinZero, false},
		{"nil bond", "a", "b", "c", "d", "e", commission2, sdk.OneInt(), valAddr1, pk1, sdk.Coin{}, false},
		{"zero min self delegation", "a", "b", "c", "d", "e", commission1, sdk.ZeroInt(), valAddr1, pk1, coinPos, false},
		{"negative min self delegation", "a", "b", "c", "d", "e", commission1, sdk.NewInt(-1), valAddr1, pk1, coinPos, false},
		{"delegation less than min self delegation", "a", "b", "c", "d", "e", commission1, coinPos.Amount.Add(sdk.OneInt()), valAddr1, pk1, coinPos, false},
	}

	for _, tc := range tests {
		t.Logf("Test: %s, pk=%t", tc.name, tc.pubkey)
		description := NewDescription(tc.moniker, tc.identity, tc.website, tc.securityContact, tc.details)
		msg, err := NewMsgCreateValidator(tc.validatorAddr, tc.pubkey, tc.bond, description, tc.CommissionRates, tc.minSelfDelegation)
		require.NoError(t, err)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

// test ValidateBasic for MsgEditValidator
func TestMsgEditValidator(t *testing.T) {
	tests := []struct {
		name, moniker, identity, website, securityContact, details string
		validatorAddr                                              sdk.ValAddress
		expectPass                                                 bool
		minSelfDelegation                                          sdk.Int
	}{
		{"basic good", "a", "b", "c", "d", "e", valAddr1, true, sdk.OneInt()},
		{"partial description", "", "", "c", "", "", valAddr1, true, sdk.OneInt()},
		{"empty description", "", "", "", "", "", valAddr1, false, sdk.OneInt()},
		{"empty address", "a", "b", "c", "d", "e", emptyAddr, false, sdk.OneInt()},
		{"nil int", "a", "b", "c", "d", "e", emptyAddr, false, sdk.Int{}},
	}

	for _, tc := range tests {
		description := NewDescription(tc.moniker, tc.identity, tc.website, tc.securityContact, tc.details)
		newRate := sdk.ZeroDec()

		msg := NewMsgEditValidator(tc.validatorAddr, description, &newRate, &tc.minSelfDelegation)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

// test ValidateBasic for MsgDelegate
func TestMsgDelegate(t *testing.T) {
	tests := []struct {
		name          string
		delegatorAddr sdk.AccAddress
		validatorAddr sdk.ValAddress
		bond          sdk.Coin
		expectPass    bool
	}{
		{"basic good", sdk.AccAddress(valAddr1), valAddr2, coinPos, true},
		{"self bond", sdk.AccAddress(valAddr1), valAddr1, coinPos, true},
		{"empty delegator", sdk.AccAddress(emptyAddr), valAddr1, coinPos, false},
		{"empty validator", sdk.AccAddress(valAddr1), emptyAddr, coinPos, false},
		{"empty bond", sdk.AccAddress(valAddr1), valAddr2, coinZero, false},
		{"nil bold", sdk.AccAddress(valAddr1), valAddr2, sdk.Coin{}, false},
	}

	for _, tc := range tests {
		msg := NewMsgDelegate(tc.delegatorAddr, tc.validatorAddr, tc.bond)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

// test ValidateBasic for MsgUnbond
func TestMsgBeginRedelegate(t *testing.T) {
	tests := []struct {
		name             string
		delegatorAddr    sdk.AccAddress
		validatorSrcAddr sdk.ValAddress
		validatorDstAddr sdk.ValAddress
		amount           sdk.Coin
		expectPass       bool
	}{
		{"regular", sdk.AccAddress(valAddr1), valAddr2, valAddr3, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), true},
		{"zero amount", sdk.AccAddress(valAddr1), valAddr2, valAddr3, sdk.NewInt64Coin(sdk.DefaultBondDenom, 0), false},
		{"nil amount", sdk.AccAddress(valAddr1), valAddr2, valAddr3, sdk.Coin{}, false},
		{"empty delegator", sdk.AccAddress(emptyAddr), valAddr1, valAddr3, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), false},
		{"empty source validator", sdk.AccAddress(valAddr1), emptyAddr, valAddr3, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), false},
		{"empty destination validator", sdk.AccAddress(valAddr1), valAddr2, emptyAddr, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), false},
	}

	for _, tc := range tests {
		msg := NewMsgBeginRedelegate(tc.delegatorAddr, tc.validatorSrcAddr, tc.validatorDstAddr, tc.amount)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

// test ValidateBasic for MsgUnbond
func TestMsgUndelegate(t *testing.T) {
	tests := []struct {
		name          string
		delegatorAddr sdk.AccAddress
		validatorAddr sdk.ValAddress
		amount        sdk.Coin
		expectPass    bool
	}{
		{"regular", sdk.AccAddress(valAddr1), valAddr2, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), true},
		{"zero amount", sdk.AccAddress(valAddr1), valAddr2, sdk.NewInt64Coin(sdk.DefaultBondDenom, 0), false},
		{"nil amount", sdk.AccAddress(valAddr1), valAddr2, sdk.Coin{}, false},
		{"empty delegator", sdk.AccAddress(emptyAddr), valAddr1, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), false},
		{"empty validator", sdk.AccAddress(valAddr1), emptyAddr, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), false},
	}

	for _, tc := range tests {
		msg := NewMsgUndelegate(tc.delegatorAddr, tc.validatorAddr, tc.amount)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}
