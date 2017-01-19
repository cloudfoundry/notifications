package domain

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

type Tokens struct {
	DefaultScopes []string
	PublicKey     string
	PrivateKey    string
}

func NewTokens(publicKey, privateKey string, defaultScopes []string) *Tokens {
	return &Tokens{
		DefaultScopes: defaultScopes,
		PublicKey:     publicKey,
		PrivateKey:    privateKey,
	}
}

func (t Tokens) Encrypt(token Token) string {
	crypt := jwt.New(jwt.SigningMethodRS256)
	crypt.Claims = token.toClaims()
	crypt.Header["kid"] = "legacy-token-key"
	encrypted, err := crypt.SignedString([]byte(t.PrivateKey))
	if err != nil {
		panic(err)
	}

	return encrypted
}

func (t Tokens) Decrypt(encryptedToken string) (Token, error) {
	tok, err := jwt.Parse(encryptedToken, jwt.Keyfunc(func(token *jwt.Token) (interface{}, error) {
		switch token.Method {
		case jwt.SigningMethodRS256, jwt.SigningMethodRS384, jwt.SigningMethodRS512:
			return []byte(t.PublicKey), nil
		default:
			return nil, errors.New("Unsupported signing method")
		}
	}))
	if err != nil {
		return Token{}, err
	}

	return newTokenFromClaims(tok.Claims), nil
}

func (t Tokens) Validate(encryptedToken string, expectedToken Token) bool {
	decryptedToken, err := t.Decrypt(encryptedToken)
	if err != nil {
		return false
	}

	return t.validate(decryptedToken, expectedToken)
}

func (t Tokens) validate(tok, expected Token) bool {
	if ok := tok.hasAudiences(expected.Audiences); !ok {
		return false
	}

	if ok := tok.hasScopes(expected.Scopes); !ok {
		return false
	}

	if ok := tok.hasAuthorities(expected.Authorities); !ok {
		return false
	}

	return true
}
