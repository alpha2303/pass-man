package vault

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/alpha2303/pass-man/internal/pkg/extras"
)

var (
	DefaultPath    string      = "./.vault"
	DefaultFileExt string      = "pmv"
	VaultPerm      os.FileMode = os.FileMode(int(0755))
	ErrAuthAccess  error       = errors.New("unauthenticated vault access")
)

func GetVaultCount() (int, error) {
	files, err := os.ReadDir(DefaultPath)

	if err != nil {
		return 0, err
	}

	return len(files), nil
}

func CreateVaultMeta(name string) VaultMeta {
	if _, err := os.Stat(DefaultPath); err != nil {
		if err := os.MkdirAll(DefaultPath, VaultPerm); err != nil {
			fmt.Println(err.Error())
		}
	}

	path, err := filepath.Abs(filepath.Join(DefaultPath, fmt.Sprintf("%s.%s", name, DefaultFileExt)))
	if err != nil {
		fmt.Println(err.Error())
	}
	return VaultMeta{
		name: name,
		path: path,
	}
}

func CreateVault(vm *VaultMeta, password string) (bool, error) {
	vault := Vault{
		meta:        *vm,
		credentials: nil,
	}

	if err := vault.Create(password); err != nil {
		return false, err
	}

	return true, nil
}

func OpenVault(vm *VaultMeta, password string) (*Vault, error) {
	vault := Vault{
		meta:        *vm,
		credentials: nil,
	}

	err := vault.Open(password)
	if err != nil {
		return nil, err
	}

	return &vault, nil
}

func Explore(v *Vault, password string) error {
	if !v.isSignedIn {
		return ErrAuthAccess
	}

	var choice int = 1

	for choice != 0 {
		fmt.Println("\nWhat would you like to do?")
		fmt.Println("1. Add New Credential")
		fmt.Println("2. Remove Credential")
		fmt.Println("3. See Credential")
		fmt.Println("0. Go Back")

		choice, err := strconv.ParseInt(extras.Input("\nEnter your choice: "), 10, 32)
		if err != nil {
			return err
		}

		switch choice {
		case 1:
			handleCreateCredentials(v)
		case 2:
			handleRemoveCredentials(v)
		case 3:
			handleSeeCredential(v)
		default:
			choice = 0
			return nil
		}

		if err := v.Save(password); err != nil {
			return err
		}
	}
	return nil
}

func handleCreateCredentials(v *Vault) {
	name := extras.Input("\nName this credential: ")
	fmt.Println("Enter credentials -")
	username := extras.Input("Username: ")
	password := extras.SilentInput("Password: ")

	credential := CreateCredential(username, password)

	if v.credentials == nil {
		v.credentials = make(map[string]Credential)
	}

	var choice string = "Y"
	if _, ok := v.credentials[name]; ok {
		choice = extras.Input("A credential with this already exists.\nWould you like to replace it [Y/n]? : ")
	}

	if choice == "Y" {
		v.add(name, credential)
	} else {
		fmt.Println("Cancelling operation, please create credential with a unique name.")
	}
}

func handleRemoveCredentials(v *Vault) {
	name := extras.Input("Enter Credential Name: ")

	if _, ok := v.credentials[name]; !ok {
		fmt.Println("A credential with the provided name does not exist.")
		return
	}

	choice := extras.Input("\nThis credential will be permanently deleted. Are you sure [Y/n]? : ")
	if choice == "Y" {
		v.remove(name)
		fmt.Printf("Credential '%s' has been deleted.", name)
	} else {
		fmt.Println("Removal operation cancelled.")
	}
}

func handleSeeCredential(v *Vault) {
	name := extras.Input("Enter Credential Name: ")

	cred, ok := v.credentials[name]

	if !ok {
		fmt.Println("A credential with the provided name does not exist.")
		return
	}

	fmt.Printf("Username: %s \n", cred.Username)
	fmt.Printf("Password: %s \n", cred.Password)
}
