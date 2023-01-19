package nodepart

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-logger-go/pipes"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-v1_2-go/ipc/common"
	"github.com/multiversx/mx-chain-vm-v1_2-go/ipc/marshaling"
)

var log = logger.GetOrCreate("vmDriver")

var _ vmcommon.VMExecutionHandler = (*VMDriver)(nil)

// VMDriver manages the execution of the VM process
type VMDriver struct {
	blockchainHook      vmcommon.BlockchainHook
	vmArguments      common.VMArguments
	config              Config
	logsMarshalizer     marshaling.Marshalizer
	messagesMarshalizer marshaling.Marshalizer

	vmInitRead    *os.File
	vmInitWrite   *os.File
	vmInputRead   *os.File
	vmInputWrite  *os.File
	vmOutputRead  *os.File
	vmOutputWrite *os.File

	counterDeploy uint64
	counterCall   uint64

	command  *exec.Cmd
	part     *NodePart
	logsPart ParentLogsPart

	// When the VMDriver is used to resolve contract queries, it might happen that a query request executes concurrently with other operations (such as "GasScheduleChange").
	// Query requests are ordered sequentially within the API layer (see the QueryService dispatcher and other related components), but this sequence of queries might
	// interleave with VM-management operations, which are or might be triggered within a different flow (e.g. the processing flow). For example, "GasScheduleChange" is triggered synchronously
	// with the processing flow (on a certain epoch change), but in asynchronicity with the querying flow.
	// This might lead to issues (such as interleaving message sequences on the communication pipes).
	// A solution is to use a mutex, and treat each operation within a critical section (in the VMDriver, thus on node's part).
	// Thus, for any two concurrent operations, the first one reaching the mutex also wins the pipe and holds ownership upon its completion.
	operationsMutex sync.Mutex
}

// NewVMDriver creates a new driver
func NewVMDriver(
	blockchainHook vmcommon.BlockchainHook,
	vmArguments common.VMArguments,
	config Config,
) (*VMDriver, error) {
	driver := &VMDriver{
		blockchainHook:      blockchainHook,
		vmArguments:      vmArguments,
		config:              config,
		logsMarshalizer:     marshaling.CreateMarshalizer(vmArguments.LogsMarshalizer),
		messagesMarshalizer: marshaling.CreateMarshalizer(vmArguments.MessagesMarshalizer),
	}

	err := driver.startVM()
	if err != nil {
		return nil, err
	}

	return driver, nil
}

func (driver *VMDriver) startVM() error {
	log.Info("VMDriver.startVM()")

	logsProfileReader, logsWriter, err := driver.resetLogsPart()
	if err != nil {
		return err
	}

	err = driver.resetPipeStreams()
	if err != nil {
		return err
	}

	vmPath, err := driver.getVMPath()
	if err != nil {
		return err
	}

	driver.command = exec.Command(vmPath)
	driver.command.ExtraFiles = []*os.File{
		driver.vmInitRead,
		driver.vmInputRead,
		driver.vmOutputWrite,
		logsProfileReader,
		logsWriter,
	}

	vmStdout, err := driver.command.StdoutPipe()
	if err != nil {
		return err
	}

	vmStderr, err := driver.command.StderrPipe()
	if err != nil {
		return err
	}

	err = driver.command.Start()
	if err != nil {
		return err
	}

	err = common.SendVMArguments(driver.vmInitWrite, driver.vmArguments)
	if err != nil {
		return err
	}

	driver.blockchainHook.ClearCompiledCodes()

	driver.part, err = NewNodePart(
		driver.vmOutputRead,
		driver.vmInputWrite,
		driver.blockchainHook,
		driver.config,
		driver.messagesMarshalizer,
	)
	if err != nil {
		return err
	}

	err = driver.logsPart.StartLoop(vmStdout, vmStderr)
	if err != nil {
		return err
	}

	return nil
}

func (driver *VMDriver) resetLogsPart() (*os.File, *os.File, error) {
	logsPart, err := pipes.NewParentPart("VM", driver.logsMarshalizer)
	if err != nil {
		return nil, nil, err
	}

	driver.logsPart = logsPart
	readProfile, writeLogs := logsPart.GetChildPipes()
	return readProfile, writeLogs, nil
}

func (driver *VMDriver) resetPipeStreams() error {
	closeFile(driver.vmInitRead)
	closeFile(driver.vmInitWrite)
	closeFile(driver.vmInputRead)
	closeFile(driver.vmInputWrite)
	closeFile(driver.vmOutputRead)
	closeFile(driver.vmOutputWrite)

	var err error

	driver.vmInitRead, driver.vmInitWrite, err = os.Pipe()
	if err != nil {
		return err
	}

	driver.vmInputRead, driver.vmInputWrite, err = os.Pipe()
	if err != nil {
		return err
	}

	driver.vmOutputRead, driver.vmOutputWrite, err = os.Pipe()
	if err != nil {
		return err
	}

	return nil
}

func closeFile(file *os.File) {
	if file != nil {
		err := file.Close()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Cannot close file.\n")
		}
	}
}

// RestartVMIfNecessary restarts VM if the process is closed
func (driver *VMDriver) RestartVMIfNecessary() error {
	if !driver.IsClosed() {
		return nil
	}

	err := driver.startVM()
	return err
}

// IsClosed checks whether the VM process is closed
func (driver *VMDriver) IsClosed() bool {
	pid := driver.command.Process.Pid
	process, err := os.FindProcess(pid)
	if err != nil {
		return true
	}

	err = process.Signal(syscall.Signal(0))
	return err != nil
}

// GetVersion gets the VM version
func (driver *VMDriver) GetVersion() string {
	driver.operationsMutex.Lock()
	defer driver.operationsMutex.Unlock()

	log.Trace("GetVersion")

	err := driver.RestartVMIfNecessary()
	if err != nil {
		log.Warn("GetVersion", "err", err)
		return ""
	}

	request := common.NewMessageVersionRequest()
	response, err := driver.part.StartLoop(request)
	if err != nil {
		log.Warn("GetVersion", "err", err)
		_ = driver.Close()
		return ""
	}

	typedResponse := response.(*common.MessageVersionResponse)

	return typedResponse.Version
}

// GasScheduleChange sends a "gas change" request to VM and waits for the output
func (driver *VMDriver) GasScheduleChange(newGasSchedule map[string]map[string]uint64) {
	driver.operationsMutex.Lock()
	defer driver.operationsMutex.Unlock()

	driver.vmArguments.GasSchedule = newGasSchedule
	err := driver.RestartVMIfNecessary()
	if err != nil {
		log.Error("GasScheduleChange RestartVMIfNecessary", "error", err)
		return
	}

	request := common.NewMessageGasScheduleChangeRequest(newGasSchedule)
	response, err := driver.part.StartLoop(request)
	if err != nil {
		log.Error("GasScheduleChange StartLoop", "error", err)
		_ = driver.Close()
		return
	}

	if response.GetError() != nil {
		log.Error("GasScheduleChange StartLoop response", "error", err)
		_ = driver.Close()
		return
	}
}

// RunSmartContractCreate sends a deploy request to VM and waits for the output
func (driver *VMDriver) RunSmartContractCreate(input *vmcommon.ContractCreateInput) (*vmcommon.VMOutput, error) {
	driver.operationsMutex.Lock()
	defer driver.operationsMutex.Unlock()

	driver.counterDeploy++
	log.Trace("RunSmartContractCreate", "counter", driver.counterDeploy)

	err := driver.RestartVMIfNecessary()
	if err != nil {
		return nil, common.WrapCriticalError(err)
	}

	request := common.NewMessageContractDeployRequest(input)
	response, err := driver.part.StartLoop(request)
	if err != nil {
		log.Warn("RunSmartContractCreate", "err", err)
		_ = driver.Close()
		return nil, common.WrapCriticalError(err)
	}

	typedResponse := response.(*common.MessageContractResponse)
	vmOutput, err := typedResponse.SerializableVMOutput.ConvertToVMOutput(), response.GetError()
	if err != nil {
		return nil, err
	}

	return vmOutput, nil
}

// RunSmartContractCall sends an execution request to VM and waits for the output
func (driver *VMDriver) RunSmartContractCall(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	driver.operationsMutex.Lock()
	defer driver.operationsMutex.Unlock()

	driver.counterCall++
	log.Trace("RunSmartContractCall", "counter", driver.counterCall, "func", input.Function, "sc", input.RecipientAddr)

	err := driver.RestartVMIfNecessary()
	if err != nil {
		return nil, common.WrapCriticalError(err)
	}

	request := common.NewMessageContractCallRequest(input)
	response, err := driver.part.StartLoop(request)
	if err != nil {
		log.Warn("RunSmartContractCall", "err", err)
		_ = driver.Close()
		return nil, common.WrapCriticalError(err)
	}

	typedResponse := response.(*common.MessageContractResponse)
	vmOutput, err := typedResponse.SerializableVMOutput.ConvertToVMOutput(), response.GetError()
	if err != nil {
		return nil, err
	}

	return vmOutput, nil
}

// DiagnoseWait sends a diagnose message to VM
func (driver *VMDriver) DiagnoseWait(milliseconds uint32) error {
	driver.operationsMutex.Lock()
	defer driver.operationsMutex.Unlock()

	err := driver.RestartVMIfNecessary()
	if err != nil {
		return common.WrapCriticalError(err)
	}

	request := common.NewMessageDiagnoseWaitRequest(milliseconds)
	response, err := driver.part.StartLoop(request)
	if err != nil {
		log.Error("DiagnoseWait", "err", err)
		_ = driver.Close()
		return common.WrapCriticalError(err)
	}

	return response.GetError()
}

// Close stops VM
func (driver *VMDriver) Close() error {
	driver.logsPart.StopLoop()

	err := driver.stopVM()
	if err != nil {
		log.Error("VMDriver.Close()", "err", err)
		return err
	}

	return nil
}

func (driver *VMDriver) stopVM() error {
	err := driver.command.Process.Kill()
	if err != nil {
		return err
	}

	_, err = driver.command.Process.Wait()
	if err != nil {
		return err
	}

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (driver *VMDriver) IsInterfaceNil() bool {
	return driver == nil
}
