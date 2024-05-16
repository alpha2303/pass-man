package vault

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/alpha2303/pass-man/internal/pkg/connector"
	"github.com/alpha2303/pass-man/internal/pkg/crypto"
)

var (
	ErrVaultNotExist error = errors.New("vault does not exist")
	ErrAuthAccess    error = errors.New("unauthenticated vault access")
	ErrCredNotExist  error = errors.New("credential does not exist in vault")
)

type Vault struct {
	meta        *VaultMeta
	isSignedIn  bool
	credentials map[string]Credential
}

func (v *Vault) Create(password *string) error {
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

func (v *Vault) Open(password *string) error {
	fileConnector, err := connector.GetConnector("file", v.meta.GetPath())
	if err != nil {
		return err
	}

	encryptedBytes, err := fileConnector.Read()

	if err != nil {
		return err
	}

	if len(encryptedBytes) == 0 {
		v.credentials = make(map[string]Credential, 0)
	}

	vaultBytes, err := crypto.Decrypt(string(encryptedBytes[:]), []byte(*password))
	if err != nil {
		return err
	}

	if err := v.fromBinary(vaultBytes); err != nil {
		return err
	}

	v.isSignedIn = true

	return nil
}

func (v *Vault) IsSignedIn() bool {
	return v.isSignedIn
}

func (v *Vault) Delete() error {
	if !v.meta.IsExists() {
		return ErrVaultNotExist
	}

	if !v.IsSignedIn() {
		return ErrAuthAccess
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

func (v *Vault) Save(password *string) error {
	if !v.meta.IsExists() {
		return ErrVaultNotExist
	}

	if !v.IsSignedIn() {
		return ErrAuthAccess
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

func (v *Vault) GetCredential(name string) (*Credential, error) {
	if !v.meta.IsExists() {
		return nil, ErrVaultNotExist
	}

	if !v.IsSignedIn() {
		return nil, ErrAuthAccess
	}

	value, ok := v.getCredential(name)
	if !ok {
		return nil, ErrCredNotExist
	}

	return value, nil
}

func (v *Vault) GetAllCredentials() (*map[string]Credential, error) {
	if !v.meta.IsExists() {
		return nil, ErrVaultNotExist
	}

	if !v.IsSignedIn() {
		return nil, ErrAuthAccess
	}

	return v.getAllCredentials(), nil
}

func (v *Vault) AddCredential(name string, credential Credential) error {
	if !v.meta.IsExists() {
		return ErrVaultNotExist
	}

	if !v.IsSignedIn() {
		return ErrAuthAccess
	}

	if v.credentials == nil {
		v.credentials = make(map[string]Credential)
	}

	v.addCredential(name, credential)

	return nil
}

func (v *Vault) RemoveCredential(name string) error {
	if !v.meta.IsExists() {
		return ErrVaultNotExist
	}

	if !v.IsSignedIn() {
		return ErrAuthAccess
	}

	v.removeCredential(name)

	return nil
}

func (v *Vault) getCredential(name string) (*Credential, bool) {
	value, ok := v.credentials[name]
	return &value, ok
}

func (v *Vault) getAllCredentials() *map[string]Credential {
	return &v.credentials
}

func (v *Vault) addCredential(name string, credential Credential) {
	v.credentials[name] = credential
}

func (v *Vault) removeCredential(name string) {
	delete(v.credentials, name)
}

func (v *Vault) toBinary() ([]byte, error) {
	var vaultBytes bytes.Buffer
	encoder := gob.NewEncoder(&vaultBytes)

	if err := encoder.Encode(v.credentials); err != nil {
		return nil, err
	}

	return vaultBytes.Bytes(), nil
}

func (v *Vault) fromBinary(vaultBytes []byte) error {
	readBuffer := bytes.NewBuffer(vaultBytes)
	decoder := gob.NewDecoder(readBuffer)

	var credentials map[string]Credential
	if err := decoder.Decode(&credentials); err != nil {
		return err
	}

	v.credentials = credentials

	return nil
}

func (v *Vault) save(password *string, connector connector.IConnector, params *crypto.CryptoParams) error {
	vaultBytes, err := v.toBinary()
	if err != nil {
		return err
	}

	encryptedVaultString, err := crypto.Encrypt(vaultBytes, []byte(*password), params)
	if err != nil {
		return err
	}

	if err = connector.Update(encryptedVaultString); err != nil {
		return err
	}

	return nil
}
