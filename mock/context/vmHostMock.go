package mock

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/wasm-vm-v1_2/crypto"
	"github.com/ElrondNetwork/wasm-vm-v1_2/wasmer"
	"github.com/ElrondNetwork/wasm-vm-v1_2/wasmvm"
)

var _ wasmvm.VMHost = (*VMHostMock)(nil)

// VMHostMock is used in tests to check the VMHost interface method calls
type VMHostMock struct {
	BlockChainHook vmcommon.BlockchainHook
	CryptoHook     crypto.VMCrypto

	EthInput []byte

	BlockchainContext wasmvm.BlockchainContext
	RuntimeContext    wasmvm.RuntimeContext
	OutputContext     wasmvm.OutputContext
	MeteringContext   wasmvm.MeteringContext
	StorageContext    wasmvm.StorageContext
	BigIntContext     wasmvm.BigIntContext

	SCAPIMethods  *wasmer.Imports
	IsBuiltinFunc bool
}

// Crypto mocked method
func (host *VMHostMock) Crypto() crypto.VMCrypto {
	return host.CryptoHook
}

// Blockchain mocked method
func (host *VMHostMock) Blockchain() wasmvm.BlockchainContext {
	return host.BlockchainContext
}

// Runtime mocked method
func (host *VMHostMock) Runtime() wasmvm.RuntimeContext {
	return host.RuntimeContext
}

// Output mocked method
func (host *VMHostMock) Output() wasmvm.OutputContext {
	return host.OutputContext
}

// Metering mocked method
func (host *VMHostMock) Metering() wasmvm.MeteringContext {
	return host.MeteringContext
}

// Storage mocked method
func (host *VMHostMock) Storage() wasmvm.StorageContext {
	return host.StorageContext
}

// BigInt mocked method
func (host *VMHostMock) BigInt() wasmvm.BigIntContext {
	return host.BigIntContext
}

// IsV2Enabled mocked method
func (host *VMHostMock) IsV2Enabled() bool {
	return true
}

// IsV3Enabled mocked method
func (host *VMHostMock) IsV3Enabled() bool {
	return true
}

// IsAheadOfTimeCompileEnabled mocked method
func (host *VMHostMock) IsAheadOfTimeCompileEnabled() bool {
	return true
}

// IsDynamicGasLockingEnabled mocked method
func (host *VMHostMock) IsDynamicGasLockingEnabled() bool {
	return true
}

// IsESDTFunctionsEnabled mocked method
func (host *VMHostMock) IsESDTFunctionsEnabled() bool {
	return true
}

// AreInSameShard mocked method
func (host *VMHostMock) AreInSameShard(_ []byte, _ []byte) bool {
	return true
}

// RevertESDTTransfer mocked method
func (host *VMHostMock) RevertESDTTransfer(_ *vmcommon.ContractCallInput) {
}

// ExecuteESDTTransfer mocked method
func (host *VMHostMock) ExecuteESDTTransfer(_ []byte, _ []byte, _ []byte, _ uint64, _ *big.Int, _ vm.CallType, _ bool) (*vmcommon.VMOutput, uint64, error) {
	return nil, 0, nil
}

// CreateNewContract mocked method
func (host *VMHostMock) CreateNewContract(_ *vmcommon.ContractCreateInput) ([]byte, error) {
	return nil, nil
}

// ExecuteOnSameContext mocked method
func (host *VMHostMock) ExecuteOnSameContext(_ *vmcommon.ContractCallInput) (*wasmvm.AsyncContextInfo, error) {
	return nil, nil
}

// ExecuteOnDestContext mocked method
func (host *VMHostMock) ExecuteOnDestContext(_ *vmcommon.ContractCallInput) (*vmcommon.VMOutput, *wasmvm.AsyncContextInfo, uint64, error) {
	return nil, nil, 0, nil
}

// InitState mocked method
func (host *VMHostMock) InitState() {
}

// PushState mocked method
func (host *VMHostMock) PushState() {
}

// PopState mocked method
func (host *VMHostMock) PopState() {
}

// ClearStateStack mocked method
func (host *VMHostMock) ClearStateStack() {
}

// GetAPIMethods mocked method
func (host *VMHostMock) GetAPIMethods() *wasmer.Imports {
	return host.SCAPIMethods
}

// GetProtocolBuiltinFunctions mocked method
func (host *VMHostMock) GetProtocolBuiltinFunctions() vmcommon.FunctionNames {
	return make(vmcommon.FunctionNames)
}

// IsBuiltinFunctionName mocked method
func (host *VMHostMock) IsBuiltinFunctionName(_ string) bool {
	return host.IsBuiltinFunc
}
