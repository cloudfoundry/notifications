package helpers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/golang-jwt/jwt/v5"
)

var (
	UAAPrivateKey   string
	UAAPublicKey    string
	UAAPublicKeyRSA *rsa.PublicKey
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

	UAAPublicKeyRSA, err = jwt.ParseRSAPublicKeyFromPEM(publicPem)
	if err != nil {
		panic(err)
	}
}

func BuildToken(header map[string]interface{}, claims map[string]interface{}) string {
	return BuildTokenWithKey(header, claims, UAAPrivateKey)
}

func BuildTokenWithKey(header map[string]interface{}, claims map[string]interface{}, signingKey string) string {
	alg := header["alg"].(string)
	signingMethod := jwt.GetSigningMethod(alg)
	token := jwt.New(signingMethod)
	token.Header = header
	jwtClaims := jwt.MapClaims{}
	for i, j := range claims {
		jwtClaims[i] = j
	}
	token.Claims = jwtClaims
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(signingKey))
	if err != nil {
		panic(err)
	}
	signed, err := token.SignedString(key)
	if err != nil {
		panic(err)
	}

	return signed
}

func BuildHSATokenWithKey(header map[string]interface{}, claims map[string]interface{}, signingKey string) string {
	alg := header["alg"].(string)
	signingMethod := jwt.GetSigningMethod(alg)
	token := jwt.New(signingMethod)
	token.Header = header
	jwtClaims := jwt.MapClaims{}
	for i, j := range claims {
		jwtClaims[i] = j
	}
	token.Claims = jwtClaims
	signed, err := token.SignedString([]byte(signingKey))
	if err != nil {
		panic(err)
	}

	return signed
}
