package fakes

import (
	"fmt"
	"time"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/models"
)

type MessagesRepo struct {
	Messages                map[string]models.Message
	DeleteBeforeError       error
	FindByIDError           error
	UpsertError             error
	DeleteBeforeInvocations []time.Time
}

func NewMessagesRepo() *MessagesRepo {
	return &MessagesRepo{
		Messages:                make(map[string]models.Message),
		DeleteBeforeInvocations: []time.Time{},
	}
}

func (fake MessagesRepo) FindByID(conn db.ConnectionInterface, messageID string) (models.Message, error) {
	if fake.FindByIDError != nil {
		return models.Message{}, fake.FindByIDError
	}

	message, ok := fake.Messages[messageID]
	if !ok {
		return message, models.RecordNotFoundError(fmt.Sprintf("We did not find the message with ID %s", messageID))
	}
	return message, fake.FindByIDError
}

func (fake MessagesRepo) Upsert(conn db.ConnectionInterface, message models.Message) (models.Message, error) {
	fake.Messages[message.ID] = message

	return message, fake.UpsertError
}

func (fake *MessagesRepo) DeleteBefore(conn db.ConnectionInterface, thresholdTime time.Time) (int, error) {
	count := 0
	for key, message := range fake.Messages {
		if message.UpdatedAt.Before(thresholdTime) {
			delete(fake.Messages, key)
			count += 1
		}
	}
	fake.DeleteBeforeInvocations = append(fake.DeleteBeforeInvocations, time.Now())
	return count, fake.DeleteBeforeError
}
