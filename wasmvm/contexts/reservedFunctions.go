package contexts

import (
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/wasm-vm-v1_2/wasmvm"
)

// reservedFunctions holds the reserved function names
type reservedFunctions struct {
	functionNames vmcommon.FunctionNames
}

// NewReservedFunctions creates a new reservedFunctions
func NewReservedFunctions(scAPINames vmcommon.FunctionNames, protocolBuiltinFunctions vmcommon.FunctionNames) *reservedFunctions {
	result := &reservedFunctions{
		functionNames: make(vmcommon.FunctionNames),
	}

	for name, value := range protocolBuiltinFunctions {
		result.functionNames[name] = value
	}

	for name, value := range scAPINames {
		result.functionNames[name] = value
	}

	var empty struct{}
	result.functionNames[wasmvm.UpgradeFunctionName] = empty

	return result
}

// IsReserved returns whether a function is reserved
func (reservedFunctions *reservedFunctions) IsReserved(functionName string) bool {
	if _, ok := reservedFunctions.functionNames[functionName]; ok {
		return true
	}

	return false
}

// GetReserved gets the reserved functions as a slice of strings
func (reservedFunctions *reservedFunctions) GetReserved() []string {
	keys := make([]string, len(reservedFunctions.functionNames))

	i := 0
	for key := range reservedFunctions.functionNames {
		keys[i] = key
		i++
	}

	return keys
}
