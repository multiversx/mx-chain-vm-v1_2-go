package contexts

import (
	"testing"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/wasm-vm-v1_2/wasmvm"
	"github.com/stretchr/testify/require"
)

func TestReservedFunctions_IsFunctionReserved(t *testing.T) {
	scAPINames := vmcommon.FunctionNames{
		"rockets": {},
	}

	fromProtocol := vmcommon.FunctionNames{
		"protocolFunctionFoo": {},
		"protocolFunctionBar": {},
	}

	reserved := NewReservedFunctions(scAPINames, fromProtocol)

	require.False(t, reserved.IsReserved("foo"))
	require.True(t, reserved.IsReserved("rockets"))
	require.True(t, reserved.IsReserved("protocolFunctionFoo"))
	require.True(t, reserved.IsReserved("protocolFunctionBar"))
	require.True(t, reserved.IsReserved(wasmvm.UpgradeFunctionName))
}
