package uaa

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/metrics"
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

	duration := time.Now().Sub(then)

	metrics.NewMetric("histogram", map[string]interface{}{
		"name":  "notifications.external-requests.uaa.client-token",
		"value": duration.Seconds(),
	}).Log()
	return token, err
}
