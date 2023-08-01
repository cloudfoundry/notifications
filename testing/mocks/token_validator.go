package mocks

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/pivotal-cf-experimental/warrant"
)

type TokenValidator struct {
	ParseCall struct {
		Receives struct {
			Token string
		}

		Returns struct {
			Token *jwt.Token
			Error error
		}
	}
}

func (t *TokenValidator) Parse(token string) (*jwt.Token, error) {
	t.ParseCall.Receives.Token = token
	return t.ParseCall.Returns.Token, t.ParseCall.Returns.Error
}

type KeyFetcher struct {
	GetSigningKeysCall struct {
		Called  bool
		Returns struct {
			Keys  []warrant.SigningKey
			Error error
		}
	}
}

func (f *KeyFetcher) GetSigningKeys() ([]warrant.SigningKey, error) {
	f.GetSigningKeysCall.Called = true
	return f.GetSigningKeysCall.Returns.Keys, f.GetSigningKeysCall.Returns.Error
}
