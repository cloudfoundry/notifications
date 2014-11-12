package fakes

import "github.com/cloudfoundry-incubator/notifications/mail"

type MailClient struct {
	Messages     []mail.Message
	SendError    error
	ConnectError error
}

func NewMailClient() MailClient {
	return MailClient{}
}

func (fake *MailClient) Connect() error {
	return fake.ConnectError
}

func (fake *MailClient) Send(msg mail.Message) error {
	err := fake.Connect()
	if err != nil {
		return err
	}

	if fake.SendError != nil {
		return fake.SendError
	}

	fake.Messages = append(fake.Messages, msg)
	return nil
}
