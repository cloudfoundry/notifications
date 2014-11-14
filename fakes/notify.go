package fakes

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/ryanmoran/stack"
)

type Notify struct {
	Response []byte
	GUID     string
	Error    error
}

func NewNotify() *Notify {
	return &Notify{}
}

func (fake *Notify) Execute(connection models.ConnectionInterface, req *http.Request, context stack.Context,
	guid string, strategy strategies.StrategyInterface, validator handlers.ValidatorInterface) ([]byte, error) {
	fake.GUID = guid

	return fake.Response, fake.Error
}
