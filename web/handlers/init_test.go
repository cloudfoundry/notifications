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
    messages    []mail.Message
    errorOnSend bool
}

func (fake *FakeMailClient) Connect() error {
    return nil
}

func (fake *FakeMailClient) Send(msg mail.Message) error {
    if fake.errorOnSend {
        return errors.New("BOOM!")
    }

    fake.messages = append(fake.messages, msg)
    return nil
}

type FakeUAAClient struct {
    UsersByID map[string]uaa.User
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
    return fake.UsersByID[id], nil
}
