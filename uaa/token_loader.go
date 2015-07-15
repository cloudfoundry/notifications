package uaa

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/metrics"
)

type TokenLoader struct {
	uaaClient UAAClientInterface
}

type UAAClientInterface interface {
	GetClientToken(string) (string, error)
}

func NewTokenLoader(uaaClient UAAClientInterface) *TokenLoader {
	return &TokenLoader{
		uaaClient: uaaClient,
	}
}

func (t *TokenLoader) Load(uaaHost string) (string, error) {
	then := time.Now()

	token, err := t.uaaClient.GetClientToken(uaaHost)

	duration := time.Now().Sub(then)

	metrics.NewMetric("histogram", map[string]interface{}{
		"name":  "notifications.external-requests.uaa.client-token",
		"value": duration.Seconds(),
	}).Log()
	return token, err
}
