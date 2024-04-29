package vault

import (
	"os"
)

type VaultMeta struct {
	name string
	path string
}

func (vm *VaultMeta) GetPath() string {
	return vm.path
}

func (vm *VaultMeta) IsExists() bool {
	_, err := os.Stat(vm.path)
	if err != nil {
		return false
	}
	return true
}
