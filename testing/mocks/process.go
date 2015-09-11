package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/pivotal-golang/lager"
)

type Process struct {
	DeliverCall struct {
		CallCount int
		Receives  struct {
			Job    *gobble.Job
			Logger lager.Logger
		}
		Returns struct {
			Error error
		}
	}
}

func NewProcess() *Process {
	return &Process{}
}

func (p *Process) Deliver(job *gobble.Job, logger lager.Logger) error {
	p.DeliverCall.Receives.Job = job
	p.DeliverCall.Receives.Logger = logger
	p.DeliverCall.CallCount++

	return p.DeliverCall.Returns.Error
}
