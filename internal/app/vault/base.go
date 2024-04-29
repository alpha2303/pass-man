package vault

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alpha2303/pass-man/internal/pkg/extras"
)

var (
	DefaultPath string = "./.vault"
)

func GetVaultCount() (int, error) {
	files, err := os.ReadDir(DefaultPath)

	if err != nil {
		return 0, err
	}

	return len(files), nil
}

func CreateVaultMeta(name string) VaultMeta {
	path := filepath.Join(DefaultPath, fmt.Sprintf("%s.txt", name))
	return VaultMeta{
		name: name,
		path: path,
	}
}

func CreateVault(vm VaultMeta, password string) (bool, error) {
	vault := Vault{
		name:        vm.name,
		credentials: nil,
	}

	if _, err := vault.Create(); err != nil {
		return false, err
	}

	return true, nil
}

func OpenVault(vm VaultMeta, password string) (*Vault, error) {
	vault := Vault{
		name:        vm.name,
		credentials: nil,
	}

	vault, err := vault.Open(password)
	if err != nil {
		return nil, err
	}

	if _, err = explore(vault); err != nil {
		return nil, err
	}

	return &vault, nil
}

func explore(v Vault) (bool, error) {
	var choice int = 1
	// todo
	for choice != 0 {
		fmt.Printf("What would you like to do?")
		fmt.Println("1. Add New Credential")
		fmt.Println("2. Remove Credential")

		extras.Input("Enter your choice: ", &choice)

		switch choice {
		case 1:
			handleCreateCredentials(&v)
		case 2:
			handleRemoveCredentials(&v)
		default:
			choice = 0
		}
		v.save()
	}
	return true, nil
}

func handleCreateCredentials(v *Vault) {
	var name, username, password string

	extras.Input("Name this credential: ", &name)
	fmt.Println("Enter credentials -")
	extras.Input("Username: ", &username)
	extras.SilentInput("Password: ", &password)

	credentials := CreateCredentials(username, password)

	v.add(name, credentials)
}

func handleRemoveCredentials(v *Vault) {
	var name string
	extras.Input("Enter Credential Name: ", &name)
	v.remove(name)
}
