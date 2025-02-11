package rpcv02

import (
	"context"
	"regexp"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/types"
	"github.com/NethermindEth/starknet.go/utils"
)

// TestCall tests Call
func TestCall(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		FunctionCall          FunctionCall
		BlockID               BlockID
		ExpectedPatternResult string
	}
	testSet := map[string][]testSetType{
		"devnet": {
			{
				FunctionCall: FunctionCall{
					// ContractAddress of predeployed devnet Feetoken
					ContractAddress:    utils.TestHexToFelt(t, "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
					EntryPointSelector: types.GetSelectorFromNameFelt("name"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:               WithBlockTag("latest"),
				ExpectedPatternResult: "^0x[0-9a-f]+$",
			},
		},
		"mock": {
			{
				FunctionCall: FunctionCall{
					ContractAddress:    utils.TestHexToFelt(t, "0xdeadbeef"),
					EntryPointSelector: types.GetSelectorFromNameFelt("decimals"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:               WithBlockTag("latest"),
				ExpectedPatternResult: "^0x12$",
			},
		},
		"testnet": {
			{
				FunctionCall: FunctionCall{
					ContractAddress:    utils.TestHexToFelt(t, "0x029260ce936efafa6d0042bc59757a653e3f992b97960c1c4f8ccd63b7a90136"),
					EntryPointSelector: utils.TestHexToFelt(t, "0x004c4fb1ab068f6039d5780c68dd0fa2f8742cceb3426d19667778ca7f3518a9"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:               WithBlockTag("latest"),
				ExpectedPatternResult: "^0x12$",
			},
			{
				FunctionCall: FunctionCall{
					ContractAddress:    utils.TestHexToFelt(t, TestNetETHAddress),
					EntryPointSelector: types.GetSelectorFromNameFelt("balanceOf"),
					Calldata:           []*felt.Felt{utils.TestHexToFelt(t, "0x0207aCC15dc241e7d167E67e30E769719A727d3E0fa47f9E187707289885Dfde")},
				},
				BlockID:               WithBlockTag("latest"),
				ExpectedPatternResult: "^0x[0-9a-f]+$",
			},
			{
				FunctionCall: FunctionCall{
					ContractAddress:    utils.TestHexToFelt(t, TestNetAccount032Address),
					EntryPointSelector: types.GetSelectorFromNameFelt("get_nonce"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:               WithBlockTag("latest"),
				ExpectedPatternResult: "^0x[0-9a-f]+$",
			},
		},
		"mainnet": {
			{
				FunctionCall: FunctionCall{
					ContractAddress:    utils.TestHexToFelt(t, "0x06a09ccb1caaecf3d9683efe335a667b2169a409d19c589ba1eb771cd210af75"),
					EntryPointSelector: types.GetSelectorFromNameFelt("decimals"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:               WithBlockTag("latest"),
				ExpectedPatternResult: "^0x12$",
			},
		},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		output, err := testConfig.provider.Call(context.Background(), FunctionCall(test.FunctionCall), test.BlockID)
		if err != nil {
			t.Fatal(err)
		}
		if diff, err := spy.Compare(output, false); err != nil || diff != "FullMatch" {
			spy.Compare(output, true)
			t.Fatal("expecting to match", err)
		}
		if len(output) == 0 {
			t.Fatal("should return an output")
		}
		match, err := regexp.Match(test.ExpectedPatternResult, []byte(output[0]))
		if err != nil || !match {
			t.Fatalf("checking output(%v) expecting %s, got: %v", err, test.ExpectedPatternResult, output[0])
		}

	}
}
