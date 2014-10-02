package cryptography

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "io"
)

type CryptoInterface interface {
    Encrypt(string) (string, error)
    Decrypt(string) (string, error)
}

type InvalidKeyError struct{}

func (err InvalidKeyError) Error() string {
    return "The key needs to be 16 bytes"
}

type CipherLengthError struct{}

func (err CipherLengthError) Error() string {
    return "Data length is too short"
}

type URLCryptoClient struct {
    key         []byte
    cipherBlock cipher.Block
}

func NewURLCryptoClient(key string) (URLCryptoClient, error) {
    if len([]byte(key)) != 16 {
        return URLCryptoClient{}, InvalidKeyError{}

    }

    cipherBlock, err := aes.NewCipher([]byte(key))

    if err != nil {
        return URLCryptoClient{}, err
    }

    return URLCryptoClient{
        key:         []byte(key),
        cipherBlock: cipherBlock,
    }, nil
}

func (crypto URLCryptoClient) Encrypt(data string) (string, error) {
    encodedText := base64.StdEncoding.EncodeToString([]byte(data))
    cipherText := make([]byte, aes.BlockSize+len(encodedText))

    initializationVector := cipherText[:aes.BlockSize]
    _, err := io.ReadFull(rand.Reader, initializationVector)

    if err != nil {
        return "", err
    }

    cipherEncrypter := cipher.NewCFBEncrypter(crypto.cipherBlock, initializationVector)
    cipherEncrypter.XORKeyStream(cipherText[aes.BlockSize:], []byte(encodedText))

    urlSafeCipherText := base64.URLEncoding.EncodeToString(cipherText)
    if err != nil {
        return "", err
    }

    return string(urlSafeCipherText), nil
}

func (crypto URLCryptoClient) Decrypt(data string) (string, error) {
    decodedData, err := base64.URLEncoding.DecodeString(data)
    if err != nil {
        return "", err
    }

    byteData := []byte(decodedData)

    if len(byteData) < aes.BlockSize {
        return "", CipherLengthError{}
    }

    initializationVector := byteData[:aes.BlockSize]
    byteData = byteData[aes.BlockSize:]

    cipherDecrypter := cipher.NewCFBDecrypter(crypto.cipherBlock, initializationVector)
    cipherDecrypter.XORKeyStream(byteData, byteData)

    decoded, err := base64.StdEncoding.DecodeString(string(byteData))
    if err != nil {
        return "", err
    }

    return string(decoded), nil
}
