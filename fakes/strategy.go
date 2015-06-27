package fakes

import "github.com/cloudfoundry-incubator/notifications/postal/strategies"

type Strategy struct {
	DispatchCall struct {
		Dispatch  strategies.Dispatch
		Responses []strategies.Response
		Error     error
	}
}

func NewStrategy() *Strategy {
	return &Strategy{}
}

func (s *Strategy) Dispatch(dispatch strategies.Dispatch) ([]strategies.Response, error) {
	s.DispatchCall.Dispatch = dispatch
	return s.DispatchCall.Responses, s.DispatchCall.Error
}
