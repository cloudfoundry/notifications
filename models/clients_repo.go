package models

import (
    "database/sql"
    "strings"
    "time"
)

type ClientsRepo struct{}

func NewClientsRepo() ClientsRepo {
    return ClientsRepo{}
}

func (repo ClientsRepo) Create(client Client) (Client, error) {
    client.CreatedAt = time.Now().Truncate(1 * time.Second).UTC()
    err := Database().Connection.Insert(&client)
    if err != nil {
        if strings.Contains(err.Error(), "Duplicate entry") {
            err = ErrDuplicateRecord
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
            err = ErrRecordNotFound
        }
        return client, err
    }
    return client, nil
}
