package collections

import (
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/v2/models"
)

type Sender struct {
	ID       string
	Name     string
	ClientID string
}

type sendersRepository interface {
	Insert(conn models.ConnectionInterface, sender models.Sender) (insertedSender models.Sender, err error)
	Update(conn models.ConnectionInterface, sender models.Sender) (models.Sender, error)
	List(conn models.ConnectionInterface, clientID string) (retrievedSenderList []models.Sender, er error)
	Get(conn models.ConnectionInterface, senderID string) (retrievedSender models.Sender, err error)
	GetByClientIDAndName(conn models.ConnectionInterface, clientID, name string) (retrievedSender models.Sender, err error)
	Delete(conn models.ConnectionInterface, sender models.Sender) error
}

type SendersCollection struct {
	senders       sendersRepository
	campaignTypes campaignTypesRepository
}

func NewSendersCollection(senders sendersRepository, campaignTypes campaignTypesRepository) SendersCollection {
	return SendersCollection{
		senders:       senders,
		campaignTypes: campaignTypes,
	}
}

func (sc SendersCollection) Set(conn ConnectionInterface, sender Sender) (Sender, error) {
	var (
		model models.Sender
		err   error
	)

	if sender.ID == "" {
		model, err = sc.senders.Insert(conn, models.Sender{
			Name:     sender.Name,
			ClientID: sender.ClientID,
		})
		if err != nil {
			switch err.(type) {
			case models.DuplicateRecordError:
				model, err = sc.senders.GetByClientIDAndName(conn, sender.ClientID, sender.Name)
				if err != nil {
					return Sender{}, PersistenceError{err}
				}
			default:
				return Sender{}, PersistenceError{err}
			}
		}
	} else {
		model, err = sc.senders.Update(conn, models.Sender{
			ID:       sender.ID,
			Name:     sender.Name,
			ClientID: sender.ClientID,
		})
		if err != nil {
			switch err.(type) {
			case models.DuplicateRecordError:
				return Sender{}, DuplicateRecordError{err}
			default:
				return Sender{}, PersistenceError{err}
			}
		}

	}

	return Sender{
		ID:       model.ID,
		Name:     model.Name,
		ClientID: model.ClientID,
	}, nil
}

func (sc SendersCollection) List(conn ConnectionInterface, clientID string) ([]Sender, error) {
	senderList := []Sender{}

	models, err := sc.senders.List(conn, clientID)
	if err != nil {
		return senderList, PersistenceError{err}
	}

	for _, model := range models {
		sender := Sender{
			ID:       model.ID,
			Name:     model.Name,
			ClientID: model.ClientID,
		}

		senderList = append(senderList, sender)
	}

	return senderList, nil
}

func (sc SendersCollection) Get(conn ConnectionInterface, senderID, clientID string) (Sender, error) {
	model, err := sc.senders.Get(conn, senderID)
	if err != nil {
		switch e := err.(type) {
		case models.RecordNotFoundError:
			return Sender{}, NotFoundError{e}
		default:
			return Sender{}, PersistenceError{err}
		}
	}

	if clientID != model.ClientID {
		return Sender{}, NotFoundError{fmt.Errorf("Sender with id %q could not be found", senderID)}
	}

	return Sender{
		ID:       model.ID,
		Name:     model.Name,
		ClientID: model.ClientID,
	}, nil
}

func (sc SendersCollection) Delete(conn ConnectionInterface, senderID, clientID string) error {
	sender, err := sc.senders.Get(conn, senderID)
	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return NotFoundError{err}
		default:
			return UnknownError{err}
		}
	}

	if sender.ClientID != clientID {
		return NotFoundError{fmt.Errorf("Sender with id %q could not be found", senderID)}
	}

	err = sc.senders.Delete(conn, sender)
	if err != nil {
		return UnknownError{err}
	}

	campaignTypes, err := sc.campaignTypes.List(conn, senderID)
	if err != nil {
		return UnknownError{err}
	}

	for _, campaignType := range campaignTypes {
		err = sc.campaignTypes.Delete(conn, campaignType)
		if err != nil {
			return UnknownError{err}
		}
	}

	return nil
}
