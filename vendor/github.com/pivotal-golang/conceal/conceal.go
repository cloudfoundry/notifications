// Package conceal provides the ability to encrypt/decrypt byte slices using aes encryption.
package conceal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
)

type CloakInterface interface {
	Veil([]byte) ([]byte, error)
	Unveil([]byte) ([]byte, error)
}

type CipherLengthError struct{}

// CipherLengthError occurs when the data passed in to Unveil is shorter than the length of the 16 byte encryption key.
func (err CipherLengthError) Error() string {
	return "Data length should be at least 16 bytes"
}

// A Cloak encrypts and decrypts []byte using Veil and Unveil.
type Cloak struct {
	cipherBlock cipher.Block
}

// NewCloak takes a key that is resized to 16 bytes and used as a key in AES encryption.
// It returns a Cloak. If the key cannot be used to create a cipherBlock, an error is returned.
func NewCloak(key []byte) (Cloak, error) {
	cipherBlock, err := aes.NewCipher(resizeKey(key))
	if err != nil {
		return Cloak{}, err
	}

	return Cloak{
		cipherBlock: cipherBlock,
	}, nil
}

func resizeKey(key []byte) []byte {
	resizedKey := md5.Sum(key)
	return resizedKey[:]
}

// Veil base64 encodes a slice of bytes and uses aes encryption. It returns an encrypted slice of bytes,
// and an error.
func (cloak Cloak) Veil(data []byte) ([]byte, error) {
	encodedData := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(encodedData, data)
	cipherText := make([]byte, aes.BlockSize+len(encodedData))

	initializationVector := cipherText[:aes.BlockSize]
	_, err := rand.Read(initializationVector)
	if err != nil {
		return []byte{}, err
	}

	cipherEncrypter := cipher.NewCFBEncrypter(cloak.cipherBlock, initializationVector)
	cipherEncrypter.XORKeyStream(cipherText[aes.BlockSize:], encodedData)

	base64CipherText := make([]byte, base64.URLEncoding.EncodedLen(len(cipherText)))
	base64.URLEncoding.Encode(base64CipherText, cipherText)
	return base64CipherText, nil
}

// Unveil base64 decodes a slice of bytes and uses aes encryption to decrypt. It returns a decrypted slice of bytes,
// and an error. A CipherLengthError is returned if the data is less than 16 bytes.
func (cloak Cloak) Unveil(data []byte) ([]byte, error) {
	decodedData := make([]byte, base64.URLEncoding.DecodedLen(len(data)))
	n, err := base64.URLEncoding.Decode(decodedData, data)
	if err != nil {
		return []byte{}, err
	}
	decodedData = decodedData[:n]

	if len(decodedData) < aes.BlockSize {
		return []byte{}, CipherLengthError{}
	}

	initializationVector := decodedData[:aes.BlockSize]
	decodedData = decodedData[aes.BlockSize:]

	cipherDecrypter := cipher.NewCFBDecrypter(cloak.cipherBlock, initializationVector)
	cipherDecrypter.XORKeyStream(decodedData, decodedData)

	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(decodedData)))
	n, err = base64.StdEncoding.Decode(decoded, decodedData)
	if err != nil {
		return []byte{}, err
	}

	return decoded[:n], nil
}
