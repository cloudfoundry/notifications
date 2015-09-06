package mocks

import "github.com/cloudfoundry-incubator/notifications/v1/models"

type UnsubscribesRepo struct {
	GetCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			UserID     string
			ClientID   string
			KindID     string
		}
		Returns struct {
			Unsubscribed bool
			Error        error
		}
	}

	SetCall struct {
		Receives struct {
			Connection  models.ConnectionInterface
			UserID      string
			ClientID    string
			KindID      string
			Unsubscribe bool
		}
		Returns struct {
			Error error
		}
	}
}

func NewUnsubscribesRepo() *UnsubscribesRepo {
	return &UnsubscribesRepo{}
}

func (ur *UnsubscribesRepo) Get(conn models.ConnectionInterface, userID, clientID, kindID string) (bool, error) {
	ur.GetCall.Receives.Connection = conn
	ur.GetCall.Receives.UserID = userID
	ur.GetCall.Receives.ClientID = clientID
	ur.GetCall.Receives.KindID = kindID

	return ur.GetCall.Returns.Unsubscribed, ur.GetCall.Returns.Error
}

func (ur *UnsubscribesRepo) Set(conn models.ConnectionInterface, userID, clientID, kindID string, unsubscribe bool) error {
	ur.SetCall.Receives.Connection = conn
	ur.SetCall.Receives.UserID = userID
	ur.SetCall.Receives.ClientID = clientID
	ur.SetCall.Receives.KindID = kindID
	ur.SetCall.Receives.Unsubscribe = unsubscribe

	return ur.SetCall.Returns.Error
}
