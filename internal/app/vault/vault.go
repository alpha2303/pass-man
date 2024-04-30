package vault

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/alpha2303/pass-man/internal/pkg/connector"
	"github.com/alpha2303/pass-man/internal/pkg/crypto"
)

var (
	ErrVaultNotExist = errors.New("vault does not exist")
)

type Vault struct {
	meta        VaultMeta
	Credentials map[string]Credential
}

func (v *Vault) Create(password string) error {
	fileConnector, err := connector.GetConnector("file", v.meta.GetPath())
	if err != nil {
		return err
	}

	if err = fileConnector.Create(); err != nil {
		return err
	}

	if err = v.save(password, fileConnector, crypto.DefaultParams); err != nil {
		return err
	}

	return nil
}

func (v *Vault) Open(password string) error {
	fileConnector, err := connector.GetConnector("file", v.meta.GetPath())
	if err != nil {
		return err
	}

	encryptedBytes, err := fileConnector.Read()

	if err != nil {
		return err
	}

	if len(encryptedBytes) == 0 {
		v.Credentials = make(map[string]Credential, 0)
	}

	vaultBytes, err := crypto.Decrypt(string(encryptedBytes[:]), []byte(password))
	if err != nil {
		return err
	}

	if err := v.fromBinary(vaultBytes); err != nil {
		return err
	}

	return nil
}

func (v *Vault) Delete() error {
	if !v.meta.IsExists() {
		return ErrVaultNotExist
	}

	fileConnector, err := connector.GetConnector("file", v.meta.GetPath())
	if err != nil {
		return err
	}

	if err := fileConnector.Delete(); err != nil {
		return err
	}

	return nil
}

func (v *Vault) Save(password string) error {
	if !v.meta.IsExists() {
		return ErrVaultNotExist
	}

	fileConnector, err := connector.GetConnector("file", v.meta.GetPath())
	if err != nil {
		return err
	}

	if err = v.save(password, fileConnector, crypto.DefaultParams); err != nil {
		return nil
	}

	return nil
}

func (v *Vault) add(name string, credential Credential) {
	v.Credentials[name] = credential
}

func (v *Vault) remove(name string) {
	delete(v.Credentials, name)
}

func (v *Vault) toBinary() ([]byte, error) {
	var vaultBytes bytes.Buffer
	encoder := gob.NewEncoder(&vaultBytes)

	if err := encoder.Encode(*v); err != nil {
		return nil, err
	}

	return vaultBytes.Bytes(), nil
}

func (v *Vault) fromBinary(vaultBytes []byte) error {
	readBuffer := bytes.NewBuffer(vaultBytes)
	decoder := gob.NewDecoder(readBuffer)

	if err := decoder.Decode(v); err != nil {
		return err
	}

	return nil
}

func (v *Vault) save(password string, connector connector.IConnector, params *crypto.CryptoParams) error {
	vaultBytes, err := v.toBinary()
	if err != nil {
		return err
	}

	encryptedVaultString, err := crypto.Encrypt(vaultBytes, []byte(password), params)
	if err != nil {
		return err
	}

	if err = connector.Update(encryptedVaultString); err != nil {
		return err
	}

	return nil
}
