package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type FakeFinder struct {
    Clients            map[string]models.Client
    Kinds              map[string]models.Kind
    ClientAndKindError error
}

func NewFakeFinder() *FakeFinder {
    return &FakeFinder{
        Clients: make(map[string]models.Client),
        Kinds:   make(map[string]models.Kind),
    }
}

func (finder *FakeFinder) ClientAndKind(clientID, kindID string) (models.Client, models.Kind, error) {
    return finder.Clients[clientID], finder.Kinds[kindID+"|"+clientID], finder.ClientAndKindError
}
