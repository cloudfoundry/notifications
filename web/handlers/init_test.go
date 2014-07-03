package handlers_test

import (
    "errors"
    "testing"

    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestWebHandlersSuite(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Web Handlers Suite")
}

type FakeMailClient struct {
    messages       []mail.Message
    errorOnSend    bool
    errorOnConnect bool
}

func (fake *FakeMailClient) Connect() error {
    if fake.errorOnConnect {
        return errors.New("BOOM!")
    }
    return nil
}

func (fake *FakeMailClient) Send(msg mail.Message) error {
    err := fake.Connect()
    if err != nil {
        return err
    }

    if fake.errorOnSend {
        return errors.New("BOOM!")
    }

    fake.messages = append(fake.messages, msg)
    return nil
}

type FakeUAAClient struct {
    UsersByID        map[string]uaa.User
    ErrorForUserByID error
}

func (fake FakeUAAClient) AuthorizeURL() string {
    return ""
}

func (fake FakeUAAClient) LoginURL() string {
    return ""
}

func (fake FakeUAAClient) SetToken(token string) {}

func (fake FakeUAAClient) Exchange(code string) (uaa.Token, error) {
    return uaa.Token{}, nil
}

func (fake FakeUAAClient) Refresh(token string) (uaa.Token, error) {
    return uaa.Token{}, nil
}

func (fake FakeUAAClient) GetClientToken() (uaa.Token, error) {
    return uaa.Token{}, nil
}

func (fake FakeUAAClient) UserByID(id string) (uaa.User, error) {
    return fake.UsersByID[id], fake.ErrorForUserByID
}
