package models

import (
	"database/sql"
	"time"
)

type GlobalUnsubscribesRepo struct{}

func NewGlobalUnsubscribesRepo() GlobalUnsubscribesRepo {
	return GlobalUnsubscribesRepo{}
}

func (repo GlobalUnsubscribesRepo) Set(conn ConnectionInterface, userGUID string, unsubscribe bool) error {
	globalUnsubscribe, err := repo.find(conn, userGUID)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}

		globalUnsubscribe = GlobalUnsubscribe{
			UserID:    userGUID,
			CreatedAt: time.Now(),
		}
	}

	switch {
	case unsubscribe && globalUnsubscribe.Primary == 0:
		err = conn.Insert(&globalUnsubscribe)
		if err != nil {
			return err
		}
	case !unsubscribe && globalUnsubscribe.Primary != 0:
		_, err = conn.Delete(&globalUnsubscribe)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repo GlobalUnsubscribesRepo) Get(conn ConnectionInterface, userGUID string) (bool, error) {
	_, err := repo.find(conn, userGUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (repo GlobalUnsubscribesRepo) find(conn ConnectionInterface, userGUID string) (GlobalUnsubscribe, error) {
	globalUnsubscribe := GlobalUnsubscribe{}
	err := conn.SelectOne(&globalUnsubscribe, "SELECT * FROM `global_unsubscribes` WHERE `user_id` = ?", userGUID)
	if err != nil {
		return GlobalUnsubscribe{}, err
	}

	return globalUnsubscribe, nil
}
