package fakes

import "github.com/cloudfoundry-incubator/notifications/v1/models"

type ClientsRepo struct {
	Clients    map[string]models.Client
	AllClients []models.Client

	UpsertCall struct {
		Error error
	}

	FindCall struct {
		Error error
	}

	UpdateCall struct {
		Error error
	}

	FindAllByTemplateIDCall struct {
		Error error
	}
}

func NewClientsRepo() *ClientsRepo {
	return &ClientsRepo{
		Clients: make(map[string]models.Client),
	}
}

func (fake *ClientsRepo) Create(conn models.ConnectionInterface, client models.Client) (models.Client, error) {
	if _, ok := fake.Clients[client.ID]; ok {
		return client, models.DuplicateRecordError{}
	}
	if client.TemplateID == "" {
		client.TemplateID = models.DefaultTemplateID
	}

	fake.Clients[client.ID] = client
	return client, nil
}

func (fake *ClientsRepo) Update(conn models.ConnectionInterface, client models.Client) (models.Client, error) {
	if client.TemplateID == "" {
		existingClient, err := fake.Find(conn, client.ID)
		if err != nil {
			return client, err
		}
		client.TemplateID = existingClient.TemplateID
	}

	fake.Clients[client.ID] = client
	return client, fake.UpdateCall.Error
}

func (fake *ClientsRepo) Upsert(conn models.ConnectionInterface, client models.Client) (models.Client, error) {
	fake.Clients[client.ID] = client
	return client, fake.UpsertCall.Error
}

func (fake *ClientsRepo) Find(conn models.ConnectionInterface, id string) (models.Client, error) {
	if fake.FindCall.Error != nil {
		return models.Client{}, fake.FindCall.Error
	}

	if client, ok := fake.Clients[id]; ok {
		return client, nil
	}

	return models.Client{}, models.NewRecordNotFoundError("Client %q could not be found", id)
}

func (fake *ClientsRepo) FindAll(conn models.ConnectionInterface) ([]models.Client, error) {
	return fake.AllClients, nil
}

func (fake *ClientsRepo) FindAllByTemplateID(conn models.ConnectionInterface, templateID string) ([]models.Client, error) {
	var clients []models.Client
	for _, client := range fake.Clients {
		if client.TemplateID == templateID {
			clients = append(clients, client)
		}
	}
	return clients, fake.FindAllByTemplateIDCall.Error
}
