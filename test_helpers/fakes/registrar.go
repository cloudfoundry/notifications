package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type FakeRegistrar struct {
    RegisterArguments []interface{}
    RegisterError     error
    PruneArguments    []interface{}
    PruneError        error
}

func NewFakeRegistrar() *FakeRegistrar {
    return &FakeRegistrar{}
}

func (fake *FakeRegistrar) Register(conn models.ConnectionInterface, client models.Client, kinds []models.Kind) error {
    fake.RegisterArguments = []interface{}{conn, client, kinds}
    return fake.RegisterError
}

func (fake *FakeRegistrar) Prune(conn models.ConnectionInterface, client models.Client, kinds []models.Kind) error {
    fake.PruneArguments = []interface{}{conn, client, kinds}
    return fake.PruneError
}
