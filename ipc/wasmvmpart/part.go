package wasmvmpart

import (
	"os"
	"time"

	logger "github.com/multiversx/mx-chain-logger-go"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go-v1_2/ipc/common"
	"github.com/multiversx/mx-chain-vm-go-v1_2/ipc/marshaling"
	"github.com/multiversx/mx-chain-vm-go-v1_2/wasmvm"
	"github.com/multiversx/mx-chain-vm-go-v1_2/wasmvm/host"
)

var log = logger.GetOrCreate("wasmvm/part")

// VMPart is the endpoint that implements the message loop on VM's side
type VMPart struct {
	Messenger *VMMessenger
	VMHost    vmcommon.VMExecutionHandler
	Repliers  []common.MessageReplier
	Version   string
}

// NewVMPart creates the VM part
func NewVMPart(
	version string,
	input *os.File,
	output *os.File,
	vmHostParameters *wasmvm.VMHostParameters,
	marshalizer marshaling.Marshalizer,
) (*VMPart, error) {
	messenger := NewVMMessenger(input, output, marshalizer)
	blockchain := NewBlockchainHookGateway(messenger)

	newVMHost, err := host.NewWASMVM(
		blockchain,
		vmHostParameters,
	)
	if err != nil {
		return nil, err
	}

	part := &VMPart{
		Messenger: messenger,
		VMHost:    newVMHost,
		Version:   version,
	}

	part.Repliers = common.CreateReplySlots(part.noopReplier)
	part.Repliers[common.ContractDeployRequest] = part.replyToRunSmartContractCreate
	part.Repliers[common.ContractCallRequest] = part.replyToRunSmartContractCall
	part.Repliers[common.DiagnoseWaitRequest] = part.replyToDiagnoseWait
	part.Repliers[common.VersionRequest] = part.replyToVersionRequest
	part.Repliers[common.GasScheduleChangeRequest] = part.replyToGasScheduleChange

	return part, nil
}

func (part *VMPart) noopReplier(_ common.MessageHandler) common.MessageHandler {
	log.Error("noopReplier called")
	return common.CreateMessage(common.UndefinedRequestOrResponse)
}

// StartLoop runs the main loop
func (part *VMPart) StartLoop() error {
	part.Messenger.Reset()
	err := part.doLoop()
	part.Messenger.Shutdown()
	log.Error("end of loop", "err", err)
	return err
}

// doLoop ends only when a critical failure takes place
func (part *VMPart) doLoop() error {
	for {
		request, err := part.Messenger.ReceiveNodeRequest()
		if err != nil {
			return err
		}
		if common.IsStopRequest(request) {
			return common.ErrStopPerNodeRequest
		}

		response := part.replyToNodeRequest(request)

		// Successful execution, send response
		err = part.Messenger.SendContractResponse(response)
		if err != nil {
			return err
		}

		part.Messenger.ResetDialogue()
	}
}

func (part *VMPart) replyToNodeRequest(request common.MessageHandler) common.MessageHandler {
	replier := part.Repliers[request.GetKind()]
	return replier(request)
}

func (part *VMPart) replyToRunSmartContractCreate(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageContractDeployRequest)
	vmOutput, err := part.VMHost.RunSmartContractCreate(typedRequest.CreateInput)
	return common.NewMessageContractResponse(vmOutput, err)
}

func (part *VMPart) replyToRunSmartContractCall(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageContractCallRequest)
	vmOutput, err := part.VMHost.RunSmartContractCall(typedRequest.CallInput)
	return common.NewMessageContractResponse(vmOutput, err)
}

func (part *VMPart) replyToDiagnoseWait(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageDiagnoseWaitRequest)
	duration := time.Duration(int64(typedRequest.Milliseconds) * int64(time.Millisecond))
	time.Sleep(duration)
	return common.NewMessageDiagnoseWaitResponse()
}

func (part *VMPart) replyToVersionRequest(_ common.MessageHandler) common.MessageHandler {
	return common.NewMessageVersionResponse(part.Version)
}

func (part *VMPart) replyToGasScheduleChange(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageGasScheduleChangeRequest)
	part.VMHost.GasScheduleChange(typedRequest.GasSchedule)
	return common.NewGasScheduleChangeResponse()
}
