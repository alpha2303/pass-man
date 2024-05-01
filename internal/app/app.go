package app

import (
	"fmt"

	vault "github.com/alpha2303/pass-man/internal/app/vault"
	extras "github.com/alpha2303/pass-man/internal/pkg/extras"
)

func HandleVaultCreate() {
	vaultName, password := requestCredentials()

	vaultMeta := vault.CreateVaultMeta(vaultName)
	if vaultMeta.IsExists() {
		fmt.Printf("Vault named '%s' already exists.", vaultName)
		return
	}

	_, err := vault.CreateVault(&vaultMeta, password)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	fmt.Printf("Vault '%s' has been successfully created.", vaultName)
}

func HandleVaultSignIn() {
	vaultName, password := requestCredentials()

	vaultMeta := vault.CreateVaultMeta(vaultName)

	if !vaultMeta.IsExists() {
		fmt.Printf("Vault named '%s' does not exist.", vaultName)
		return
	}

	vaultObj, err := vault.OpenVault(&vaultMeta, password)

	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}

	if err = vault.Explore(vaultObj, password); err != nil {
		fmt.Printf("%s", err.Error())
	}
}

func HandleVaultDelete() {
	vaultName, password := requestCredentials()

	vaultMeta := vault.CreateVaultMeta(vaultName)

	if !vaultMeta.IsExists() {
		fmt.Printf("Vault named '%s' does not exist.", vaultName)
		return
	}

	vaultObj, err := vault.OpenVault(&vaultMeta, password)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}

	choice := extras.Input("Are you sure you want to delete this vault [Y/n]? : ")

	if choice == "Y" {
		vaultObj.Delete()
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
	}
}

func requestCredentials() (string, string) {
	vaultName := extras.Input("Enter vault name: ")
	password := extras.SilentInput("Enter password: ")

	return vaultName, password
}
