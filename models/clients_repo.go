package models

import (
    "database/sql"
    "strings"
    "time"
)

type ClientsRepo struct{}

type ClientsRepoInterface interface {
    Create(Client) (Client, error)
    Find(string) (Client, error)
    Update(Client) (Client, error)
    Upsert(Client) (Client, error)
}

func NewClientsRepo() ClientsRepo {
    return ClientsRepo{}
}

func (repo ClientsRepo) Create(client Client) (Client, error) {
    client.CreatedAt = time.Now().Truncate(1 * time.Second).UTC()
    err := Database().Connection.Insert(&client)
    if err != nil {
        if strings.Contains(err.Error(), "Duplicate entry") {
            err = ErrDuplicateRecord{}
        }
        return client, err
    }
    return client, nil
}

func (repo ClientsRepo) Find(id string) (Client, error) {
    client := Client{}
    err := Database().Connection.SelectOne(&client, "SELECT * FROM `clients` WHERE `id` = ?", id)
    if err != nil {
        if err == sql.ErrNoRows {
            err = ErrRecordNotFound{}
        }
        return client, err
    }
    return client, nil
}

func (repo ClientsRepo) Update(client Client) (Client, error) {
    _, err := Database().Connection.Update(&client)
    if err != nil {
        return client, err
    }

    return repo.Find(client.ID)
}

func (repo ClientsRepo) Upsert(client Client) (Client, error) {
    existingClient, err := repo.Find(client.ID)
    client.CreatedAt = existingClient.CreatedAt

    if err != nil {
        if (err == ErrRecordNotFound{}) {
            return repo.Create(client)
        } else {
            return client, err
        }
    }

    return repo.Update(client)
}
