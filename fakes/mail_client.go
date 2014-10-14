package fakes

import "github.com/cloudfoundry-incubator/notifications/mail"

type FakeMailClient struct {
    Messages     []mail.Message
    SendError    error
    ConnectError error
}

func (fake *FakeMailClient) Connect() error {
    return fake.ConnectError
}

func (fake *FakeMailClient) Send(msg mail.Message) error {
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
