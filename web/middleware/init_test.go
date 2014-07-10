package middleware_test

import (
    "errors"
    "testing"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/dgrijalva/jwt-go"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestWebMiddlewareSuite(t *testing.T) {
    RegisterFastTokenSigningMethod()

    RegisterFailHandler(Fail)
    RunSpecs(t, "Web Middleware Suite")
}

const (
    UAAPrivateKey = "PRIVATE-KEY"
    UAAPublicKey  = "PUBLIC-KEY"
)

type SigningMethodFast struct{}

func (m SigningMethodFast) Alg() string {
    return "FAST"
}

func (m SigningMethodFast) Sign(signingString string, key []byte) (string, error) {
    signature := jwt.EncodeSegment([]byte(signingString + "SUPERFAST"))
    return signature, nil
}

func (m SigningMethodFast) Verify(signingString, signature string, key []byte) (err error) {
    if signature != jwt.EncodeSegment([]byte(signingString+"SUPERFAST")) {
        return errors.New("Signature is invalid")
    }

    return nil
}

func RegisterFastTokenSigningMethod() {
    jwt.RegisterSigningMethod("FAST", func() jwt.SigningMethod {
        return SigningMethodFast{}
    })
}

func BuildToken(header map[string]interface{}, claims map[string]interface{}) string {
    config.UAAPublicKey = UAAPublicKey

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
