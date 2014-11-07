package fakes

import (
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
)

type MailStrategy struct {
    DispatchArguments []interface{}
    Responses         []postal.Response
    Error             error
    TrimCalled        bool
}

func NewMailStrategy() *MailStrategy {
    return &MailStrategy{}
}

func (fake *MailStrategy) Dispatch(clientID string, guid postal.TypedGUID,
    options postal.Options, conn models.ConnectionInterface) ([]postal.Response, error) {

    fake.DispatchArguments = []interface{}{clientID, guid, options}
    return fake.Responses, fake.Error
}

func (fake *MailStrategy) Trim(response []byte) []byte {
    fake.TrimCalled = true
    return response
}
