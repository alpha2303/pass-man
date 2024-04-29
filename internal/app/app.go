package app

import (
	"fmt"

	vault "github.com/alpha2303/pass-man/internal/app/vault"
	extras "github.com/alpha2303/pass-man/internal/pkg/extras"
)

func HandleVaultCreate() {
	var vaultName string
	extras.Input("Enter vault name: ", &vaultName)
	vaultMeta := vault.CreateVaultMeta(vaultName)

	if vaultMeta.IsExists() {
		fmt.Printf("Vault named '%s' already exists.", vaultName)
		return
	}

	var password string
	extras.SilentInput("Enter password for vault: ", &password)

	_, err := vault.CreateVault(vaultMeta, password)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
}

func HandleVaultSignIn() {
	var vaultName string
	extras.Input("Enter vault name: ", &vaultName)
	vaultMeta := vault.CreateVaultMeta(vaultName)

	if !vaultMeta.IsExists() {
		fmt.Printf("Vault named '%s' does not exist.", vaultName)
		return
	}

	var password string
	extras.SilentInput("Enter password for vault: ", &password)

	if _, err := vault.OpenVault(vaultMeta, password); err != nil {
		fmt.Printf("%s", err.Error())
	}
}

func HandleVaultDelete() {
	var vaultName string
	extras.Input("Enter vault name: ", &vaultName)
	vaultMeta := vault.CreateVaultMeta(vaultName)

	if !vaultMeta.IsExists() {
		fmt.Printf("Vault named '%s' does not exist.", vaultName)
		return
	}

	var password string
	extras.SilentInput("Enter password for vault: ", &password)

	vault_obj, err := vault.OpenVault(vaultMeta, password)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}

	var choice string
	extras.Input("Are you sure you want to delete this vault? [Y/n]", &choice)

	if choice == "Y" {
		vault_obj.Delete()
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
	}
}

func requestCredentials() (string, string) {
	var vaultName, password string
	extras.Input("Enter vault name: ", &vaultName)
	extras.SilentInput("Enter password: ", &password)

	return vaultName, password
}
