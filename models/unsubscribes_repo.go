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
	Destroy(ConnectionInterface, Unsubscribe) (int, error)
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

func (repo UnsubscribesRepo) Upsert(conn ConnectionInterface, unsubscribe Unsubscribe) (Unsubscribe, error) {
	_, err := repo.Find(conn, unsubscribe.ClientID, unsubscribe.KindID, unsubscribe.UserID)
	if err != nil {
		if (err == ErrRecordNotFound{}) {
			return repo.Create(conn, unsubscribe)
		}
		return unsubscribe, err
	}
	return unsubscribe, nil
}

func (repo UnsubscribesRepo) FindAllByUserID(conn ConnectionInterface, userID string) ([]Unsubscribe, error) {
	unsubscribes := []Unsubscribe{}
	results, err := conn.Select(Unsubscribe{}, "SELECT * FROM `unsubscribes` WHERE `user_id` = ?", userID)
	if err != nil {
		return unsubscribes, err
	}

	for _, result := range results {
		unsubscribes = append(unsubscribes, *(result.(*Unsubscribe)))
	}

	return unsubscribes, nil
}

func (repo UnsubscribesRepo) Destroy(conn ConnectionInterface, unsubscribe Unsubscribe) (int, error) {
	unsubscribe, err := repo.Find(conn, unsubscribe.ClientID, unsubscribe.KindID, unsubscribe.UserID)
	if err != nil {
		if (err == ErrRecordNotFound{}) {
			return 0, nil
		}
		return 0, err
	}
	rowsAffected, err := conn.Delete(&unsubscribe)
	return int(rowsAffected), err
}
