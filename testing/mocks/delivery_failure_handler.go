package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/postal/common"
	"github.com/pivotal-golang/lager"
)

type DeliveryFailureHandler struct {
	HandleCall struct {
		WasCalled bool
		Receives  struct {
			Job    common.Retryable
			Logger lager.Logger
		}
	}
}

func NewDeliveryFailureHandler() *DeliveryFailureHandler {
	return &DeliveryFailureHandler{}
}

func (h *DeliveryFailureHandler) Handle(job common.Retryable, logger lager.Logger) {
	h.HandleCall.WasCalled = true
	h.HandleCall.Receives.Job = job
	h.HandleCall.Receives.Logger = logger
}
