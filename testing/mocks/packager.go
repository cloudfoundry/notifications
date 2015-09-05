package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/mail"
	"github.com/cloudfoundry-incubator/notifications/postal"
)

type Packager struct {
	PrepareContextCall struct {
		Receives struct {
			Delivery postal.Delivery
			Sender   string
			Domain   string
		}
		Returns struct {
			MessageContext postal.MessageContext
			Error          error
		}
	}

	PackCall struct {
		Receives struct {
			MessageContext postal.MessageContext
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

func (p *Packager) PrepareContext(delivery postal.Delivery, sender, domain string) (postal.MessageContext, error) {
	p.PrepareContextCall.Receives.Delivery = delivery
	p.PrepareContextCall.Receives.Sender = sender
	p.PrepareContextCall.Receives.Domain = domain

	return p.PrepareContextCall.Returns.MessageContext, p.PrepareContextCall.Returns.Error
}

func (p *Packager) Pack(context postal.MessageContext) (mail.Message, error) {
	p.PackCall.Receives.MessageContext = context

	return p.PackCall.Returns.Message, p.PackCall.Returns.Error
}
