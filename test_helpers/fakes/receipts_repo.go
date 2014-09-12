package fakes

import (
    "errors"

    "github.com/cloudfoundry-incubator/notifications/models"
)

type FakeReceiptsRepo struct {
    CreateUserGUIDs     []string
    ClientID            string
    KindID              string
    CreateReceiptsError bool
}

func NewFakeReceiptsRepo() FakeReceiptsRepo {
    return FakeReceiptsRepo{}
}

func (fake *FakeReceiptsRepo) CreateReceipts(conn models.ConnectionInterface, userGUIDs []string, clientID, kindID string) error {
    if fake.CreateReceiptsError {
        return errors.New("a database error")
    }

    fake.CreateUserGUIDs = userGUIDs
    fake.ClientID = clientID
    fake.KindID = kindID

    return nil
}
