package models

import (
	"database/sql"
	"strings"
	"time"
)

type ClientsRepo struct{}

type ClientsRepoInterface interface {
	Create(ConnectionInterface, Client) (Client, error)
	Find(ConnectionInterface, string) (Client, error)
	FindAll(ConnectionInterface) ([]Client, error)
	FindAllByTemplateID(ConnectionInterface, string) ([]Client, error)
	Update(ConnectionInterface, Client) (Client, error)
	Upsert(ConnectionInterface, Client) (Client, error)
}

func NewClientsRepo() ClientsRepo {
	return ClientsRepo{}
}

func (repo ClientsRepo) Create(conn ConnectionInterface, client Client) (Client, error) {
	client.CreatedAt = time.Now().Truncate(1 * time.Second).UTC()
	if client.TemplateID == "" {
		client.TemplateID = DefaultTemplateID
	}
	err := conn.Insert(&client)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			err = DuplicateRecordError{}
		}
		return client, err
	}
	return client, nil
}

func (repo ClientsRepo) Find(conn ConnectionInterface, id string) (Client, error) {
	client := Client{}
	err := conn.SelectOne(&client, "SELECT * FROM `clients` WHERE `id` = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = NewRecordNotFoundError("Client with ID %q could not be found", id)
		}
		return client, err
	}
	return client, nil
}

func (repo ClientsRepo) FindAll(conn ConnectionInterface) ([]Client, error) {
	clients := []Client{}
	_, err := conn.Select(&clients, "SELECT * FROM `clients`")
	if err != nil {
		return []Client{}, err
	}

	return clients, nil
}

func (repo ClientsRepo) Update(conn ConnectionInterface, client Client) (Client, error) {
	if client.TemplateID == DoNotSetTemplateID {
		existingClient, err := repo.Find(conn, client.ID)
		if err != nil {
			return client, err
		}
		client.TemplateID = existingClient.TemplateID
	}

	_, err := conn.Update(&client)
	if err != nil {
		return client, err
	}

	return repo.Find(conn, client.ID)
}

func (repo ClientsRepo) Upsert(conn ConnectionInterface, client Client) (Client, error) {
	existingClient, err := repo.Find(conn, client.ID)
	client.Primary = existingClient.Primary
	client.CreatedAt = existingClient.CreatedAt

	switch err.(type) {
	case RecordNotFoundError:
		return repo.Create(conn, client)
	case nil:
		return repo.Update(conn, client)
	default:
		return client, err
	}
}

func (repo ClientsRepo) FindAllByTemplateID(conn ConnectionInterface, templateID string) ([]Client, error) {
	clients := []Client{}
	_, err := conn.Select(&clients, "SELECT * FROM `clients` WHERE `template_id` = ?", templateID)
	if err != nil {
		return clients, err
	}

	return clients, nil
}
