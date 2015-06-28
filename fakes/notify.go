package fakes

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/ryanmoran/stack"
)

type Notify struct {
	ExecuteCall struct {
		Args struct {
			Connection    models.ConnectionInterface
			Request       *http.Request
			Context       stack.Context
			GUID          string
			Strategy      services.StrategyInterface
			Validator     handlers.ValidatorInterface
			VCAPRequestID string
		}
		Response []byte
		Error    error
	}
}

func NewNotify() *Notify {
	return &Notify{}
}

func (fake *Notify) Execute(connection models.ConnectionInterface, req *http.Request, context stack.Context,
	guid string, strategy services.StrategyInterface, validator handlers.ValidatorInterface, vcapRequestID string) ([]byte, error) {

	fake.ExecuteCall.Args.Connection = connection
	fake.ExecuteCall.Args.Request = req
	fake.ExecuteCall.Args.Context = context
	fake.ExecuteCall.Args.GUID = guid
	fake.ExecuteCall.Args.Strategy = strategy
	fake.ExecuteCall.Args.Validator = validator
	fake.ExecuteCall.Args.VCAPRequestID = vcapRequestID

	return fake.ExecuteCall.Response, fake.ExecuteCall.Error
}
