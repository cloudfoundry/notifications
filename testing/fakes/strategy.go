package fakes

import "github.com/cloudfoundry-incubator/notifications/v1/services"

type Strategy struct {
	DispatchCall struct {
		Receives struct {
			Dispatch services.Dispatch
		}
		Returns struct {
			Responses []services.Response
			Error     error
		}
	}
}

func NewStrategy() *Strategy {
	return &Strategy{}
}

func (s *Strategy) Dispatch(dispatch services.Dispatch) ([]services.Response, error) {
	s.DispatchCall.Receives.Dispatch = dispatch

	return s.DispatchCall.Returns.Responses, s.DispatchCall.Returns.Error
}
