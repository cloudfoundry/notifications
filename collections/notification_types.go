package collections

import (
	"errors"
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/models"
)

type CampaignType struct {
	ID          string
	Name        string
	Description string
	Critical    bool
	TemplateID  string
	SenderID    string
}

type CampaignTypesCollection struct {
	notificationTypesRepository notificationTypesRepository
	sendersRepository           sendersRepository
}

type notificationTypesRepository interface {
	Insert(models.ConnectionInterface, models.NotificationType) (models.NotificationType, error)
	GetBySenderIDAndName(models.ConnectionInterface, string, string) (models.NotificationType, error)
	List(models.ConnectionInterface, string) ([]models.NotificationType, error)
	Get(models.ConnectionInterface, string) (models.NotificationType, error)
}

func NewCampaignTypesCollection(nr notificationTypesRepository, sr sendersRepository) CampaignTypesCollection {
	return CampaignTypesCollection{
		notificationTypesRepository: nr,
		sendersRepository:           sr,
	}
}

func (nc CampaignTypesCollection) Add(conn models.ConnectionInterface, notificationType CampaignType, clientID string) (CampaignType, error) {
	senderModel, err := nc.sendersRepository.Get(conn, notificationType.SenderID)

	if err != nil {
		switch e := err.(type) {
		case models.RecordNotFoundError:
			return CampaignType{}, NotFoundError{
				Err:     e,
				Message: string(e),
			}
		default:
			return CampaignType{}, PersistenceError{err}
		}
	}

	if senderModel.ClientID != clientID {
		return CampaignType{}, NotFoundError{
			Err:     errors.New("sender not found"),
			Message: "sender not found",
		}
	}

	if notificationType.Name == "" {
		return CampaignType{}, ValidationError{
			Err:     errors.New("missing campaign type name"),
			Message: "missing campaign type name",
		}
	}

	if notificationType.Description == "" {
		return CampaignType{}, ValidationError{
			Err:     errors.New("missing campaign type description"),
			Message: "missing campaign type description",
		}
	}

	returnNotificationType, err := nc.notificationTypesRepository.Insert(conn, models.NotificationType{
		Name:        notificationType.Name,
		Description: notificationType.Description,
		Critical:    notificationType.Critical,
		TemplateID:  notificationType.TemplateID,
		SenderID:    notificationType.SenderID,
	})
	if err != nil {
		switch err.(type) {
		case models.DuplicateRecordError:
			returnNotificationType, err = nc.notificationTypesRepository.GetBySenderIDAndName(conn, notificationType.SenderID, notificationType.Name)
			if err != nil {
				return CampaignType{}, PersistenceError{err}
			}
		default:
			return CampaignType{}, PersistenceError{err}
		}
	}

	return CampaignType{
		ID:          returnNotificationType.ID,
		Name:        returnNotificationType.Name,
		Description: returnNotificationType.Description,
		Critical:    returnNotificationType.Critical,
		TemplateID:  returnNotificationType.TemplateID,
		SenderID:    returnNotificationType.SenderID,
	}, err
}

func (nc CampaignTypesCollection) List(conn models.ConnectionInterface, senderID, clientID string) ([]CampaignType, error) {
	if senderID == "" {
		return []CampaignType{}, ValidationError{
			Err:     errors.New("missing sender id"),
			Message: "missing sender id",
		}
	}

	if clientID == "" {
		return []CampaignType{}, NewValidationError("missing client id")
	}

	senderModel, err := nc.sendersRepository.Get(conn, senderID)

	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return []CampaignType{}, NotFoundError{
				Err:     err,
				Message: "sender not found",
			}
		default:
			return []CampaignType{}, PersistenceError{err}
		}
	}

	if senderModel.ClientID != clientID {
		return []CampaignType{}, NewNotFoundError("sender not found")
	}

	modelList, err := nc.notificationTypesRepository.List(conn, senderID)
	if err != nil {
		return []CampaignType{}, PersistenceError{err}
	}

	notificationTypeList := []CampaignType{}

	for _, model := range modelList {
		notificationType := CampaignType{
			ID:          model.ID,
			Name:        model.Name,
			Description: model.Description,
			Critical:    model.Critical,
			TemplateID:  model.TemplateID,
			SenderID:    model.SenderID,
		}
		notificationTypeList = append(notificationTypeList, notificationType)
	}

	return notificationTypeList, nil
}

func (nc CampaignTypesCollection) Get(conn models.ConnectionInterface, notificationTypeID, senderID, clientID string) (CampaignType, error) {
	if notificationTypeID == "" {
		return CampaignType{}, ValidationError{
			Err:     errors.New("missing campaign type id"),
			Message: "missing campaign type id",
		}
	}

	if senderID == "" {
		return CampaignType{}, ValidationError{
			Err:     errors.New("missing sender id"),
			Message: "missing sender id",
		}
	}

	if clientID == "" {
		return CampaignType{}, ValidationError{
			Err:     errors.New("missing client id"),
			Message: "missing client id",
		}
	}

	sender, err := nc.sendersRepository.Get(conn, senderID)
	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return CampaignType{}, NewNotFoundError(fmt.Sprintf("sender %s not found", notificationTypeID))
		default:
			return CampaignType{}, PersistenceError{err}
		}
	}

	if clientID != sender.ClientID {
		message := fmt.Sprintf("sender %s not found", senderID)
		return CampaignType{}, NotFoundError{
			Err:     errors.New(message),
			Message: message,
		}
	}

	notificationType, err := nc.notificationTypesRepository.Get(conn, notificationTypeID)
	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return CampaignType{}, NotFoundError{
				Err:     err,
				Message: fmt.Sprintf("campaign type %s not found", notificationTypeID),
			}
		default:
			return CampaignType{}, PersistenceError{err}
		}
	}

	if senderID != notificationType.SenderID {
		message := fmt.Sprintf("campaign type %s not found", notificationTypeID)
		return CampaignType{}, NotFoundError{
			Err:     errors.New(message),
			Message: message,
		}
	}

	return CampaignType{
		ID:          notificationType.ID,
		Name:        notificationType.Name,
		Description: notificationType.Description,
		Critical:    notificationType.Critical,
		TemplateID:  notificationType.TemplateID,
		SenderID:    notificationType.SenderID,
	}, nil
}
