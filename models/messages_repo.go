package models

import (
	"database/sql"
	"time"

	"github.com/cloudfoundry-incubator/notifications/db"
)

type MessagesRepo struct {
}

func NewMessagesRepo() MessagesRepo {
	return MessagesRepo{}
}

func (repo MessagesRepo) Create(conn db.ConnectionInterface, message Message) (Message, error) {
	err := conn.Insert(&message)
	if err != nil {
		return Message{}, err
	}
	return message, nil
}

func (repo MessagesRepo) FindByID(conn db.ConnectionInterface, messageID string) (Message, error) {
	message := Message{}
	err := conn.SelectOne(&message, "SELECT * FROM `messages` WHERE `id`=?", messageID)
	if err != nil {
		if err == sql.ErrNoRows {
			return Message{}, NewRecordNotFoundError("Message with ID %q could not be found", messageID)
		}
		return Message{}, err
	}
	return message, nil
}

func (repo MessagesRepo) Update(conn db.ConnectionInterface, message Message) (Message, error) {
	_, err := conn.Update(&message)
	if err != nil {
		return message, err
	}

	return repo.FindByID(conn, message.ID)
}

func (repo MessagesRepo) Upsert(conn db.ConnectionInterface, message Message) (Message, error) {
	_, err := repo.FindByID(conn, message.ID)

	switch err.(type) {
	case RecordNotFoundError:
		return repo.Create(conn, message)
	case nil:
		return repo.Update(conn, message)
	default:
		return message, err
	}
}

func (repo MessagesRepo) DeleteBefore(conn db.ConnectionInterface, threshold time.Time) (int, error) {
	result, err := conn.Exec("DELETE FROM `messages` WHERE `updated_at` < ?", threshold.UTC())
	if err != nil {
		return 0, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
