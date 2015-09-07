package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
)

type Registrar struct {
	RegisterCall struct {
		Receives struct {
			Connection services.ConnectionInterface
			Client     models.Client
			Kinds      []models.Kind
		}
		Returns struct {
			Error error
		}
	}

	PruneCall struct {
		Called   bool
		Receives struct {
			Connection services.ConnectionInterface
			Client     models.Client
			Kinds      []models.Kind
		}
		Returns struct {
			Error error
		}
	}
}

func NewRegistrar() *Registrar {
	return &Registrar{}
}

func (r *Registrar) Register(conn services.ConnectionInterface, client models.Client, kinds []models.Kind) error {
	r.RegisterCall.Receives.Connection = conn
	r.RegisterCall.Receives.Client = client
	r.RegisterCall.Receives.Kinds = kinds

	return r.RegisterCall.Returns.Error
}

func (r *Registrar) Prune(conn services.ConnectionInterface, client models.Client, kinds []models.Kind) error {
	r.PruneCall.Called = true
	r.PruneCall.Receives.Connection = conn
	r.PruneCall.Receives.Client = client
	r.PruneCall.Receives.Kinds = kinds

	return r.PruneCall.Returns.Error
}
