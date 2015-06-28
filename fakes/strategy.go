package fakes

import "github.com/cloudfoundry-incubator/notifications/services"

type Strategy struct {
	DispatchCall struct {
		Dispatch  services.Dispatch
		Responses []services.Response
		Error     error
	}
}

func NewStrategy() *Strategy {
	return &Strategy{}
}

func (s *Strategy) Dispatch(dispatch services.Dispatch) ([]services.Response, error) {
	s.DispatchCall.Dispatch = dispatch
	return s.DispatchCall.Responses, s.DispatchCall.Error
}
