package tests

import (
	"testing"

	coreMock "github.com/ElrondNetwork/elrond-go-core/core/mock"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/elrond-vm-common/mock"
	"github.com/ElrondNetwork/wasm-vm-v1_2/arwen"
	"github.com/ElrondNetwork/wasm-vm-v1_2/config"
	"github.com/ElrondNetwork/wasm-vm-v1_2/ipc/common"
	"github.com/ElrondNetwork/wasm-vm-v1_2/ipc/nodepart"
	contextmock "github.com/ElrondNetwork/wasm-vm-v1_2/mock/context"
	worldmock "github.com/ElrondNetwork/wasm-vm-v1_2/mock/world"
	"github.com/stretchr/testify/require"
)

var arwenVirtualMachine = []byte{5, 0}

func TestArwenDriver_DiagnoseWait(t *testing.T) {
	t.Skip("driver not supported anymore")
	blockchain := &contextmock.BlockchainHookStub{}
	adressGenerator := &coreMock.AddressGeneratorStub{}
	driver := newDriver(t, blockchain, adressGenerator)

	err := driver.DiagnoseWait(100)
	require.Nil(t, err)
}

func TestArwenDriver_DiagnoseWaitWithTimeout(t *testing.T) {
	t.Skip("driver not supported anymore: requires standalone binary")
	blockchain := &contextmock.BlockchainHookStub{}
	adressGenerator := &coreMock.AddressGeneratorStub{}
	driver := newDriver(t, blockchain, adressGenerator)

	err := driver.DiagnoseWait(5000)
	require.True(t, common.IsCriticalError(err))
	require.Contains(t, err.Error(), "timeout")
	require.True(t, driver.IsClosed())
}

func TestArwenDriver_RestartsIfStopped(t *testing.T) {
	t.Skip("driver not supported anymore")
	logger.ToggleLoggerName(true)
	_ = logger.SetLogLevel("*:TRACE")

	blockchain := &contextmock.BlockchainHookStub{}
	adressGenerator := &coreMock.AddressGeneratorStub{}
	driver := newDriver(t, blockchain, adressGenerator)

	blockchain.GetUserAccountCalled = func(address []byte) (vmcommon.UserAccountHandler, error) {
		return &worldmock.Account{Code: bytecodeCounter}, nil
	}

	vmOutput, err := driver.RunSmartContractCreate(createDeployInput(bytecodeCounter))
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	vmOutput, err = driver.RunSmartContractCall(createCallInput("increment"))
	require.Nil(t, err)
	require.NotNil(t, vmOutput)

	require.False(t, driver.IsClosed())
	_ = driver.Close()
	require.True(t, driver.IsClosed())

	// Per this request, Arwen is restarted
	vmOutput, err = driver.RunSmartContractCreate(createDeployInput(bytecodeCounter))
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.False(t, driver.IsClosed())
}

func BenchmarkArwenDriver_RestartsIfStopped(b *testing.B) {
	blockchain := &contextmock.BlockchainHookStub{}
	adressGenerator := &coreMock.AddressGeneratorStub{}
	driver := newDriver(b, blockchain, adressGenerator)

	for i := 0; i < b.N; i++ {
		_ = driver.Close()
		require.True(b, driver.IsClosed())
		_ = driver.RestartArwenIfNecessary()
		require.False(b, driver.IsClosed())
	}
}

func BenchmarkArwenDriver_RestartArwenIfNecessary(b *testing.B) {
	blockchain := &contextmock.BlockchainHookStub{}
	adressGenerator := &coreMock.AddressGeneratorStub{}
	driver := newDriver(b, blockchain, adressGenerator)

	for i := 0; i < b.N; i++ {
		_ = driver.RestartArwenIfNecessary()
	}
}

func TestArwenDriver_GetVersion(t *testing.T) {
	t.Skip("driver not supported anymore")
	// This test requires `make arwen` before running, or must be run directly
	// with `make test`
	blockchain := &contextmock.BlockchainHookStub{}
	adressGenerator := &coreMock.AddressGeneratorStub{}
	driver := newDriver(t, blockchain, adressGenerator)
	version := driver.GetVersion()
	require.NotZero(t, len(version))
	require.NotEqual(t, "undefined", version)
}

func newDriver(tb testing.TB, blockchain *contextmock.BlockchainHookStub, addressGenerator *coreMock.AddressGeneratorStub) *nodepart.ArwenDriver {
	driver, err := nodepart.NewArwenDriver(
		blockchain,
		addressGenerator,
		common.ArwenArguments{
			VMHostParameters: arwen.VMHostParameters{
				VMType:                   arwenVirtualMachine,
				BlockGasLimit:            uint64(10000000),
				GasSchedule:              config.MakeGasMapForTests(),
				ElrondProtectedKeyPrefix: []byte("ELROND"),
				EnableEpochsHandler: &mock.EnableEpochsHandlerStub{
					IsSCDeployFlagEnabledField:            true,
					IsAheadOfTimeGasUsageFlagEnabledField: true,
					IsRepairCallbackFlagEnabledField:      true,
					IsBuiltInFunctionsFlagEnabledField:    true,
				},
			},
		},
		nodepart.Config{MaxLoopTime: 1000},
	)
	require.Nil(tb, err)
	require.NotNil(tb, driver)
	require.False(tb, driver.IsClosed())
	return driver
}
