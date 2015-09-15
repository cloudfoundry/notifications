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
	DeleteCall struct {
		Receives struct {
			Unsubscriber collections.Unsubscriber
			Connection   db.ConnectionInterface
		}
		Returns struct {
			Error error
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

func (u *UnsubscribersCollection) Delete(connection collections.ConnectionInterface, unsubscriber collections.Unsubscriber) error {
	u.DeleteCall.Receives.Unsubscriber = unsubscriber
	u.DeleteCall.Receives.Connection = connection

	return u.DeleteCall.Returns.Error
}
