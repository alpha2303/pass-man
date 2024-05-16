package vault

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	DefaultPath    string      = "./.vault"
	DefaultFileExt string      = "pmv"
	VaultPerm      os.FileMode = os.FileMode(int(0755))
)

func GetVaultCount() (int, error) {
	files, err := os.ReadDir(DefaultPath)

	if err != nil {
		return 0, err
	}

	return len(files), nil
}

func CreateVaultMeta(name *string) VaultMeta {
	if _, err := os.Stat(DefaultPath); err != nil {
		if err := os.MkdirAll(DefaultPath, VaultPerm); err != nil {
			fmt.Println(err.Error())
		}
	}

	path, err := filepath.Abs(filepath.Join(DefaultPath, fmt.Sprintf("%s.%s", *name, DefaultFileExt)))
	if err != nil {
		fmt.Println(err.Error())
	}
	return VaultMeta{
		name: name,
		path: &path,
	}
}

func CreateVault(vm *VaultMeta, password *string) error {

	vault := Vault{
		meta:        vm,
		isSignedIn:  false,
		credentials: nil,
	}

	if err := vault.Create(password); err != nil {
		return err
	}

	return nil
}

func OpenVault(vm *VaultMeta, password *string) (*Vault, error) {
	vault := Vault{
		meta:        vm,
		isSignedIn:  false,
		credentials: nil,
	}

	if err := vault.Open(password); err != nil {
		return nil, err
	}

	return &vault, nil
}
