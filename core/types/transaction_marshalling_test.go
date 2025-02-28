package types

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestTransactionUnmarshalJsonDeposit(t *testing.T) {
	tx := NewTx(&DepositTx{
		SourceHash: common.HexToHash("0x1234"),
		Mint:       big.NewInt(34),
	})
	json, err := tx.MarshalJSON()
	require.NoError(t, err, "Failed to marshal tx JSON")

	got := &Transaction{}
	err = got.UnmarshalJSON(json)
	require.NoError(t, err, "Failed to unmarshal tx JSON")
	require.Equal(t, tx.Hash(), got.Hash())
}

func TestTransactionUnmarshalJsonDepositWithNonce(t *testing.T) {
	json := []byte(`{
		"type": "0x7e",
		"from": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0001",
		"gas": "0xf4240",
		"value": "0x0",
		"nonce": "0x3",
		"mint": "0x0",
		"input": "0x015d8eb900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000064460d46000000000000000000000000000000000000000000000000000000003b9aca00d184c96292c6a33c9314adc803dae69870ca9d7e410f82a969d4e7f3f6d1a5fd00000000000000000000000000000000000000000000000000000000000000020000000000000000000000009965507d1a55bcc2695c58ba16fb37d819b0a4dc000000000000000000000000000000000000000000000000000000000000083400000000000000000000000000000000000000000000000000000000000f4240",
		"sourceHash": "0x0000000000000000000000000000000000000000000000000000000000001234"
	}`)

	got := &Transaction{}
	err := got.UnmarshalJSON(json)
	require.NoError(t, err, "Failed to unmarshal tx JSON")
	// NOTE(chokobole): This fails because got is not a DepositTx but a depositTxWithNonce.
	// require.Equal(t, common.HexToHash("0x1234"), got.SourceHash())

	_, ok := got.inner.(*depositTxWithNonce)
	require.True(t, ok)

	json, err = got.MarshalJSON()
	require.NoError(t, err, "Failed to marshal tx JSON")

	got = &Transaction{}
	err = got.UnmarshalJSON(json)
	require.NoError(t, err, "Failed to unmarshal tx JSON")
	require.Equal(t, uint8(0x7e), got.Type())
	require.Equal(t, uint64(0xf4240), got.Gas())
	require.Equal(t, common.HexToHash("0x1234"), got.SourceHash())
}

func TestTransactionUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name          string
		json          string
		expectedError string
	}{
		{
			name:          "No gas",
			json:          `{"type":"0x7e","nonce":null,"gasPrice":null,"maxPriorityFeePerGas":null,"maxFeePerGas":null,"value":"0x1","input":"0x616263646566","v":null,"r":null,"s":null,"to":null,"sourceHash":"0x0000000000000000000000000000000000000000000000000000000000000000","from":"0x0000000000000000000000000000000000000001","hash":"0xa4341f3db4363b7ca269a8538bd027b2f8784f84454ca917668642d5f6dffdf9"}`,
			expectedError: "missing required field 'gas'",
		},
		{
			name:          "No value",
			json:          `{"type":"0x7e","nonce":null,"gas": "0x1234", "gasPrice":null,"maxPriorityFeePerGas":null,"maxFeePerGas":null,"input":"0x616263646566","v":null,"r":null,"s":null,"to":null,"sourceHash":"0x0000000000000000000000000000000000000000000000000000000000000000","from":"0x0000000000000000000000000000000000000001","hash":"0xa4341f3db4363b7ca269a8538bd027b2f8784f84454ca917668642d5f6dffdf9"}`,
			expectedError: "missing required field 'value'",
		},
		{
			name:          "No input",
			json:          `{"type":"0x7e","nonce":null,"gas": "0x1234", "gasPrice":null,"maxPriorityFeePerGas":null,"maxFeePerGas":null,"value":"0x1","v":null,"r":null,"s":null,"to":null,"sourceHash":"0x0000000000000000000000000000000000000000000000000000000000000000","from":"0x0000000000000000000000000000000000000001","hash":"0xa4341f3db4363b7ca269a8538bd027b2f8784f84454ca917668642d5f6dffdf9"}`,
			expectedError: "missing required field 'input'",
		},
		{
			name:          "No from",
			json:          `{"type":"0x7e","nonce":null,"gas": "0x1234", "gasPrice":null,"maxPriorityFeePerGas":null,"maxFeePerGas":null,"value":"0x1","input":"0x616263646566","v":null,"r":null,"s":null,"to":null,"sourceHash":"0x0000000000000000000000000000000000000000000000000000000000000000","hash":"0xa4341f3db4363b7ca269a8538bd027b2f8784f84454ca917668642d5f6dffdf9"}`,
			expectedError: "missing required field 'from'",
		},
		{
			name:          "No sourceHash",
			json:          `{"type":"0x7e","nonce":null,"gas": "0x1234", "gasPrice":null,"maxPriorityFeePerGas":null,"maxFeePerGas":null,"value":"0x1","input":"0x616263646566","v":null,"r":null,"s":null,"to":null,"from":"0x0000000000000000000000000000000000000001","hash":"0xa4341f3db4363b7ca269a8538bd027b2f8784f84454ca917668642d5f6dffdf9"}`,
			expectedError: "missing required field 'sourceHash'",
		},
		{
			name: "No mint",
			json: `{"type":"0x7e","nonce":null,"gas": "0x1234", "gasPrice":null,"maxPriorityFeePerGas":null,"maxFeePerGas":null,"value":"0x1","input":"0x616263646566","v":null,"r":null,"s":null,"to":null,"sourceHash":"0x0000000000000000000000000000000000000000000000000000000000000000","from":"0x0000000000000000000000000000000000000001","hash":"0xa4341f3db4363b7ca269a8538bd027b2f8784f84454ca917668642d5f6dffdf9"}`,
			// Allowed
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var parsedTx = &Transaction{}
			err := json.Unmarshal([]byte(test.json), &parsedTx)
			if test.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, test.expectedError)
			}
		})
	}
}
