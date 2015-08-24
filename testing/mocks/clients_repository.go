package mocks

import "github.com/cloudfoundry-incubator/notifications/v1/models"

type ClientsRepository struct {
	FindCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			ClientID   string
		}
		Returns struct {
			Client models.Client
			Error  error
		}
	}

	FindAllCall struct {
		Receives struct {
			Connection models.ConnectionInterface
		}
		Returns struct {
			Clients []models.Client
			Error   error
		}
	}

	FindAllByTemplateIDCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			TemplateID string
		}
		Returns struct {
			Clients []models.Client
			Error   error
		}
	}

	UpdateCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			Client     models.Client
		}
		Returns struct {
			Client models.Client
			Error  error
		}
	}

	UpsertCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			Client     models.Client
		}
		Returns struct {
			Client models.Client
			Error  error
		}
	}
}

func NewClientsRepository() *ClientsRepository {
	return &ClientsRepository{}
}

func (cr *ClientsRepository) Find(conn models.ConnectionInterface, clientID string) (models.Client, error) {
	cr.FindCall.Receives.Connection = conn
	cr.FindCall.Receives.ClientID = clientID

	return cr.FindCall.Returns.Client, cr.FindCall.Returns.Error
}

func (cr *ClientsRepository) FindAll(conn models.ConnectionInterface) ([]models.Client, error) {
	cr.FindAllCall.Receives.Connection = conn

	return cr.FindAllCall.Returns.Clients, cr.FindAllCall.Returns.Error
}

func (cr *ClientsRepository) FindAllByTemplateID(conn models.ConnectionInterface, templateID string) ([]models.Client, error) {
	cr.FindAllByTemplateIDCall.Receives.Connection = conn
	cr.FindAllByTemplateIDCall.Receives.TemplateID = templateID

	return cr.FindAllByTemplateIDCall.Returns.Clients, cr.FindAllByTemplateIDCall.Returns.Error
}

func (cr *ClientsRepository) Update(conn models.ConnectionInterface, client models.Client) (models.Client, error) {
	cr.UpdateCall.Receives.Connection = conn
	cr.UpdateCall.Receives.Client = client

	return cr.UpdateCall.Returns.Client, cr.UpdateCall.Returns.Error
}

func (cr *ClientsRepository) Upsert(conn models.ConnectionInterface, client models.Client) (models.Client, error) {
	cr.UpsertCall.Receives.Connection = conn
	cr.UpsertCall.Receives.Client = client

	return cr.UpsertCall.Returns.Client, cr.UpsertCall.Returns.Error
}
