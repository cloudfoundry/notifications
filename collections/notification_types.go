package collections

import (
	"errors"
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/models"
)

type NotificationType struct {
	ID          string
	Name        string
	Description string
	Critical    bool
	TemplateID  string
	SenderID    string
}

type NotificationTypesCollection struct {
	notificationTypesRepository notificationTypesRepository
	sendersRepository           sendersRepository
}

type notificationTypesRepository interface {
	Insert(models.ConnectionInterface, models.NotificationType) (models.NotificationType, error)
	GetBySenderIDAndName(models.ConnectionInterface, string, string) (models.NotificationType, error)
	List(models.ConnectionInterface, string) ([]models.NotificationType, error)
	Get(models.ConnectionInterface, string) (models.NotificationType, error)
}

func NewNotificationTypesCollection(nr notificationTypesRepository, sr sendersRepository) NotificationTypesCollection {
	return NotificationTypesCollection{
		notificationTypesRepository: nr,
		sendersRepository:           sr,
	}
}

func (nc NotificationTypesCollection) Add(conn models.ConnectionInterface, notificationType NotificationType, clientID string) (NotificationType, error) {
	senderModel, err := nc.sendersRepository.Get(conn, notificationType.SenderID)

	if err != nil {
		switch e := err.(type) {
		case models.RecordNotFoundError:
			return NotificationType{}, NotFoundError{
				Err:     e,
				Message: string(e),
			}
		default:
			return NotificationType{}, PersistenceError{err}
		}
	}

	if senderModel.ClientID != clientID {
		return NotificationType{}, NotFoundError{
			Err:     errors.New("sender not found"),
			Message: "sender not found",
		}
	}

	if notificationType.Name == "" {
		return NotificationType{}, ValidationError{
			Err:     errors.New("missing campaign type name"),
			Message: "missing campaign type name",
		}
	}

	if notificationType.Description == "" {
		return NotificationType{}, ValidationError{
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
				return NotificationType{}, PersistenceError{err}
			}
		default:
			return NotificationType{}, PersistenceError{err}
		}
	}

	return NotificationType{
		ID:          returnNotificationType.ID,
		Name:        returnNotificationType.Name,
		Description: returnNotificationType.Description,
		Critical:    returnNotificationType.Critical,
		TemplateID:  returnNotificationType.TemplateID,
		SenderID:    returnNotificationType.SenderID,
	}, err
}

func (nc NotificationTypesCollection) List(conn models.ConnectionInterface, senderID, clientID string) ([]NotificationType, error) {
	if senderID == "" {
		return []NotificationType{}, ValidationError{
			Err:     errors.New("missing sender id"),
			Message: "missing sender id",
		}
	}

	if clientID == "" {
		return []NotificationType{}, NewValidationError("missing client id")
	}

	senderModel, err := nc.sendersRepository.Get(conn, senderID)

	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return []NotificationType{}, NotFoundError{
				Err:     err,
				Message: "sender not found",
			}
		default:
			return []NotificationType{}, PersistenceError{err}
		}
	}

	if senderModel.ClientID != clientID {
		return []NotificationType{}, NewNotFoundError("sender not found")
	}

	modelList, err := nc.notificationTypesRepository.List(conn, senderID)
	if err != nil {
		return []NotificationType{}, PersistenceError{err}
	}

	notificationTypeList := []NotificationType{}

	for _, model := range modelList {
		notificationType := NotificationType{
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

func (nc NotificationTypesCollection) Get(conn models.ConnectionInterface, notificationTypeID, senderID, clientID string) (NotificationType, error) {
	if notificationTypeID == "" {
		return NotificationType{}, ValidationError{
			Err:     errors.New("missing campaign type id"),
			Message: "missing campaign type id",
		}
	}

	if senderID == "" {
		return NotificationType{}, ValidationError{
			Err:     errors.New("missing sender id"),
			Message: "missing sender id",
		}
	}

	if clientID == "" {
		return NotificationType{}, ValidationError{
			Err:     errors.New("missing client id"),
			Message: "missing client id",
		}
	}

	sender, err := nc.sendersRepository.Get(conn, senderID)
	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return NotificationType{}, NewNotFoundError(fmt.Sprintf("sender %s not found", notificationTypeID))
		default:
			return NotificationType{}, PersistenceError{err}
		}
	}

	if clientID != sender.ClientID {
		message := fmt.Sprintf("sender %s not found", senderID)
		return NotificationType{}, NotFoundError{
			Err:     errors.New(message),
			Message: message,
		}
	}

	notificationType, err := nc.notificationTypesRepository.Get(conn, notificationTypeID)
	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return NotificationType{}, NotFoundError{
				Err:     err,
				Message: fmt.Sprintf("campaign type %s not found", notificationTypeID),
			}
		default:
			return NotificationType{}, PersistenceError{err}
		}
	}

	if senderID != notificationType.SenderID {
		message := fmt.Sprintf("campaign type %s not found", notificationTypeID)
		return NotificationType{}, NotFoundError{
			Err:     errors.New(message),
			Message: message,
		}
	}

	return NotificationType{
		ID:          notificationType.ID,
		Name:        notificationType.Name,
		Description: notificationType.Description,
		Critical:    notificationType.Critical,
		TemplateID:  notificationType.TemplateID,
		SenderID:    notificationType.SenderID,
	}, nil
}
