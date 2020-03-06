package nodepart

import (
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
)

// NodeMessenger is
type NodeMessenger struct {
	common.Messenger
}

// NewNodeMessenger creates
func NewNodeMessenger(reader *os.File, writer *os.File) *NodeMessenger {
	return &NodeMessenger{
		Messenger: *common.NewMessenger("NODE", reader, writer),
	}
}

// SendContractRequest sends
func (messenger *NodeMessenger) SendContractRequest(request common.MessageHandler) error {
	err := messenger.Send(request)
	if err != nil {
		return common.ErrCannotSendContractRequest
	}

	return nil
}

// SendHookCallResponse sends
func (messenger *NodeMessenger) SendHookCallResponse(response common.MessageHandler) error {
	err := messenger.Send(response)
	if err != nil {
		return common.ErrCannotSendHookCallResponse
	}

	return nil
}

// ReceiveHookCallRequestOrContractResponse waits
func (messenger *NodeMessenger) ReceiveHookCallRequestOrContractResponse(timeout int) (common.MessageHandler, error) {
	message, err := messenger.Receive(timeout)
	if err != nil {
		return nil, err
	}

	return message, nil
}
