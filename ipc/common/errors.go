package common

import (
	"fmt"
)

// CriticalError signals a critical error
type CriticalError struct {
	InnerErr error
}

// WrapCriticalError wraps an error
func WrapCriticalError(err error) *CriticalError {
	return &CriticalError{InnerErr: err}
}

func (err *CriticalError) Error() string {
	return fmt.Sprintf("critical error: %v", err.InnerErr)
}

// Unwrap unwraps the inner error
func (err *CriticalError) Unwrap() error {
	return err.InnerErr
}

// IsCriticalError returns whether the error is critical
func IsCriticalError(err error) bool {
	_, ok := err.(*CriticalError)
	return ok
}

// ErrBadVMArguments signals a critical error
var ErrBadVMArguments = &CriticalError{InnerErr: fmt.Errorf("bad arguments passed to vm")}

// ErrVMClosed signals a critical error
var ErrVMClosed = &CriticalError{InnerErr: fmt.Errorf("vm closed")}

// ErrVMTimeExpired signals a critical error
var ErrVMTimeExpired = &CriticalError{InnerErr: fmt.Errorf("vm time expired")}

// ErrVMNotFound signals a critical error
var ErrVMNotFound = &CriticalError{InnerErr: fmt.Errorf("vm binary not found")}

// ErrInvalidMessageNonce signals a critical error
var ErrInvalidMessageNonce = &CriticalError{InnerErr: fmt.Errorf("invalid dialogue nonce")}

// ErrStopPerNodeRequest signals a critical error
var ErrStopPerNodeRequest = &CriticalError{InnerErr: fmt.Errorf("vm will stop, as requested")}

// ErrBadRequestFromNode signals a critical error
var ErrBadRequestFromNode = &CriticalError{InnerErr: fmt.Errorf("bad message from node")}

// ErrBadMessageFromVM signals a critical error
var ErrBadMessageFromVM = &CriticalError{InnerErr: fmt.Errorf("bad message from vm")}

// ErrCannotSendContractRequest signals a critical error
var ErrCannotSendContractRequest = &CriticalError{InnerErr: fmt.Errorf("cannot send contract request")}

// ErrCannotSendHookCallResponse signals a critical error
var ErrCannotSendHookCallResponse = &CriticalError{InnerErr: fmt.Errorf("cannot send hook call response")}

// ErrCannotSendHookCallRequest signals a critical error
var ErrCannotSendHookCallRequest = &CriticalError{InnerErr: fmt.Errorf("cannot send hook call request")}

// ErrCannotReceiveHookCallResponse signals a critical error
var ErrCannotReceiveHookCallResponse = &CriticalError{InnerErr: fmt.Errorf("cannot receive hook call response")}

// ErrBadHookResponseFromNode signals a critical error
var ErrBadHookResponseFromNode = &CriticalError{InnerErr: fmt.Errorf("bad hook response from node")}

const (
	// ErrCodeSuccess signals success
	ErrCodeSuccess = iota
	// ErrCodeCannotCreateFile signals a critical error
	ErrCodeCannotCreateFile
	// ErrCodeInit signals a critical error
	ErrCodeInit
	// ErrCodeTerminated signals a critical error
	ErrCodeTerminated
)
