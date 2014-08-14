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
    Update(ConnectionInterface, Client) (Client, error)
    Upsert(ConnectionInterface, Client) (Client, error)
}

func NewClientsRepo() ClientsRepo {
    return ClientsRepo{}
}

func (repo ClientsRepo) Create(conn ConnectionInterface, client Client) (Client, error) {
    client.CreatedAt = time.Now().Truncate(1 * time.Second).UTC()
    err := conn.Insert(&client)
    if err != nil {
        if strings.Contains(err.Error(), "Duplicate entry") {
            err = ErrDuplicateRecord{}
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
            err = ErrRecordNotFound{}
        }
        return client, err
    }
    return client, nil
}

func (repo ClientsRepo) Update(conn ConnectionInterface, client Client) (Client, error) {
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

    if err != nil {
        if (err == ErrRecordNotFound{}) {
            return repo.Create(conn, client)
        } else {
            return client, err
        }
    }

    return repo.Update(conn, client)
}
