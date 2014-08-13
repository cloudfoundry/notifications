package handlers

import "github.com/cloudfoundry-incubator/notifications/models"

type Finder struct {
    clientsRepo models.ClientsRepoInterface
    kindsRepo   models.KindsRepoInterface
}

type FinderInterface interface {
    ClientAndKind(string, string) (models.Client, models.Kind, error)
}

func NewFinder(clientsRepo models.ClientsRepoInterface, kindsRepo models.KindsRepoInterface) Finder {
    return Finder{
        clientsRepo: clientsRepo,
        kindsRepo:   kindsRepo,
    }
}

func (finder Finder) ClientAndKind(clientID, kindID string) (models.Client, models.Kind, error) {
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

func (finder Finder) client(clientID string) (models.Client, error) {
    client, err := finder.clientsRepo.Find(models.Database().Connection, clientID)
    if err != nil {
        if _, ok := err.(models.ErrRecordNotFound); ok {
            return models.Client{}, nil
        } else {
            return models.Client{}, err
        }
    }
    return client, nil
}

func (finder Finder) kind(clientID, kindID string) (models.Kind, error) {
    kind, err := finder.kindsRepo.Find(models.Database().Connection, kindID, clientID)
    if err != nil {
        if _, ok := err.(models.ErrRecordNotFound); ok {
            return models.Kind{}, nil
        } else {
            return models.Kind{}, err
        }
    }
    return kind, nil
}
