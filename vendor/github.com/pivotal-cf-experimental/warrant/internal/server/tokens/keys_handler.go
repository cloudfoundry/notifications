package tokens

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"math/big"
	"net/http"

	"github.com/pivotal-cf-experimental/warrant/internal/documents"
)

type keysHandler struct {
	publicKey string
}

func (h keysHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	pem, _ := pem.Decode([]byte(h.publicKey))

	if pem == nil {
		panic("No PEM data was included in the public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(pem.Bytes)

	if err != nil {
		panic(err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)

	if !ok {
		panic("public key is not rsa")
	}

	exponentBytes := big.NewInt(int64(rsaPublicKey.E)).Bytes()
	modulusBytes := rsaPublicKey.N.Bytes()

	response, err := json.Marshal(documents.TokenKeysResponse{
		Keys: []documents.TokenKeyResponse{
			{
				Kid:   "legacy-token-key",
				Alg:   "SHA256withRSA",
				Value: h.publicKey,
				Kty:   "RSA",
				Use:   "sig",
				N:     base64.RawURLEncoding.EncodeToString(modulusBytes),
				E:     base64.RawURLEncoding.EncodeToString(exponentBytes),
			},
			{
				Kid:   "token-key",
				Alg:   "SHA256withRSA",
				Value: h.publicKey,
				Kty:   "RSA",
				Use:   "sig",
				N:     base64.RawURLEncoding.EncodeToString(modulusBytes),
				E:     base64.RawURLEncoding.EncodeToString(exponentBytes),
			},
		},
	})

	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
