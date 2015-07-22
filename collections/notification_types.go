package collections

import (
	"errors"

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
}

func NewNotificationTypesCollection(nr notificationTypesRepository, sr sendersRepository) NotificationTypesCollection {
	return NotificationTypesCollection{
		notificationTypesRepository: nr,
		sendersRepository:           sr,
	}
}

func (nc NotificationTypesCollection) Add(conn models.ConnectionInterface, notificationType NotificationType) (NotificationType, error) {
	if notificationType.Name == "" {
		return NotificationType{}, ValidationError{
			Err: errors.New("missing notification type name"),
		}
	}

	if notificationType.Description == "" {
		return NotificationType{}, ValidationError{
			Err: errors.New("missing notification type description"),
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
	senderModel, err := nc.sendersRepository.Get(conn, senderID)

	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return []NotificationType{}, NotFoundError{err}
		default:
			return []NotificationType{}, PersistenceError{err}
		}
	}

	if senderModel.ClientID != clientID {
		return []NotificationType{}, NotFoundError{errors.New("sender not found")}
	}

	modelList, err := nc.notificationTypesRepository.List(conn, senderID)
	if err != nil {
		panic(err)
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
