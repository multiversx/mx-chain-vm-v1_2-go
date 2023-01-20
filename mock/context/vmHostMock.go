package mock

import (
	"math/big"

	"github.com/multiversx/mx-chain-core-go/data/vm"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-v1_2-go/crypto"
	"github.com/multiversx/mx-chain-vm-v1_2-go/vmhost"
	"github.com/multiversx/mx-chain-vm-v1_2-go/wasmer"
)

var _ vmhost.VMHost = (*VMHostMock)(nil)

// VMHostMock is used in tests to check the VMHost interface method calls
type VMHostMock struct {
	BlockChainHook vmcommon.BlockchainHook
	CryptoHook     crypto.VMCrypto

	EthInput []byte

	BlockchainContext vmhost.BlockchainContext
	RuntimeContext    vmhost.RuntimeContext
	OutputContext     vmhost.OutputContext
	MeteringContext   vmhost.MeteringContext
	StorageContext    vmhost.StorageContext
	BigIntContext     vmhost.BigIntContext

	SCAPIMethods  *wasmer.Imports
	IsBuiltinFunc bool
}

// Crypto mocked method
func (host *VMHostMock) Crypto() crypto.VMCrypto {
	return host.CryptoHook
}

// Blockchain mocked method
func (host *VMHostMock) Blockchain() vmhost.BlockchainContext {
	return host.BlockchainContext
}

// Runtime mocked method
func (host *VMHostMock) Runtime() vmhost.RuntimeContext {
	return host.RuntimeContext
}

// Output mocked method
func (host *VMHostMock) Output() vmhost.OutputContext {
	return host.OutputContext
}

// Metering mocked method
func (host *VMHostMock) Metering() vmhost.MeteringContext {
	return host.MeteringContext
}

// Storage mocked method
func (host *VMHostMock) Storage() vmhost.StorageContext {
	return host.StorageContext
}

// BigInt mocked method
func (host *VMHostMock) BigInt() vmhost.BigIntContext {
	return host.BigIntContext
}

// IsVMV2Enabled mocked method
func (host *VMHostMock) IsVMV2Enabled() bool {
	return true
}

// IsVMV3Enabled mocked method
func (host *VMHostMock) IsVMV3Enabled() bool {
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
func (host *VMHostMock) ExecuteOnSameContext(_ *vmcommon.ContractCallInput) (*vmhost.AsyncContextInfo, error) {
	return nil, nil
}

// ExecuteOnDestContext mocked method
func (host *VMHostMock) ExecuteOnDestContext(_ *vmcommon.ContractCallInput) (*vmcommon.VMOutput, *vmhost.AsyncContextInfo, uint64, error) {
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
