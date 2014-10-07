package fakes

import (
    "errors"

    "github.com/cloudfoundry-incubator/notifications/mail"
)

type FakeMailClient struct {
    Messages       []mail.Message
    ErrorOnSend    bool
    ErrorOnConnect bool
}

func (fake *FakeMailClient) Connect() error {
    if fake.ErrorOnConnect {
        return errors.New("BOOM!")
    }
    return nil
}

func (fake *FakeMailClient) Send(msg mail.Message) error {
    err := fake.Connect()
    if err != nil {
        return err
    }

    if fake.ErrorOnSend {
        return errors.New("BOOM!")
    }

    fake.Messages = append(fake.Messages, msg)
    return nil
}
