package services

import "github.com/cloudfoundry-incubator/notifications/models"

type NotificationsFinder struct {
	clientsRepo clientsRepo
	kindsRepo   models.KindsRepoInterface
	database    models.DatabaseInterface
}

type clientsRepo interface {
	FindAll(models.ConnectionInterface) ([]models.Client, error)
	Find(models.ConnectionInterface, string) (models.Client, error)
	Upsert(models.ConnectionInterface, models.Client) (models.Client, error)
	Update(models.ConnectionInterface, models.Client) (models.Client, error)
	FindAllByTemplateID(models.ConnectionInterface, string) ([]models.Client, error)
}

type NotificationsFinderInterface interface {
	AllClientsAndNotifications(models.DatabaseInterface) ([]models.Client, []models.Kind, error)
	ClientAndKind(models.DatabaseInterface, string, string) (models.Client, models.Kind, error)
}

func NewNotificationsFinder(clientsRepo clientsRepo, kindsRepo models.KindsRepoInterface) NotificationsFinder {
	return NotificationsFinder{
		clientsRepo: clientsRepo,
		kindsRepo:   kindsRepo,
	}
}

func (finder NotificationsFinder) AllClientsAndNotifications(database models.DatabaseInterface) ([]models.Client, []models.Kind, error) {
	var clients []models.Client
	var notifications []models.Kind
	var err error

	clients, err = finder.clientsRepo.FindAll(database.Connection())
	if err != nil {
		return clients, notifications, err
	}

	notifications, err = finder.kindsRepo.FindAll(database.Connection())
	if err != nil {
		return clients, notifications, err
	}

	return clients, notifications, nil
}

func (finder NotificationsFinder) ClientAndKind(database models.DatabaseInterface, clientID, kindID string) (models.Client, models.Kind, error) {
	client, err := finder.client(database, clientID)
	if err != nil {
		return models.Client{}, models.Kind{}, err
	}

	kind, err := finder.kind(database, clientID, kindID)
	if err != nil {
		return client, models.Kind{}, err
	}

	return client, kind, nil
}

func (finder NotificationsFinder) client(database models.DatabaseInterface, clientID string) (models.Client, error) {
	client, err := finder.clientsRepo.Find(database.Connection(), clientID)
	if err != nil {
		if _, ok := err.(models.RecordNotFoundError); ok {
			return models.Client{ID: clientID}, nil
		} else {
			return models.Client{}, err
		}
	}
	return client, nil
}

func (finder NotificationsFinder) kind(database models.DatabaseInterface, clientID, kindID string) (models.Kind, error) {
	kind, err := finder.kindsRepo.Find(database.Connection(), kindID, clientID)
	if err != nil {
		if _, ok := err.(models.RecordNotFoundError); ok {
			return models.Kind{ID: kindID, ClientID: clientID}, nil
		} else {
			return models.Kind{}, err
		}
	}
	return kind, nil
}
