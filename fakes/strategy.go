package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"
)

type Strategy struct {
	DispatchArguments []interface{}
	Responses         []strategies.Response
	Error             error
	TrimCalled        bool
}

func NewStrategy() *Strategy {
	return &Strategy{}
}

func (s *Strategy) Dispatch(clientID, guid, vcapRequestID string,
	options postal.Options, conn models.ConnectionInterface) ([]strategies.Response, error) {

	s.DispatchArguments = []interface{}{clientID, guid, vcapRequestID, options}
	return s.Responses, s.Error
}
