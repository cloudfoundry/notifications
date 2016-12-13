package helpers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/dgrijalva/jwt-go"
)

var (
	UAAPrivateKey string
	UAAPublicKey  string
)

func init() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)

	if err != nil {
		panic(err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(privateKey.Public())

	privatePem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	publicPem := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	UAAPrivateKey = string(privatePem)
	UAAPublicKey = string(publicPem)
}

func BuildToken(header map[string]interface{}, claims map[string]interface{}) string {
	return BuildTokenWithKey(header, claims, UAAPrivateKey)
}

func BuildTokenWithKey(header map[string]interface{}, claims map[string]interface{}, signingKey string) string {
	application.UAAPublicKey = UAAPublicKey

	alg := header["alg"].(string)
	signingMethod := jwt.GetSigningMethod(alg)
	token := jwt.New(signingMethod)
	token.Header = header
	token.Claims = claims

	signed, err := token.SignedString([]byte(signingKey))
	if err != nil {
		panic(err)
	}

	return signed
}
