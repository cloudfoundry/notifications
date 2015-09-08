package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
)

type StrategyDeterminer struct {
	DetermineCall struct {
		Receives struct {
			Connection services.ConnectionInterface
			UAAHost    string
			Job        gobble.Job
		}

		Returns struct {
			Error error
		}

		WasCalled bool
	}
}

func NewStrategyDeterminer() *StrategyDeterminer {
	return &StrategyDeterminer{}
}

func (s *StrategyDeterminer) Determine(conn services.ConnectionInterface, uaaHost string, job gobble.Job) error {
	s.DetermineCall.Receives.Connection = conn
	s.DetermineCall.Receives.UAAHost = uaaHost
	s.DetermineCall.Receives.Job = job
	s.DetermineCall.WasCalled = true

	return s.DetermineCall.Returns.Error
}
