package models

import (
	"database/sql"
	"fmt"
	"strings"
)

type SendersRepository struct {
	generateGUID guidGeneratorFunc
}

type guidGeneratorFunc func() (string, error)

type Sender struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	ClientID string `db:"client_id"`
}

func NewSendersRepository(guidGenerator guidGeneratorFunc) SendersRepository {
	return SendersRepository{
		generateGUID: guidGenerator,
	}
}

func (r SendersRepository) Insert(conn ConnectionInterface, sender Sender) (Sender, error) {
	var err error
	sender.ID, err = r.generateGUID()
	if err != nil {
		return Sender{}, err
	}

	err = conn.Insert(&sender)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return Sender{}, DuplicateRecordError{fmt.Errorf("Sender with name %q already exists", sender.Name)}
		}

		return Sender{}, err
	}

	return sender, nil
}

func (r SendersRepository) Update(conn ConnectionInterface, sender Sender) (Sender, error) {
	_, err := conn.Update(&sender)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			err = DuplicateRecordError{fmt.Errorf("Sender with name %q already exists", sender.Name)}
		}
		return sender, err
	}

	return sender, nil
}

func (r SendersRepository) List(conn ConnectionInterface, clientID string) ([]Sender, error) {
	senders := []Sender{}
	_, err := conn.Select(&senders, "SELECT * FROM `senders` WHERE `client_id` = ?", clientID)
	return senders, err
}

func (r SendersRepository) Get(conn ConnectionInterface, senderID string) (Sender, error) {
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

func (r SendersRepository) GetByClientIDAndName(conn ConnectionInterface, clientID, name string) (Sender, error) {
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

func (r SendersRepository) Delete(conn ConnectionInterface, sender Sender) error {
	_, err := conn.Delete(&sender)
	if err != nil {
		return err
	}

	return nil
}
