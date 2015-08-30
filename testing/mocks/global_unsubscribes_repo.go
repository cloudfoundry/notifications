package mocks

import "github.com/cloudfoundry-incubator/notifications/v1/models"

type GlobalUnsubscribesRepo struct {
	GetCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			UserID     string
		}
		Returns struct {
			Unsubscribed bool
			Error        error
		}
	}

	SetCall struct {
		Receives struct {
			Connection   models.ConnectionInterface
			UserID       string
			Unsubscribed bool
		}
		Returns struct {
			Error error
		}
	}
}

func NewGlobalUnsubscribesRepo() *GlobalUnsubscribesRepo {
	return &GlobalUnsubscribesRepo{}
}

func (r *GlobalUnsubscribesRepo) Get(conn models.ConnectionInterface, userID string) (bool, error) {
	r.GetCall.Receives.Connection = conn
	r.GetCall.Receives.UserID = userID

	return r.GetCall.Returns.Unsubscribed, r.GetCall.Returns.Error
}

func (r *GlobalUnsubscribesRepo) Set(conn models.ConnectionInterface, userID string, unsubscribed bool) error {
	r.SetCall.Receives.Connection = conn
	r.SetCall.Receives.UserID = userID
	r.SetCall.Receives.Unsubscribed = unsubscribed

	return r.SetCall.Returns.Error
}
