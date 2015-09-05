package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/pivotal-golang/lager"
)

type Workflow struct {
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

func NewWorkflow() *Workflow {
	return &Workflow{}
}

func (w *Workflow) Deliver(job *gobble.Job, logger lager.Logger) error {
	w.DeliverCall.Receives.Job = job
	w.DeliverCall.Receives.Logger = logger
	w.DeliverCall.CallCount++

	return w.DeliverCall.Returns.Error
}
