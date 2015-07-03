package postal

import (
	"sync"
	"time"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/uaa"
)

type uaaTokenInterface interface {
	SetToken(string)
	GetClientToken() (string, error)
}

var accessToken string
var mutex sync.Mutex

type TokenLoader struct {
	uaaClient uaaTokenInterface
}

type TokenLoaderInterface interface {
	Load() (string, error)
}

func NewTokenLoader(uaaClient uaaTokenInterface) TokenLoader {
	return TokenLoader{
		uaaClient: uaaClient,
	}
}

func (loader TokenLoader) Load() (string, error) {
	var err error

	mutex.Lock()
	defer mutex.Unlock()

	if loader.newTokenRequired() {
		accessToken, err = loader.getNewClientToken()
		if err != nil {
			err = UAAErrorFor(err)
			return "", err
		}
	}
	loader.uaaClient.SetToken(accessToken)

	return accessToken, nil
}

func ResetLoader() {
	accessToken = ""
}

func (loader TokenLoader) expired() bool {
	toleratedRequestTime := time.Duration(30 * time.Second)
	expired, err := uaa.AccessTokenExpiresBefore(accessToken, time.Duration(toleratedRequestTime))
	if err != nil {
		return true
	}

	return expired
}

func (loader TokenLoader) newTokenRequired() bool {
	return accessToken == "" || loader.expired()
}

func (loader TokenLoader) getNewClientToken() (string, error) {
	then := time.Now()

	at, err := loader.uaaClient.GetClientToken()

	duration := time.Now().Sub(then)

	metrics.NewMetric("histogram", map[string]interface{}{
		"name":  "notifications.external-requests.uaa.client-token",
		"value": duration.Seconds(),
	}).Log()

	return at, err
}
