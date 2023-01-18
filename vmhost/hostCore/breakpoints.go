package hostCore

import (
	"github.com/multiversx/mx-chain-vm-v1_2-go/vmhost"
)

func (host *vmHost) handleBreakpointIfAny(executionErr error) error {
	if executionErr == nil {
		return nil
	}

	runtime := host.Runtime()
	breakpointValue := runtime.GetRuntimeBreakpointValue()
	if breakpointValue != vmhost.BreakpointNone {
		executionErr = host.handleBreakpoint(breakpointValue)
	}

	return executionErr
}

func (host *vmHost) handleBreakpoint(breakpointValue vmhost.BreakpointValue) error {
	if breakpointValue == vmhost.BreakpointAsyncCall {
		return host.handleAsyncCallBreakpoint()
	}
	if breakpointValue == vmhost.BreakpointExecutionFailed {
		return vmhost.ErrExecutionFailed
	}
	if breakpointValue == vmhost.BreakpointSignalError {
		return vmhost.ErrSignalError
	}
	if breakpointValue == vmhost.BreakpointOutOfGas {
		return vmhost.ErrNotEnoughGas
	}

	return vmhost.ErrUnhandledRuntimeBreakpoint
}
