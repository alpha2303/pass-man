package app

import (
	"fmt"

	vault "github.com/alpha2303/pass-man/internal/app/vault"
	extras "github.com/alpha2303/pass-man/internal/pkg/extras"
)

func HandleVaultCreate() {
	vaultName := extras.Input("Enter vault name: ")
	vaultMeta := vault.CreateVaultMeta(vaultName)

	if vaultMeta.IsExists() {
		fmt.Printf("Vault named '%s' already exists.", vaultName)
		return
	}

	var password string
	extras.SilentInput("Enter password for vault: ", &password)

	_, err := vault.CreateVault(&vaultMeta, password)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	fmt.Printf("Vault '%s' has been successfully created.", vaultName)
}

func HandleVaultSignIn() {
	vaultName := extras.Input("Enter vault name: ")
	vaultMeta := vault.CreateVaultMeta(vaultName)

	if !vaultMeta.IsExists() {
		fmt.Printf("Vault named '%s' does not exist.", vaultName)
		return
	}

	var password string
	extras.SilentInput("Enter password for vault: ", &password)

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
	vaultName := extras.Input("Enter vault name: ")
	vaultMeta := vault.CreateVaultMeta(vaultName)

	if !vaultMeta.IsExists() {
		fmt.Printf("Vault named '%s' does not exist.", vaultName)
		return
	}

	var password string
	extras.SilentInput("Enter password for vault: ", &password)

	vault_obj, err := vault.OpenVault(&vaultMeta, password)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}

	choice := extras.Input("Are you sure you want to delete this vault [Y/n]? : ")

	if choice == "Y" {
		vault_obj.Delete()
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
	}
}

func requestCredentials() (string, string) {
	var password string
	vaultName := extras.Input("Enter vault name: ")
	extras.SilentInput("Enter password: ", &password)

	return vaultName, password
}
