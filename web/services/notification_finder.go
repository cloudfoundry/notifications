package services

import "github.com/cloudfoundry-incubator/notifications/models"

type NotificationFinder struct {
    clientsRepo models.ClientsRepoInterface
    kindsRepo   models.KindsRepoInterface
}

type NotificationFinderInterface interface {
    ClientAndKind(string, string) (models.Client, models.Kind, error)
}

func NewNotificationFinder(clientsRepo models.ClientsRepoInterface, kindsRepo models.KindsRepoInterface) NotificationFinder {
    return NotificationFinder{
        clientsRepo: clientsRepo,
        kindsRepo:   kindsRepo,
    }
}

func (finder NotificationFinder) ClientAndKind(clientID, kindID string) (models.Client, models.Kind, error) {
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

func (finder NotificationFinder) client(clientID string) (models.Client, error) {
    client, err := finder.clientsRepo.Find(models.Database().Connection, clientID)
    if err != nil {
        if _, ok := err.(models.ErrRecordNotFound); ok {
            return models.Client{ID: clientID}, nil
        } else {
            return models.Client{}, err
        }
    }
    return client, nil
}

func (finder NotificationFinder) kind(clientID, kindID string) (models.Kind, error) {
    kind, err := finder.kindsRepo.Find(models.Database().Connection, kindID, clientID)
    if err != nil {
        if _, ok := err.(models.ErrRecordNotFound); ok {
            return models.Kind{ID: kindID, ClientID: clientID}, nil
        } else {
            return models.Kind{}, err
        }
    }
    return kind, nil
}
