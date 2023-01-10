package host

import (
	"github.com/multiversx/mx-chain-vm-go-v1_2/wasmvm"
)

func (host *vmHost) handleBreakpointIfAny(executionErr error) error {
	if executionErr == nil {
		return nil
	}

	runtime := host.Runtime()
	breakpointValue := runtime.GetRuntimeBreakpointValue()
	if breakpointValue != wasmvm.BreakpointNone {
		executionErr = host.handleBreakpoint(breakpointValue)
	}

	return executionErr
}

func (host *vmHost) handleBreakpoint(breakpointValue wasmvm.BreakpointValue) error {
	if breakpointValue == wasmvm.BreakpointAsyncCall {
		return host.handleAsyncCallBreakpoint()
	}
	if breakpointValue == wasmvm.BreakpointExecutionFailed {
		return wasmvm.ErrExecutionFailed
	}
	if breakpointValue == wasmvm.BreakpointSignalError {
		return wasmvm.ErrSignalError
	}
	if breakpointValue == wasmvm.BreakpointOutOfGas {
		return wasmvm.ErrNotEnoughGas
	}

	return wasmvm.ErrUnhandledRuntimeBreakpoint
}
