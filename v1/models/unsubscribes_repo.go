package models

import (
	"database/sql"
	"errors"
	"strings"
)

type UnsubscribesRepo struct{}

func NewUnsubscribesRepo() UnsubscribesRepo {
	return UnsubscribesRepo{}
}

func (repo UnsubscribesRepo) Get(conn ConnectionInterface, userID, clientID, kindID string) (bool, error) {
	err := conn.SelectOne(&Unsubscribe{}, "SELECT * FROM `unsubscribes` WHERE `client_id` = ? AND `kind_id` = ? AND `user_id` = ?", clientID, kindID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (repo UnsubscribesRepo) Set(conn ConnectionInterface, userID, clientID, kindID string, unsubscribe bool) error {
	var record Unsubscribe
	err := conn.SelectOne(&record, "SELECT * FROM `unsubscribes` WHERE `client_id` = ? AND `kind_id` = ? AND `user_id` = ?", clientID, kindID, userID)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}

		record = Unsubscribe{
			UserID:   userID,
			ClientID: clientID,
			KindID:   kindID,
		}
	}

	switch {
	case unsubscribe && record.Primary == 0:
		_, err = repo.create(conn, record)
		if err != nil {
			return err
		}

	case !unsubscribe && record.Primary != 0:
		_, err = repo.delete(conn, record)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repo UnsubscribesRepo) create(conn ConnectionInterface, unsubscribe Unsubscribe) (Unsubscribe, error) {
	err := conn.Insert(&unsubscribe)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			err = DuplicateError{errors.New("duplicate record")}
		}
		return unsubscribe, err
	}
	return unsubscribe, nil
}

func (repo UnsubscribesRepo) delete(conn ConnectionInterface, unsubscribe Unsubscribe) (int, error) {
	rowsAffected, err := conn.Delete(&unsubscribe)
	return int(rowsAffected), err
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
