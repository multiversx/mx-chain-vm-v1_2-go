package main

import (
	"fmt"
	"os"

	"github.com/ElrondNetwork/elrond-go-logger/pipes"
	"github.com/ElrondNetwork/wasm-vm-v1_2/ipc/common"
	"github.com/ElrondNetwork/wasm-vm-v1_2/ipc/marshaling"
	"github.com/ElrondNetwork/wasm-vm-v1_2/ipc/wasmvmpart"
	"github.com/ElrondNetwork/wasm-vm-v1_2/wasmvm"
)

const (
	fileDescriptorArwenInit      = 3
	fileDescriptorNodeToArwen    = 4
	fileDescriptorArwenToNode    = 5
	fileDescriptorReadLogProfile = 6
	fileDescriptorLogToNode      = 7
)

var appVersion = "undefined"

func main() {
	errCode, errMessage := doMain()
	if errCode != common.ErrCodeSuccess {
		fmt.Fprintln(os.Stderr, errMessage)
		os.Exit(errCode)
	}
}

// doMain returns (error code, error message)
func doMain() (int, string) {
	wasmvmInitFile := getPipeFile(fileDescriptorArwenInit)
	if wasmvmInitFile == nil {
		return common.ErrCodeCannotCreateFile, "Cannot get pipe file: [wasmvmInitFile]"
	}

	nodeToArwenFile := getPipeFile(fileDescriptorNodeToArwen)
	if nodeToArwenFile == nil {
		return common.ErrCodeCannotCreateFile, "Cannot get pipe file: [nodeToArwenFile]"
	}

	wasmvmToNodeFile := getPipeFile(fileDescriptorArwenToNode)
	if wasmvmToNodeFile == nil {
		return common.ErrCodeCannotCreateFile, "Cannot get pipe file: [wasmvmToNodeFile]"
	}

	readLogProfileFile := getPipeFile(fileDescriptorReadLogProfile)
	if readLogProfileFile == nil {
		return common.ErrCodeCannotCreateFile, "Cannot get pipe file: [readLogProfileFile]"
	}

	logToNodeFile := getPipeFile(fileDescriptorLogToNode)
	if logToNodeFile == nil {
		return common.ErrCodeCannotCreateFile, "Cannot get pipe file: [logToNodeFile]"
	}

	wasmvmArguments, err := common.GetArwenArguments(wasmvmInitFile)
	if err != nil {
		return common.ErrCodeInit, fmt.Sprintf("Cannot receive gasSchedule: %v", err)
	}

	messagesMarshalizer := marshaling.CreateMarshalizer(wasmvmArguments.MessagesMarshalizer)
	logsMarshalizer := marshaling.CreateMarshalizer(wasmvmArguments.LogsMarshalizer)

	logsPart, err := pipes.NewChildPart(readLogProfileFile, logToNodeFile, logsMarshalizer)
	if err != nil {
		return common.ErrCodeInit, fmt.Sprintf("Cannot create logs part: %v", err)
	}

	err = logsPart.StartLoop()
	if err != nil {
		return common.ErrCodeInit, fmt.Sprintf("Cannot start logs loop: %v", err)
	}

	defer logsPart.StopLoop()

	appVersion = wasmvm.WASMVMVersion
	part, err := wasmvmpart.NewArwenPart(
		appVersion,
		nodeToArwenFile,
		wasmvmToNodeFile,
		&wasmvmArguments.VMHostParameters,
		messagesMarshalizer,
	)
	if err != nil {
		return common.ErrCodeInit, fmt.Sprintf("Cannot create ArwenPart: %v", err)
	}

	err = part.StartLoop()
	if err != nil {
		return common.ErrCodeTerminated, fmt.Sprintf("Ended Arwen loop: %v", err)
	}

	// This is never reached, actually. Arwen is supposed to run an infinite message loop.
	return common.ErrCodeSuccess, ""
}

func getPipeFile(fileDescriptor uintptr) *os.File {
	file := os.NewFile(fileDescriptor, fmt.Sprintf("/proc/self/fd/%d", fileDescriptor))
	return file
}
