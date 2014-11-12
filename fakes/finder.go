package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type Finder struct {
	Clients            map[string]models.Client
	Kinds              map[string]models.Kind
	ClientAndKindError error
}

func NewFinder() *Finder {
	return &Finder{
		Clients: make(map[string]models.Client),
		Kinds:   make(map[string]models.Kind),
	}
}

func (finder *Finder) ClientAndKind(clientID, kindID string) (models.Client, models.Kind, error) {
	return finder.Clients[clientID], finder.Kinds[kindID+"|"+clientID], finder.ClientAndKindError
}
