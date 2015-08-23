package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/mail"
	"github.com/pivotal-golang/lager"
)

type MailClient struct {
	Messages      []mail.Message
	SendLogger    lager.Logger
	SendError     error
	ConnectLogger lager.Logger
	ConnectError  error
}

func NewMailClient() *MailClient {
	return &MailClient{}
}

func (fake *MailClient) Connect(logger lager.Logger) error {
	fake.ConnectLogger = logger
	return fake.ConnectError
}

func (fake *MailClient) Send(msg mail.Message, logger lager.Logger) error {
	fake.SendLogger = logger

	if fake.ConnectError != nil {
		return fake.ConnectError
	}

	if fake.SendError != nil {
		return fake.SendError
	}

	fake.Messages = append(fake.Messages, msg)
	return nil
}
