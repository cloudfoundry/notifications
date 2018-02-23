package uaa

import (
	"time"

	metrics "github.com/rcrowley/go-metrics"
)

type uaaClient interface {
	GetClientToken(string) (string, error)
}

type TokenLoader struct {
	uaa uaaClient
}

func NewTokenLoader(uaa uaaClient) *TokenLoader {
	return &TokenLoader{
		uaa: uaa,
	}
}

func (t *TokenLoader) Load(uaaHost string) (string, error) {
	then := time.Now()

	token, err := t.uaa.GetClientToken(uaaHost)

	metrics.GetOrRegisterTimer("notifications.external-requests.uaa.client-token", nil).Update(time.Since(then))
	return token, err
}
