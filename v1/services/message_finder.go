package services

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
)

type Message struct {
	Status string
}

type MessagesRepoInterface interface {
	FindByID(db.ConnectionInterface, string) (models.Message, error)
	Upsert(db.ConnectionInterface, models.Message) (models.Message, error)
}

type MessageFinder struct {
	repo MessagesRepoInterface
}

func NewMessageFinder(repo MessagesRepoInterface) MessageFinder {
	return MessageFinder{
		repo: repo,
	}
}

func (finder MessageFinder) Find(database DatabaseInterface, messageID string) (Message, error) {
	message, err := finder.repo.FindByID(database.Connection(), messageID)
	if err != nil {
		return Message{}, err
	}

	return Message{Status: message.Status}, nil
}
