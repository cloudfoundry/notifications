package models

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/nu7hatch/gouuid"
)

type SendersRepository struct {
	generateGUID guidGeneratorFunc
}

type guidGeneratorFunc func() (*uuid.UUID, error)

func NewSendersRepository(guidGenerator guidGeneratorFunc) SendersRepository {
	return SendersRepository{
		generateGUID: guidGenerator,
	}
}

func (r SendersRepository) Insert(conn db.ConnectionInterface, sender Sender) (Sender, error) {
	guid, err := r.generateGUID()
	if err != nil {
		panic(err)
	}

	sender.ID = guid.String()
	err = conn.Insert(&sender)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			err = DuplicateRecordError{}
		}
		return sender, err
	}

	return sender, nil
}

func (r SendersRepository) Get(conn db.ConnectionInterface, senderID string) (Sender, error) {
	sender := Sender{}
	err := conn.SelectOne(&sender, "SELECT * FROM `senders` WHERE `id` = ?", senderID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = RecordNotFoundError{fmt.Errorf("Sender with id %q could not be found", senderID)}
		}
		return sender, err
	}

	return sender, nil
}

func (r SendersRepository) GetByClientIDAndName(conn db.ConnectionInterface, clientID, name string) (Sender, error) {
	sender := Sender{}
	err := conn.SelectOne(&sender, "SELECT * FROM `senders` WHERE `client_id` = ? AND `name` = ?", clientID, name)
	if err != nil {
		if err == sql.ErrNoRows {
			err = RecordNotFoundError{fmt.Errorf("Sender with client_id %q and name %q could not be found", clientID, name)}
		}
		return sender, err
	}

	return sender, nil
}
