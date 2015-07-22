package collections

import "github.com/cloudfoundry-incubator/notifications/models"

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
}

type notificationTypesRepository interface {
	Insert(models.ConnectionInterface, models.NotificationType) (models.NotificationType, error)
	GetBySenderIDAndName(models.ConnectionInterface, string, string) (models.NotificationType, error)
}

func NewNotificationTypesCollection(repo notificationTypesRepository) NotificationTypesCollection {
	return NotificationTypesCollection{
		notificationTypesRepository: repo,
	}
}

func (nc NotificationTypesCollection) Add(conn models.ConnectionInterface, notificationType NotificationType) (NotificationType, error) {
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
