package models

import (
    "database/sql"
    "strings"
    "time"
)

type UnsubscribesRepoInterface interface {
    Create(ConnectionInterface, Unsubscribe) (Unsubscribe, error)
    Upsert(ConnectionInterface, Unsubscribe) (Unsubscribe, error)
    Find(ConnectionInterface, string, string, string) (Unsubscribe, error)
}

type UnsubscribesRepo struct{}

func NewUnsubscribesRepo() UnsubscribesRepo {
    return UnsubscribesRepo{}
}

func (repo UnsubscribesRepo) Create(conn ConnectionInterface, unsubscribe Unsubscribe) (Unsubscribe, error) {
    unsubscribe.CreatedAt = time.Now().Truncate(1 * time.Second).UTC()
    err := conn.Insert(&unsubscribe)
    if err != nil {
        if strings.Contains(err.Error(), "Duplicate entry") {
            err = ErrDuplicateRecord{}
        }
        return unsubscribe, err
    }
    return unsubscribe, nil
}

//TODO: this is not needed until resubscribing to notifications story 76994722
func (repo UnsubscribesRepo) Upsert(conn ConnectionInterface, unsubscribe Unsubscribe) (Unsubscribe, error) {
    return unsubscribe, nil
}

func (repo UnsubscribesRepo) Find(conn ConnectionInterface, clientID string, kindID string, userID string) (Unsubscribe, error) {
    unsubscribe := Unsubscribe{}
    err := conn.SelectOne(&unsubscribe, "SELECT * FROM `unsubscribes` WHERE `client_id` = ? AND `kind_id` = ? AND `user_id` = ?", clientID, kindID, userID)
    if err != nil {
        if err == sql.ErrNoRows {
            err = ErrRecordNotFound{}
        }
        return unsubscribe, err
    }
    return unsubscribe, nil
}
