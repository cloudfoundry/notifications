package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type Registrar struct {
	RegisterCall struct {
		Arguments struct {
			Connection models.ConnectionInterface
			Client     models.Client
			Kinds      []models.Kind
		}
		Error error
	}

	PruneCall struct {
		Arguments struct {
			Connection models.ConnectionInterface
			Client     models.Client
			Kinds      []models.Kind
		}
		Called bool
		Error  error
	}
}

func NewRegistrar() *Registrar {
	return &Registrar{}
}

func (fake *Registrar) Register(conn models.ConnectionInterface, client models.Client, kinds []models.Kind) error {
	fake.RegisterCall.Arguments.Connection = conn
	fake.RegisterCall.Arguments.Client = client
	fake.RegisterCall.Arguments.Kinds = kinds

	return fake.RegisterCall.Error
}

func (fake *Registrar) Prune(conn models.ConnectionInterface, client models.Client, kinds []models.Kind) error {
	fake.PruneCall.Called = true
	fake.PruneCall.Arguments.Connection = conn
	fake.PruneCall.Arguments.Client = client
	fake.PruneCall.Arguments.Kinds = kinds

	return fake.PruneCall.Error
}
