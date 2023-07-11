package hostCore

import (
	"fmt"
	"sync"

	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-v1_2-go/config"
	"github.com/multiversx/mx-chain-vm-v1_2-go/crypto"
	"github.com/multiversx/mx-chain-vm-v1_2-go/crypto/factory"
	"github.com/multiversx/mx-chain-vm-v1_2-go/vmhost"
	"github.com/multiversx/mx-chain-vm-v1_2-go/vmhost/contexts"
	"github.com/multiversx/mx-chain-vm-v1_2-go/vmhost/cryptoapi"
	"github.com/multiversx/mx-chain-vm-v1_2-go/vmhost/vmhooks"
	"github.com/multiversx/mx-chain-vm-v1_2-go/wasmer"
)

var log = logger.GetOrCreate("vm/host")

// MaximumWasmerInstanceCount represents the maximum number of Wasmer instances that can be active at the same time
var MaximumWasmerInstanceCount = uint64(10)

// TryFunction corresponds to the try() part of a try / catch block
type TryFunction func()

// CatchFunction corresponds to the catch() part of a try / catch block
type CatchFunction func(error)

// vmHost implements HostContext interface.
type vmHost struct {
	blockChainHook vmcommon.BlockchainHook
	cryptoHook     crypto.VMCrypto
	mutExecution   sync.RWMutex

	ethInput []byte

	blockchainContext vmhost.BlockchainContext
	runtimeContext    vmhost.RuntimeContext
	outputContext     vmhost.OutputContext
	meteringContext   vmhost.MeteringContext
	storageContext    vmhost.StorageContext
	bigIntContext     vmhost.BigIntContext

	gasSchedule              config.GasScheduleMap
	scAPIMethods             *wasmer.Imports
	protocolBuiltinFunctions vmcommon.FunctionNames
	enableEpochsHandler      vmhost.EnableEpochsHandler
}

// NewVMHost creates a new VM vmHost
func NewVMHost(
	blockChainHook vmcommon.BlockchainHook,
	hostParameters *vmhost.VMHostParameters,
) (*vmHost, error) {
	if check.IfNil(hostParameters.EnableEpochsHandler) {
		return nil, vmhost.ErrNilEnableEpochsHandler
	}

	cryptoHook := factory.NewVMCrypto()
	host := &vmHost{
		blockChainHook:           blockChainHook,
		cryptoHook:               cryptoHook,
		meteringContext:          nil,
		runtimeContext:           nil,
		blockchainContext:        nil,
		storageContext:           nil,
		bigIntContext:            nil,
		gasSchedule:              hostParameters.GasSchedule,
		scAPIMethods:             nil,
		protocolBuiltinFunctions: hostParameters.ProtocolBuiltinFunctions,
		enableEpochsHandler:      hostParameters.EnableEpochsHandler,
	}

	var err error

	imports, err := vmhooks.BaseOpsAPIImports()
	if err != nil {
		return nil, err
	}

	imports, err = vmhooks.BigIntImports(imports)
	if err != nil {
		return nil, err
	}

	imports, err = vmhooks.SmallIntImports(imports)
	if err != nil {
		return nil, err
	}

	imports, err = cryptoapi.CryptoImports(imports)
	if err != nil {
		return nil, err
	}

	err = wasmer.SetImports(imports)
	if err != nil {
		return nil, err
	}

	host.scAPIMethods = imports

	host.blockchainContext, err = contexts.NewBlockchainContext(host, blockChainHook)
	if err != nil {
		return nil, err
	}

	host.runtimeContext, err = contexts.NewRuntimeContext(
		host,
		hostParameters.VMType,
		hostParameters.UseWarmInstance,
	)
	if err != nil {
		return nil, err
	}

	host.meteringContext, err = contexts.NewMeteringContext(host, hostParameters.GasSchedule, hostParameters.BlockGasLimit)
	if err != nil {
		return nil, err
	}

	host.outputContext, err = contexts.NewOutputContext(host)
	if err != nil {
		return nil, err
	}

	host.storageContext, err = contexts.NewStorageContext(host, blockChainHook, hostParameters.ProtectedKeyPrefix)
	if err != nil {
		return nil, err
	}

	host.bigIntContext, err = contexts.NewBigIntContext()
	if err != nil {
		return nil, err
	}

	gasCostConfig, err := config.CreateGasConfig(host.gasSchedule)
	if err != nil {
		return nil, err
	}

	host.runtimeContext.SetMaxInstanceCount(MaximumWasmerInstanceCount)

	opcodeCosts := gasCostConfig.WASMOpcodeCost.ToOpcodeCostsArray()
	wasmer.SetOpcodeCosts(&opcodeCosts)

	if hostParameters.WasmerSIGSEGVPassthrough {
		wasmer.SetSIGSEGVPassthrough()
	}

	wasmer.ForceInstallSighandlers()

	host.initContexts()

	return host, nil
}

// GetVersion returns the VM version string
func (host *vmHost) GetVersion() string {
	return vmhost.VMVersion
}

// Crypto returns the VMCrypto instance of the host
func (host *vmHost) Crypto() crypto.VMCrypto {
	return host.cryptoHook
}

// Blockchain returns the BlockchainContext instance of the host
func (host *vmHost) Blockchain() vmhost.BlockchainContext {
	return host.blockchainContext
}

// Runtime returns the RuntimeContext instance of the host
func (host *vmHost) Runtime() vmhost.RuntimeContext {
	return host.runtimeContext
}

// Output returns the OutputContext instance of the host
func (host *vmHost) Output() vmhost.OutputContext {
	return host.outputContext
}

// Metering returns the MeteringContext instance of the host
func (host *vmHost) Metering() vmhost.MeteringContext {
	return host.meteringContext
}

// Storage returns the StorageContext instance of the host
func (host *vmHost) Storage() vmhost.StorageContext {
	return host.storageContext
}

// BigInt returns the BigIntContext instance of the host
func (host *vmHost) BigInt() vmhost.BigIntContext {
	return host.bigIntContext
}

// IsVMV2Enabled returns whether the VM V2 mode is enabled
func (host *vmHost) IsVMV2Enabled() bool {
	return host.enableEpochsHandler.IsSCDeployFlagEnabledInEpoch(host.enableEpochsHandler.GetCurrentEpoch())
}

// IsVMV3Enabled returns whether the V3 features are enabled
func (host *vmHost) IsVMV3Enabled() bool {
	return host.enableEpochsHandler.IsRepairCallbackFlagEnabledInEpoch(host.enableEpochsHandler.GetCurrentEpoch())
}

// IsAheadOfTimeCompileEnabled returns whether ahead-of-time compilation is enabled
func (host *vmHost) IsAheadOfTimeCompileEnabled() bool {
	return host.enableEpochsHandler.IsAheadOfTimeGasUsageFlagEnabledInEpoch(host.enableEpochsHandler.GetCurrentEpoch())
}

// IsDynamicGasLockingEnabled returns whether dynamic gas locking mode is enabled
func (host *vmHost) IsDynamicGasLockingEnabled() bool {
	return host.enableEpochsHandler.IsSCDeployFlagEnabledInEpoch(host.enableEpochsHandler.GetCurrentEpoch())
}

// IsESDTFunctionsEnabled returns whether ESDT functions are enabled
func (host *vmHost) IsESDTFunctionsEnabled() bool {
	return host.enableEpochsHandler.IsBuiltInFunctionsFlagEnabledInEpoch(host.enableEpochsHandler.GetCurrentEpoch())
}

// GetContexts returns the main contexts of the host
func (host *vmHost) GetContexts() (
	vmhost.BigIntContext,
	vmhost.BlockchainContext,
	vmhost.MeteringContext,
	vmhost.OutputContext,
	vmhost.RuntimeContext,
	vmhost.StorageContext,
) {
	return host.bigIntContext,
		host.blockchainContext,
		host.meteringContext,
		host.outputContext,
		host.runtimeContext,
		host.storageContext
}

// InitState resets the contexts of the host and reconfigures its flags
func (host *vmHost) InitState() {
	host.initContexts()
}

func (host *vmHost) initContexts() {
	host.ClearContextStateStack()
	host.bigIntContext.InitState()
	host.outputContext.InitState()
	host.meteringContext.InitState()
	host.runtimeContext.InitState()
	host.storageContext.InitState()
	host.ethInput = nil
}

// ClearContextStateStack cleans the state stacks of all the contexts of the host
func (host *vmHost) ClearContextStateStack() {
	host.bigIntContext.ClearStateStack()
	host.outputContext.ClearStateStack()
	host.meteringContext.ClearStateStack()
	host.runtimeContext.ClearStateStack()
	host.storageContext.ClearStateStack()
}

// Clean closes the currently running Wasmer instance
func (host *vmHost) Clean() {
	if host.runtimeContext.IsWarmInstance() {
		return
	}
	host.runtimeContext.CleanWasmerInstance()
}

// GetAPIMethods returns the EEI as a set of imports for Wasmer
func (host *vmHost) GetAPIMethods() *wasmer.Imports {
	return host.scAPIMethods
}

// GetProtocolBuiltinFunctions returns the names of the built-in functions, reserved by the protocol
func (host *vmHost) GetProtocolBuiltinFunctions() vmcommon.FunctionNames {
	return host.protocolBuiltinFunctions
}

// GasScheduleChange applies a new gas schedule to the host
func (host *vmHost) GasScheduleChange(newGasSchedule config.GasScheduleMap) {
	host.mutExecution.Lock()
	defer host.mutExecution.Unlock()

	host.gasSchedule = newGasSchedule
	gasCostConfig, err := config.CreateGasConfig(newGasSchedule)
	if err != nil {
		log.Error("cannot apply new gas config remained with old one")
		return
	}

	opcodeCosts := gasCostConfig.WASMOpcodeCost.ToOpcodeCostsArray()
	wasmer.SetOpcodeCosts(&opcodeCosts)

	host.meteringContext.SetGasSchedule(newGasSchedule)
}

// GetGasScheduleMap returns the currently stored gas schedule
func (host *vmHost) GetGasScheduleMap() config.GasScheduleMap {
	return host.gasSchedule
}

// RunSmartContractCreate executes the deployment of a new contract
func (host *vmHost) RunSmartContractCreate(input *vmcommon.ContractCreateInput) (vmOutput *vmcommon.VMOutput, err error) {
	host.mutExecution.RLock()
	defer host.mutExecution.RUnlock()

	log.Trace("RunSmartContractCreate begin", "len(code)", len(input.ContractCode), "metadata", input.ContractCodeMetadata)

	try := func() {
		vmOutput = host.doRunSmartContractCreate(input)
	}

	catch := func(caught error) {
		err = caught
		log.Error("RunSmartContractCreate", "error", err)
	}

	TryCatch(try, catch, "vmhost.RunSmartContractCreate")
	if vmOutput != nil {
		log.Trace("RunSmartContractCreate end", "returnCode", vmOutput.ReturnCode, "returnMessage", vmOutput.ReturnMessage)
	}

	return
}

// RunSmartContractCall executes the call of an existing contract
func (host *vmHost) RunSmartContractCall(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput, err error) {
	host.mutExecution.RLock()
	defer host.mutExecution.RUnlock()

	log.Trace("RunSmartContractCall begin", "function", input.Function)

	tryUpgrade := func() {
		vmOutput = host.doRunSmartContractUpgrade(input)
	}

	tryCall := func() {
		vmOutput = host.doRunSmartContractCall(input)

		if host.hasRetriableExecutionError(vmOutput) {
			log.Error("Retriable execution error detected. Will reset warm Wasmer instance.")
			host.runtimeContext.ResetWarmInstance()
		}
	}

	catch := func(caught error) {
		err = caught
		log.Error("RunSmartContractCall", "error", err)
	}

	isUpgrade := input.Function == vmhost.UpgradeFunctionName
	if isUpgrade {
		TryCatch(tryUpgrade, catch, "vmhost.RunSmartContractUpgrade")
	} else {
		TryCatch(tryCall, catch, "vmhost.RunSmartContractCall")
	}

	return
}

// Close closes all internal instances of the vm
func (host *vmHost) Close() error {
	return nil
}

// TryCatch simulates a try/catch block using golang's recover() functionality
func TryCatch(try TryFunction, catch CatchFunction, catchFallbackMessage string) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%s, panic: %v", catchFallbackMessage, r)
			}

			catch(err)
		}
	}()

	try()
}

func (host *vmHost) hasRetriableExecutionError(vmOutput *vmcommon.VMOutput) bool {
	if !host.runtimeContext.IsWarmInstance() {
		return false
	}

	return vmOutput.ReturnMessage == "allocation error"
}

// AreInSameShard returns true if the provided addresses are part of the same shard
func (host *vmHost) AreInSameShard(leftAddress []byte, rightAddress []byte) bool {
	blockchain := host.Blockchain()
	leftShard := blockchain.GetShardOfAddress(leftAddress)
	rightShard := blockchain.GetShardOfAddress(rightAddress)

	return leftShard == rightShard
}

// IsInterfaceNil returns true if there is no value under the interface
func (host *vmHost) IsInterfaceNil() bool {
	return host == nil
}
