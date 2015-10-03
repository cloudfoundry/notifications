package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/mail"
	"github.com/cloudfoundry-incubator/notifications/postal/common"
)

type Packager struct {
	PrepareContextCall struct {
		Receives struct {
			Delivery common.Delivery
			Sender   string
			Domain   string
		}
		Returns struct {
			MessageContext common.MessageContext
			Error          error
		}
	}

	PackCall struct {
		Receives struct {
			MessageContext common.MessageContext
		}
		Returns struct {
			Message mail.Message
			Error   error
		}
	}
}

func NewPackager() *Packager {
	return &Packager{}
}

func (p *Packager) PrepareContext(delivery common.Delivery, sender, domain string) (common.MessageContext, error) {
	p.PrepareContextCall.Receives.Delivery = delivery
	p.PrepareContextCall.Receives.Sender = sender
	p.PrepareContextCall.Receives.Domain = domain

	return p.PrepareContextCall.Returns.MessageContext, p.PrepareContextCall.Returns.Error
}

func (p *Packager) Pack(context common.MessageContext) (mail.Message, error) {
	p.PackCall.Receives.MessageContext = context

	return p.PackCall.Returns.Message, p.PackCall.Returns.Error
}
