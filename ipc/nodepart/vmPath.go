package nodepart

import (
	"os"
	"path"

	"github.com/multiversx/mx-chain-vm-go-v1_2/ipc/common"
)

func (driver *VMDriver) getVMPath() (string, error) {
	wasmvmPath, err := driver.getVMPathInCurrentDirectory()
	if err == nil {
		return wasmvmPath, nil
	}

	wasmvmPath, err = driver.getVMPathFromEnvironment()
	if err == nil {
		return wasmvmPath, nil
	}

	return "", common.ErrVMNotFound
}

func (driver *VMDriver) getVMPathInCurrentDirectory() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	wasmvmPath := path.Join(cwd, "wasmvm")
	if fileExists(wasmvmPath) {
		return wasmvmPath, nil
	}

	return "", common.ErrVMNotFound
}

func (driver *VMDriver) getVMPathFromEnvironment() (string, error) {
	wasmvmPath := os.Getenv(common.EnvVarWASMVMPath)
	if fileExists(wasmvmPath) {
		return wasmvmPath, nil
	}

	return "", common.ErrVMNotFound
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
