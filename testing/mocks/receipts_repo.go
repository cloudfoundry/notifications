package mocks

import "github.com/cloudfoundry-incubator/notifications/v1/models"

type ReceiptsRepo struct {
	CreateReceiptsCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			UserGUIDs  []string
			ClientID   string
			KindID     string
		}
		Returns struct {
			Error error
		}
	}
}

func NewReceiptsRepo() *ReceiptsRepo {
	return &ReceiptsRepo{}
}

func (rr *ReceiptsRepo) CreateReceipts(conn models.ConnectionInterface, userGUIDs []string, clientID, kindID string) error {
	rr.CreateReceiptsCall.Receives.Connection = conn
	rr.CreateReceiptsCall.Receives.UserGUIDs = userGUIDs
	rr.CreateReceiptsCall.Receives.ClientID = clientID
	rr.CreateReceiptsCall.Receives.KindID = kindID

	return rr.CreateReceiptsCall.Returns.Error
}
