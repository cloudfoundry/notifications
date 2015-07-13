package uaa

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/metrics"
)

type ZonedTokenLoader struct {
	uaaClient UAAClientInterface
}

type UAAClientInterface interface {
	ZonedGetClientToken(string) (string, error)
}

func NewZonedTokenLoader(uaaClient UAAClientInterface) *ZonedTokenLoader {
	return &ZonedTokenLoader{
		uaaClient: uaaClient,
	}
}

func (z *ZonedTokenLoader) Load(uaaHost string) (string, error) {
	then := time.Now()

	token, err := z.uaaClient.ZonedGetClientToken(uaaHost)

	duration := time.Now().Sub(then)

	metrics.NewMetric("histogram", map[string]interface{}{
		"name":  "notifications.external-requests.uaa.client-token",
		"value": duration.Seconds(),
	}).Log()
	return token, err
}
