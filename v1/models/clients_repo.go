package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type ClientsRepo struct{}

func NewClientsRepo() ClientsRepo {
	return ClientsRepo{}
}

func (repo ClientsRepo) Find(conn ConnectionInterface, id string) (Client, error) {
	client := Client{}
	err := conn.SelectOne(&client, "SELECT * FROM `clients` WHERE `id` = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = NotFoundError{fmt.Errorf("Client with ID %q could not be found", id)}
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
	case NotFoundError:
		client, err := repo.create(conn, client)
		if _, ok := err.(DuplicateError); ok {
			return repo.Update(conn, client)
		}

		return client, err
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

func (repo ClientsRepo) create(conn ConnectionInterface, client Client) (Client, error) {
	err := conn.Insert(&client)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			err = DuplicateError{errors.New("duplicate record")}
		}
		return client, err
	}
	return client, nil
}
