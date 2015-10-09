package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/pivotal-golang/lager"
)

type V1DeliveryJobProcessor struct {
	ProcessCall struct {
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

func NewV1DeliveryJobProcessor() *V1DeliveryJobProcessor {
	return &V1DeliveryJobProcessor{}
}

func (p *V1DeliveryJobProcessor) Process(job *gobble.Job, logger lager.Logger) error {
	p.ProcessCall.Receives.Job = job
	p.ProcessCall.Receives.Logger = logger
	p.ProcessCall.CallCount++

	return p.ProcessCall.Returns.Error
}
