package nodepart

import (
	"os"
	"path"

	"github.com/ElrondNetwork/wasm-vm-v1_2/ipc/common"
)

func (driver *ArwenDriver) getArwenPath() (string, error) {
	wasmvmPath, err := driver.getArwenPathInCurrentDirectory()
	if err == nil {
		return wasmvmPath, nil
	}

	wasmvmPath, err = driver.getArwenPathFromEnvironment()
	if err == nil {
		return wasmvmPath, nil
	}

	return "", common.ErrArwenNotFound
}

func (driver *ArwenDriver) getArwenPathInCurrentDirectory() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	wasmvmPath := path.Join(cwd, "wasmvm")
	if fileExists(wasmvmPath) {
		return wasmvmPath, nil
	}

	return "", common.ErrArwenNotFound
}

func (driver *ArwenDriver) getArwenPathFromEnvironment() (string, error) {
	wasmvmPath := os.Getenv(common.EnvVarWASMVMPath)
	if fileExists(wasmvmPath) {
		return wasmvmPath, nil
	}

	return "", common.ErrArwenNotFound
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
