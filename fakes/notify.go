package fakes

import (
    "net/http"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/ryanmoran/stack"
)

type Notify struct {
    Response []byte
    GUID     postal.TypedGUID
    Error    error
}

func NewNotify() *Notify {
    return &Notify{}
}

func (fake *Notify) Execute(connection models.ConnectionInterface, req *http.Request, context stack.Context,
    guid postal.TypedGUID, strategy postal.StrategyInterface) ([]byte, error) {
    fake.GUID = guid

    return fake.Response, fake.Error
}
