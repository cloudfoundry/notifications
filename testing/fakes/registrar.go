package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type Registrar struct {
	RegisterCall struct {
		Receives struct {
			Connection models.ConnectionInterface
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
			Connection models.ConnectionInterface
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

func (r *Registrar) Register(conn models.ConnectionInterface, client models.Client, kinds []models.Kind) error {
	r.RegisterCall.Receives.Connection = conn
	r.RegisterCall.Receives.Client = client
	r.RegisterCall.Receives.Kinds = kinds

	return r.RegisterCall.Returns.Error
}

func (r *Registrar) Prune(conn models.ConnectionInterface, client models.Client, kinds []models.Kind) error {
	r.PruneCall.Called = true
	r.PruneCall.Receives.Connection = conn
	r.PruneCall.Receives.Client = client
	r.PruneCall.Receives.Kinds = kinds

	return r.PruneCall.Returns.Error
}
