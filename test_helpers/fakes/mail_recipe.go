package fakes

import (
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
)

type FakeMailRecipe struct {
    DeliverMailArguments []interface{}
    Responses            []postal.Response
    Error                error
    TrimCalled           bool
}

func (fake *FakeMailRecipe) DeliverMail(clientID string, guid postal.TypedGUID,
    options postal.Options, conn models.ConnectionInterface) ([]postal.Response, error) {

    fake.DeliverMailArguments = []interface{}{clientID, guid, options}
    return fake.Responses, fake.Error
}

func (fake *FakeMailRecipe) Trim(response []byte) []byte {
    fake.TrimCalled = true
    return response
}
