package main

import (
	"os"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/wasm-vm-v1_2/wasmvmdebug"
)

var log = logger.GetOrCreate("wasmvmdebug")

const (
	// ErrCodeSuccess signals success
	ErrCodeSuccess = iota
	// ErrCodeCriticalError signals a critical error
	ErrCodeCriticalError
)

func main() {
	logger.ToggleLoggerName(true)
	_ = logger.SetLogLevel("*:TRACE")

	facade := wasmvmdebug.NewDebugFacade()
	app := initializeCLI(facade)

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(ErrCodeCriticalError)
	}

	os.Exit(ErrCodeSuccess)
}
