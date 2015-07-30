package collections

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/models"
)

type Sender struct {
	ID       string
	Name     string
	ClientID string
}

type sendersRepository interface {
	Insert(conn models.ConnectionInterface, sender models.Sender) (insertedSender models.Sender, err error)
	Get(conn models.ConnectionInterface, senderID string) (retrievedSender models.Sender, err error)
	GetByClientIDAndName(conn models.ConnectionInterface, clientID, name string) (retrievedSender models.Sender, err error)
}

type SendersCollection struct {
	repo sendersRepository
}

func NewSendersCollection(repo sendersRepository) SendersCollection {
	return SendersCollection{
		repo: repo,
	}
}

func (sc SendersCollection) Add(conn models.ConnectionInterface, sender Sender) (Sender, error) {
	model, err := sc.repo.Insert(conn, models.Sender{
		Name:     sender.Name,
		ClientID: sender.ClientID,
	})
	if err != nil {
		switch err.(type) {
		case models.DuplicateRecordError:
			model, err = sc.repo.GetByClientIDAndName(conn, sender.ClientID, sender.Name)
			if err != nil {
				return Sender{}, PersistenceError{err}
			}
		default:
			return Sender{}, PersistenceError{err}
		}
	}

	return Sender{
		ID:       model.ID,
		Name:     model.Name,
		ClientID: model.ClientID,
	}, nil
}

func (sc SendersCollection) Get(conn models.ConnectionInterface, senderID, clientID string) (Sender, error) {
	model, err := sc.repo.Get(conn, senderID)
	if err != nil {
		switch e := err.(type) {
		case models.RecordNotFoundError:
			return Sender{}, NotFoundError{
				Err:     e,
				Message: string(e),
			}
		default:
			return Sender{}, PersistenceError{err}
		}
	}

	if clientID != model.ClientID {
		return Sender{}, NotFoundError{
			Err:     errors.New("sender not found"),
			Message: "sender not found",
		}
	}

	return Sender{
		ID:       model.ID,
		Name:     model.Name,
		ClientID: model.ClientID,
	}, nil
}
