package wasmvmpart

import (
	"os"

	"github.com/multiversx/wasm-vm-v1_2/ipc/common"
	"github.com/multiversx/wasm-vm-v1_2/ipc/marshaling"
)

// VMMessenger is the messenger on VM's part of the pipe
type VMMessenger struct {
	common.Messenger
}

// NewVMMessenger creates a new messenger
func NewVMMessenger(reader *os.File, writer *os.File, marshalizer marshaling.Marshalizer) *VMMessenger {
	return &VMMessenger{
		Messenger: *common.NewMessengerPipes("WASMVM", reader, writer, marshalizer),
	}
}

// ReceiveNodeRequest waits for a request from Node
func (messenger *VMMessenger) ReceiveNodeRequest() (common.MessageHandler, error) {
	message, err := messenger.Receive(0)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// SendContractResponse sends a contract response to the Node
func (messenger *VMMessenger) SendContractResponse(response common.MessageHandler) error {
	log.Trace("[WASMVM]: SendContractResponse", "response", response.DebugString())

	err := messenger.Send(response)
	if err != nil {
		return err
	}

	return nil
}

// SendHookCallRequest makes a hook call (over the pipe) and waits for the response
func (messenger *VMMessenger) SendHookCallRequest(request common.MessageHandler) (common.MessageHandler, error) {
	log.Trace("[WASMVM]: SendHookCallRequest", "request", request.DebugString())

	err := messenger.Send(request)
	if err != nil {
		return nil, common.ErrCannotSendHookCallRequest
	}

	response, err := messenger.Receive(0)
	if err != nil {
		return nil, common.ErrCannotReceiveHookCallResponse
	}

	return response, nil
}
