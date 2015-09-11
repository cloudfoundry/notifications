package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/pivotal-golang/lager"
)

type Workflow struct {
	DeliverCall struct {
		CallCount int
		Receives  struct {
			Delivery postal.Delivery
			Logger   lager.Logger
		}
		Returns struct {
			Error error
		}
	}
}

func NewWorkflow() *Workflow {
	return &Workflow{}
}

func (w *Workflow) Deliver(delivery postal.Delivery, logger lager.Logger) error {
	w.DeliverCall.Receives.Delivery = delivery
	w.DeliverCall.Receives.Logger = logger
	w.DeliverCall.CallCount++

	return w.DeliverCall.Returns.Error
}
