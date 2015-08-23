package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/pivotal-golang/lager"
)

type DeliveryFailureHandler struct {
	HandleCall struct {
		Receives struct {
			Job    postal.Retryable
			Logger lager.Logger
		}
	}
}

func NewDeliveryFailureHandler() *DeliveryFailureHandler {
	return &DeliveryFailureHandler{}
}

func (h *DeliveryFailureHandler) Handle(job postal.Retryable, logger lager.Logger) {
	h.HandleCall.Receives.Job = job
	h.HandleCall.Receives.Logger = logger
}
