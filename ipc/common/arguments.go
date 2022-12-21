package common

import (
	"os"

	"github.com/ElrondNetwork/wasm-vm-v1_2/ipc/marshaling"
	"github.com/ElrondNetwork/wasm-vm-v1_2/wasmvm"
)

// VMArguments represents the initialization arguments required by VM, passed through the initialization pipe
type VMArguments struct {
	wasmvm.VMHostParameters
	LogsMarshalizer     marshaling.MarshalizerKind
	MessagesMarshalizer marshaling.MarshalizerKind
}

// SendVMArguments sends initialization arguments through a pipe
func SendVMArguments(pipe *os.File, pipeArguments VMArguments) error {
	sender := NewSender(pipe, createArgumentsMarshalizer())
	message := NewMessageInitialize(pipeArguments)
	_, err := sender.Send(message)
	return err
}

// GetVMArguments reads initialization arguments from the pipe
func GetVMArguments(pipe *os.File) (*VMArguments, error) {
	receiver := NewReceiver(pipe, createArgumentsMarshalizer())
	message, _, err := receiver.Receive(0)
	if err != nil {
		return nil, err
	}

	typedMessage := message.(*MessageInitialize)
	return &typedMessage.Arguments, nil
}

// For the arguments, the marshalizer is fixed to JSON
func createArgumentsMarshalizer() marshaling.Marshalizer {
	return marshaling.CreateMarshalizer(marshaling.JSON)
}
