package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
)

type UnsubscribersCollection struct {
	SetCall struct {
		Receives struct {
			Unsubscriber collections.Unsubscriber
			Connection   db.ConnectionInterface
		}
		Returns struct {
			Unsubscriber collections.Unsubscriber
			Error        error
		}
	}
}

func NewUnsubscribersCollection() *UnsubscribersCollection {
	return &UnsubscribersCollection{}
}

func (u *UnsubscribersCollection) Set(connection collections.ConnectionInterface, unsubscriber collections.Unsubscriber) (collections.Unsubscriber, error) {
	u.SetCall.Receives.Unsubscriber = unsubscriber
	u.SetCall.Receives.Connection = connection

	return u.SetCall.Returns.Unsubscriber, u.SetCall.Returns.Error
}
