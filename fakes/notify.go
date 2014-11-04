package fakes

import (
    "net/http"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/ryanmoran/stack"
)

type FakeNotify struct {
    Response []byte
    GUID     postal.TypedGUID
    Error    error
}

func (fake *FakeNotify) Execute(connection models.ConnectionInterface, req *http.Request, context stack.Context,
    guid postal.TypedGUID, recipe postal.RecipeInterface) ([]byte, error) {
    fake.GUID = guid

    return fake.Response, fake.Error
}
