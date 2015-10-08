package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/postal/common"
	"github.com/pivotal-golang/lager"
)

type V2DeliveryJobProcessor struct {
	ProcessCall struct {
		CallCount int
		Receives  struct {
			Delivery common.Delivery
			Logger   lager.Logger
		}
		Returns struct {
			Error error
		}
	}
}

func NewV2DeliveryJobProcessor() *V2DeliveryJobProcessor {
	return &V2DeliveryJobProcessor{}
}

func (p *V2DeliveryJobProcessor) Process(delivery common.Delivery, logger lager.Logger) error {
	p.ProcessCall.Receives.Delivery = delivery
	p.ProcessCall.Receives.Logger = logger
	p.ProcessCall.CallCount++

	return p.ProcessCall.Returns.Error
}
