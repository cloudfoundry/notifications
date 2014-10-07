package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type FakeClientsRepo struct {
    Clients     map[string]models.Client
    UpsertError error
    FindError   error
}

func NewFakeClientsRepo() *FakeClientsRepo {
    return &FakeClientsRepo{
        Clients: make(map[string]models.Client),
    }
}

func (fake *FakeClientsRepo) Create(conn models.ConnectionInterface, client models.Client) (models.Client, error) {
    if _, ok := fake.Clients[client.ID]; ok {
        return client, models.ErrDuplicateRecord{}
    }
    fake.Clients[client.ID] = client
    return client, nil
}

func (fake *FakeClientsRepo) Update(conn models.ConnectionInterface, client models.Client) (models.Client, error) {
    fake.Clients[client.ID] = client
    return client, nil
}

func (fake *FakeClientsRepo) Upsert(conn models.ConnectionInterface, client models.Client) (models.Client, error) {
    fake.Clients[client.ID] = client
    return client, fake.UpsertError
}

func (fake *FakeClientsRepo) Find(conn models.ConnectionInterface, id string) (models.Client, error) {
    if client, ok := fake.Clients[id]; ok {
        return client, fake.FindError
    }
    return models.Client{}, models.ErrRecordNotFound{}
}
