package fakes

import (
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/models"
)

type MessagesRepo struct {
	Messages      map[string]models.Message
	FindByIDError error
	UpsertError   error
}

func NewMessagesRepo() *MessagesRepo {
	return &MessagesRepo{
		Messages: make(map[string]models.Message),
	}
}

func (fake MessagesRepo) FindByID(conn models.ConnectionInterface, messageID string) (models.Message, error) {
	if fake.FindByIDError != nil {
		return models.Message{}, fake.FindByIDError
	}

	message, ok := fake.Messages[messageID]
	if !ok {
		return message, models.RecordNotFoundError(fmt.Sprintf("We did not find the message with ID %s", messageID))
	}
	return message, fake.FindByIDError
}

func (fake MessagesRepo) Upsert(conn models.ConnectionInterface, message models.Message) (models.Message, error) {
	fake.Messages[message.ID] = message

	return message, fake.UpsertError
}
