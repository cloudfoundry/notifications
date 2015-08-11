package fakes

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/db"
)

type ReceiptsRepo struct {
	CreateUserGUIDs     []string
	ClientID            string
	KindID              string
	CreateReceiptsError bool
}

func NewReceiptsRepo() *ReceiptsRepo {
	return &ReceiptsRepo{}
}

func (fake *ReceiptsRepo) CreateReceipts(conn db.ConnectionInterface, userGUIDs []string, clientID, kindID string) error {
	if fake.CreateReceiptsError {
		return errors.New("a database error")
	}

	fake.CreateUserGUIDs = userGUIDs
	fake.ClientID = clientID
	fake.KindID = kindID

	return nil
}
