package fakes

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/dgrijalva/jwt-go"
)

const (
	UAAPrivateKey = "PRIVATE-KEY"
	UAAPublicKey  = "PUBLIC-KEY"
)

func RegisterFastTokenSigningMethod() {
	jwt.RegisterSigningMethod("FAST", func() jwt.SigningMethod {
		return SigningMethodFast{}
	})
}

type SigningMethodFast struct{}

func (m SigningMethodFast) Alg() string {
	return "FAST"
}

func (m SigningMethodFast) Sign(signingString string, key interface{}) (string, error) {
	signature := jwt.EncodeSegment([]byte(signingString + "SUPERFAST"))
	return signature, nil
}

func (m SigningMethodFast) Verify(signingString, signature string, key interface{}) (err error) {
	if signature != jwt.EncodeSegment([]byte(signingString+"SUPERFAST")) {
		return errors.New("Signature is invalid")
	}

	return nil
}

func BuildToken(header map[string]interface{}, claims map[string]interface{}) string {
	application.UAAPublicKey = UAAPublicKey

	alg := header["alg"].(string)
	signingMethod := jwt.GetSigningMethod(alg)
	token := jwt.New(signingMethod)
	token.Header = header
	token.Claims = claims

	signed, err := token.SignedString([]byte(UAAPrivateKey))
	if err != nil {
		panic(err)
	}

	return signed
}
