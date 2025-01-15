package app

import (
	"errors"
	"fmt"
	"strconv"

	vault "github.com/alpha2303/pass-man/internal/app/vault"
	extras "github.com/alpha2303/pass-man/internal/pkg/extras"
)

func HandleVaultCreate() {
	vaultName, password := requestCredentials()

	vaultMeta := vault.CreateVaultMeta(&vaultName)
	if vaultMeta.IsExists() {
		fmt.Printf("Vault named '%s' already exists.", vaultName)
		return
	}

	if err := vault.CreateVault(&vaultMeta, &password); err != nil {
		fmt.Printf("%s", err.Error())
		return
	}

	fmt.Printf("Vault '%s' has been successfully created.", vaultName)
}

func HandleVaultSignIn() {
	vaultName, password := requestCredentials()

	vaultMeta := vault.CreateVaultMeta(&vaultName)

	if !vaultMeta.IsExists() {
		fmt.Printf("Vault named '%s' does not exist.", vaultName)
		return
	}

	vaultObj, err := vault.OpenVault(&vaultMeta, &password)

	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}

	fmt.Printf("\nSigned in to vault %s\n", vaultName)

	if err = explore(vaultObj, password); err != nil {
		fmt.Printf("%s", err.Error())
	}
}

func HandleVaultDelete() {
	vaultName := extras.Input("Enter vault name: ")
	vaultMeta := vault.CreateVaultMeta(&vaultName)

	if !vaultMeta.IsExists() {
		fmt.Printf("Vault named '%s' does not exist.", vaultName)
		return
	}

	password := extras.SilentInput("Enter password: ")

	vaultObj, err := vault.OpenVault(&vaultMeta, &password)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}

	choice := extras.Input("Are you sure you want to delete this vault [Y/n]? : ")

	if choice == "Y" {
		err := vaultObj.Delete()
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
	}
}

func explore(v *vault.Vault, password string) error {
	if !v.IsSignedIn() {
		return vault.ErrAuthAccess
	}

	var choice int = 1

	for choice != 0 {
		fmt.Println("\nWhat would you like to do?")
		fmt.Println("1. Add New Credential")
		fmt.Println("2. Remove Credential")
		fmt.Println("3. See Credential")
		fmt.Println("4. Display Credential Names")
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
		case 4:
			handleSeeAllCredNames(v)
		default:
			choice = 0
			return nil
		}

		if err := v.Save(&password); err != nil {
			return err
		}
	}
	return nil
}

func handleCreateCredentials(v *vault.Vault) {
	name := extras.Input("\nName this credential: ")
	fmt.Println("** Enter credentials **")
	domain := extras.Input("Domain: ")
	username := extras.Input("Username: ")
	password := extras.SilentInput("Password: ")

	credential := vault.CreateCredential(domain, username, password)

	var choice string = "Y"
	if _, error := v.GetCredential(name); !errors.Is(error, vault.ErrCredNotExist) {
		choice = extras.Input("A credential with this already exists.\nWould you like to replace it [Y/n]? : ")
	}

	if choice == "Y" {
		if err := v.AddCredential(name, credential); err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("Credential '%s' has been created.", name)
		}
		return
	}

	fmt.Println("Cancelling operation, please create credential with a unique name.")
}

func handleRemoveCredentials(v *vault.Vault) {
	name := extras.Input("Enter Credential Name: ")

	if _, error := v.GetCredential(name); errors.Is(error, vault.ErrCredNotExist) {
		fmt.Println("A credential with the provided name does not exist.")
		return
	}

	choice := extras.Input("\nThis credential will be permanently deleted. Are you sure [Y/n]? : ")
	if choice == "Y" {
		if err := v.RemoveCredential(name); err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("Credential '%s' has been deleted.", name)
		}
		return
	}

	fmt.Println("Removal operation cancelled.")
}

func handleSeeCredential(v *vault.Vault) {
	name := extras.Input("Enter Credential Name: ")

	cred, err := v.GetCredential(name)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Username: %s \n", cred.Username)
	fmt.Printf("Password: %s \n", cred.Password)
}

func handleSeeAllCredNames(v *vault.Vault) {
	credentials, err := v.GetAllCredentials()
	count := 1
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for key := range *credentials {
		fmt.Printf("%d. %s\n", count, key)
		count++
	}
}

func requestCredentials() (string, string) {
	vaultName := extras.Input("Enter vault name: ")
	password := extras.SilentInput("Enter password: ")

	return vaultName, password
}
