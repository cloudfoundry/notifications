package domain

import "github.com/dgrijalva/jwt-go"

type Tokens struct {
	DefaultScopes []string
	PublicKey     string
}

func NewTokens(publicKey string, defaultScopes []string) *Tokens {
	return &Tokens{
		DefaultScopes: defaultScopes,
		PublicKey:     publicKey,
	}
}

func (t Tokens) Encrypt(token Token) string {
	crypt := jwt.New(jwt.SigningMethodHS256)
	crypt.Claims = token.toClaims()
	encrypted, err := crypt.SignedString([]byte(t.PublicKey))
	if err != nil {
		panic(err)
	}

	return encrypted
}

func (t Tokens) Decrypt(encryptedToken string) (Token, error) {
	tok, err := jwt.Parse(encryptedToken, jwt.Keyfunc(func(*jwt.Token) (interface{}, error) {
		return []byte(t.PublicKey), nil
	}))
	if err != nil {
		return Token{}, err
	}

	return newTokenFromClaims(tok.Claims), nil
}

func (t Tokens) Validate(encryptedToken string, audiences, scopes []string) bool {
	decryptedToken, err := t.Decrypt(encryptedToken)
	if err != nil {
		return false
	}

	return t.validate(decryptedToken, Token{
		Audiences: audiences,
		Scopes:    scopes,
	})
}

func (t Tokens) validate(tok, expected Token) bool {
	if ok := tok.hasAudiences(expected.Audiences); !ok {
		return false
	}

	if ok := tok.hasScopes(expected.Scopes); !ok {
		return false
	}

	return true
}
