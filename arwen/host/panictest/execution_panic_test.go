package panictest

import (
	"math/big"
	"testing"

	arwenHost "github.com/multiversx/mx-chain-vm-v1_2-go/arwen/host"
	"github.com/stretchr/testify/require"
)

func TestExecution_PanicInGoWithSilentWasmer_SIGSEGV(t *testing.T) {
	code := arwenHost.GetTestSCCode("counter", "../../../")
	host, blockchain := arwenHost.DefaultTestArwenForCallSigSegv(t, code, big.NewInt(1), true)

	blockchain.GetStorageDataCalled = func(_ []byte, _ []byte) ([]byte, uint32, error) {
		var i *int
		i = nil

		// dereference a nil pointer
		*i = *i + 1
		return nil, 0, nil
	}

	input := arwenHost.DefaultTestContractCallInput()
	input.GasProvided = 10000000
	input.Function = "increment"

	// Ensure that no more panic
	defer func() {
		r := recover()
		require.Nil(t, r)
	}()

	expectedError := "runtime error: invalid memory address or nil pointer dereference"

	_, err := host.RunSmartContractCall(input)
	require.Equal(t, expectedError, err.Error())
}
