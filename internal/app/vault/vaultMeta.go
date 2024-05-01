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
	if _, err := os.Stat(vm.path); err != nil {
		return false
	}

	return true
}
