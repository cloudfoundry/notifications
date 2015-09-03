package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/mail"

	"github.com/pivotal-golang/lager"
)

type MailClient struct {
	ConnectCall struct {
		Receives struct {
			Logger lager.Logger
		}
		Returns struct {
			Error error
		}
	}

	SendCall struct {
		CallCount int
		Receives  struct {
			Message mail.Message
			Logger  lager.Logger
		}
		Returns struct {
			Error error
		}
	}
}

func NewMailClient() *MailClient {
	return &MailClient{}
}

func (mc *MailClient) Connect(logger lager.Logger) error {
	mc.ConnectCall.Receives.Logger = logger

	return mc.ConnectCall.Returns.Error
}

func (mc *MailClient) Send(message mail.Message, logger lager.Logger) error {
	mc.SendCall.Receives.Message = message
	mc.SendCall.Receives.Logger = logger
	mc.SendCall.CallCount++

	return mc.SendCall.Returns.Error
}
