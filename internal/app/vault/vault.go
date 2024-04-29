package vault

import (
	"fmt"

	"github.com/alpha2303/pass-man/internal/pkg/extras"
)

type Vault struct {
	name        string
	credentials map[string]Credentials
}

func (v *Vault) Create() (bool, error) {}

func (v *Vault) Open(password string) (Vault, error) {}

func (v *Vault) Delete() (bool, error) {}

func (v *Vault) load() (Vault, error) {}

func (v *Vault) save() {}

func (v *Vault) add(name string, credential Credentials) {
	var choice string = "Y"
	if _, ok := v.credentials[name]; ok == true {
		extras.Input("A credential with this already exists.\nWould you like to replace it? [Y/n]", &choice)
	}

	if choice == "Y" {
		v.credentials[name] = credential
	} else {
		fmt.Println("Cancelling operation, please create credential with a unique name.")
	}
}

func (v *Vault) remove(name string) {
	var choice string
	if _, ok := v.credentials[name]; ok != true {
		fmt.Println("A credential with the provided name does not exist.")
	}
	extras.Input("This credential will be permanently deleted. Are you sure? [Y/n]", &choice)
	if choice == "Y" {
		delete(v.credentials, name)
		fmt.Printf("Credential '%s' has been deleted.", name)
	} else {
		fmt.Println("Removal operation cancelled.")
	}
}

func (v *Vault) toBinary() {}
