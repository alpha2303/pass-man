package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	argon2 "golang.org/x/crypto/argon2"
)

type CryptoParams struct {
	memory     uint32
	iterations uint32
	threads    uint8
	saltLength uint32
	keyLength  uint32
}

var (
	DefaultParams = &CryptoParams{
		memory:     64 * 1024,
		iterations: 1,
		threads:    4,
		saltLength: 16,
		keyLength:  32,
	}
	ErrInvalidData         = errors.New("data provided is not a valid encoded data")
	ErrIncompatibleVersion = errors.New("encryption version is incompatible")
)

func Encrypt(srcBytes []byte, masterPassword []byte, params *CryptoParams) (string, error) {
	salt, err := generateRandomBytes(params.saltLength)

	if err != nil {
		return "", err
	}

	key, err := generateKey(masterPassword, salt, params)

	if err != nil {
		return "", err
	}

	cipherBlock, err := aes.NewCipher(key)

	if err != nil {
		return "", err
	}

	aesGCMCipher, err := cipher.NewGCM(cipherBlock)

	if err != nil {
		return "", err
	}

	nonce, err := generateNonce(aesGCMCipher)

	if err != nil {
		return "", err
	}

	var destBytes []byte = aesGCMCipher.Seal(nonce, nonce, srcBytes, nil)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Data := base64.RawStdEncoding.EncodeToString(destBytes)

	encodedData := fmt.Sprintf("$argon2i$%d$%s$%s$%d,%d,%d", argon2.Version, b64Data, b64Salt, params.iterations, params.memory, params.threads)

	return encodedData, nil
}

func Decrypt(encodedData string, masterPassword []byte) ([]byte, error) {

	params, encryptedData, salt, err := decodeData(encodedData)

	if err != nil {
		return nil, err
	}

	key, err := generateKey(masterPassword, salt, params)

	if err != nil {
		return nil, err
	}

	cipherBlock, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	aesGCMCipher, err := cipher.NewGCM(cipherBlock)

	if err != nil {
		return nil, err
	}

	nonce, cipherText := encryptedData[:aesGCMCipher.NonceSize()], encryptedData[aesGCMCipher.NonceSize():]

	destBytes, err := aesGCMCipher.Open(nil, nonce, cipherText, nil)

	if err != nil {
		return nil, err
	}

	return destBytes, nil
}

func decodeData(encodedData string) (*CryptoParams, []byte, []byte, error) {
	values := strings.Split(encodedData, "$")

	if len(values) != 6 {
		return nil, nil, nil, ErrInvalidData
	}

	var version int
	if _, err := fmt.Sscanf(values[2], "%d", &version); err != nil {
		return nil, nil, nil, err
	}

	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	data, err := base64.RawStdEncoding.DecodeString(values[3])

	if err != nil {
		return nil, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(values[4])

	if err != nil {
		return nil, nil, nil, err
	}

	params := &CryptoParams{}
	if _, err := fmt.Sscanf(values[5], "%d,%d,%d", &(params.iterations), &(params.memory), &(params.threads)); err != nil {
		return nil, nil, nil, err
	}

	params.saltLength = uint32(len(salt))
	params.keyLength = uint32(32)

	return params, data, salt, nil
}

func generateKey(masterPassword []byte, salt []byte, params *CryptoParams) ([]byte, error) {

	return argon2.IDKey(
		masterPassword,
		salt,
		params.iterations,
		params.memory,
		params.threads,
		params.keyLength,
	), nil
}

func generateRandomBytes(bytesLength uint32) ([]byte, error) {
	var bytes []byte = make([]byte, bytesLength)

	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}

	return bytes, nil
}

func generateNonce(cipherAEAD cipher.AEAD) ([]byte, error) {
	var nonce []byte = make([]byte, cipherAEAD.NonceSize())

	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	return nonce, nil
}
