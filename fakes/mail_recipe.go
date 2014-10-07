package fakes

import (
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
)

type FakeMailRecipe struct {
    DispatchArguments []interface{}
    Responses         []postal.Response
    Error             error
    TrimCalled        bool
}

func (fake *FakeMailRecipe) Dispatch(clientID string, guid postal.TypedGUID,
    options postal.Options, conn models.ConnectionInterface) ([]postal.Response, error) {

    fake.DispatchArguments = []interface{}{clientID, guid, options}
    return fake.Responses, fake.Error
}

func (fake *FakeMailRecipe) Trim(response []byte) []byte {
    fake.TrimCalled = true
    return response
}
