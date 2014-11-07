package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type Registrar struct {
    RegisterArguments []interface{}
    RegisterError     error
    PruneArguments    []interface{}
    PruneError        error
}

func NewRegistrar() *Registrar {
    return &Registrar{}
}

func (fake *Registrar) Register(conn models.ConnectionInterface, client models.Client, kinds []models.Kind) error {
    fake.RegisterArguments = []interface{}{conn, client, kinds}
    return fake.RegisterError
}

func (fake *Registrar) Prune(conn models.ConnectionInterface, client models.Client, kinds []models.Kind) error {
    fake.PruneArguments = []interface{}{conn, client, kinds}
    return fake.PruneError
}
