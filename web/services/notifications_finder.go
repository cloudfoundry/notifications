package services

import "github.com/cloudfoundry-incubator/notifications/models"

type NotificationsFinder struct {
	clientsRepo models.ClientsRepoInterface
	kindsRepo   models.KindsRepoInterface
	database    models.DatabaseInterface
}

type NotificationsFinderInterface interface {
	AllClientsAndNotifications() ([]models.Client, []models.Kind, error)
	ClientAndKind(string, string) (models.Client, models.Kind, error)
}

func NewNotificationsFinder(clientsRepo models.ClientsRepoInterface, kindsRepo models.KindsRepoInterface, database models.DatabaseInterface) NotificationsFinder {
	return NotificationsFinder{
		clientsRepo: clientsRepo,
		kindsRepo:   kindsRepo,
		database:    database,
	}
}

func (finder NotificationsFinder) AllClientsAndNotifications() ([]models.Client, []models.Kind, error) {
	var clients []models.Client
	var notifications []models.Kind
	var err error

	clients, err = finder.clientsRepo.FindAll(finder.database.Connection())
	if err != nil {
		return clients, notifications, err
	}

	notifications, err = finder.kindsRepo.FindAll(finder.database.Connection())
	if err != nil {
		return clients, notifications, err
	}

	return clients, notifications, nil
}

func (finder NotificationsFinder) ClientAndKind(clientID, kindID string) (models.Client, models.Kind, error) {
	client, err := finder.client(clientID)
	if err != nil {
		return models.Client{}, models.Kind{}, err
	}

	kind, err := finder.kind(clientID, kindID)
	if err != nil {
		return client, models.Kind{}, err
	}

	return client, kind, nil
}

func (finder NotificationsFinder) client(clientID string) (models.Client, error) {
	client, err := finder.clientsRepo.Find(finder.database.Connection(), clientID)
	if err != nil {
		if _, ok := err.(models.ErrRecordNotFound); ok {
			return models.Client{ID: clientID}, nil
		} else {
			return models.Client{}, err
		}
	}
	return client, nil
}

func (finder NotificationsFinder) kind(clientID, kindID string) (models.Kind, error) {
	kind, err := finder.kindsRepo.Find(finder.database.Connection(), kindID, clientID)
	if err != nil {
		if _, ok := err.(models.ErrRecordNotFound); ok {
			return models.Kind{ID: kindID, ClientID: clientID}, nil
		} else {
			return models.Kind{}, err
		}
	}
	return kind, nil
}
