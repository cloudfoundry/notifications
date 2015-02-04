package postal

import (
	"sync"
	"time"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type uaaTokenInterface interface {
	uaa.SetTokenInterface
	uaa.GetClientTokenInterface
}

var token uaa.Token
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
		token, err = loader.getNewClientToken()
		if err != nil {
			err = UAAErrorFor(err)
			return "", err
		}
	}
	loader.uaaClient.SetToken(token.Access)

	return token.Access, nil
}

func ResetLoader() {
	token = uaa.Token{}
}

func (loader TokenLoader) expired() bool {
	toleratedRequestTime := time.Duration(30 * time.Second)
	expired, err := token.ExpiresBefore(time.Duration(toleratedRequestTime))
	if err != nil {
		return true
	}

	return expired
}

func (loader TokenLoader) newTokenRequired() bool {
	return token.Access == "" || loader.expired()
}

func (loader TokenLoader) getNewClientToken() (uaa.Token, error) {
	then := time.Now()

	token, err := loader.uaaClient.GetClientToken()

	duration := time.Now().Sub(then)

	metrics.NewMetric("histogram", map[string]interface{}{
		"name":  "notifications.external-requests.uaa.client-token",
		"value": duration.Seconds(),
	}).Log()

	return token, err
}
