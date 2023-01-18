package nodepart

import (
	"os"
	"path"

	"github.com/multiversx/mx-chain-vm-v1_2-go/ipc/common"
)

func (driver *VMDriver) getVMPath() (string, error) {
	vmPath, err := driver.getVMPathInCurrentDirectory()
	if err == nil {
		return vmPath, nil
	}

	vmPath, err = driver.getVMPathFromEnvironment()
	if err == nil {
		return vmPath, nil
	}

	return "", common.ErrVMNotFound
}

func (driver *VMDriver) getVMPathInCurrentDirectory() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	vmPath := path.Join(cwd, "vm")
	if fileExists(vmPath) {
		return vmPath, nil
	}

	return "", common.ErrVMNotFound
}

func (driver *VMDriver) getVMPathFromEnvironment() (string, error) {
	vmPath := os.Getenv(common.EnvVarVMPath)
	if fileExists(vmPath) {
		return vmPath, nil
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
