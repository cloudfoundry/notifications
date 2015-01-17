package services

import "github.com/cloudfoundry-incubator/notifications/models"

type Message struct {
	Status string
}

type MessagesRepoInterface interface {
	FindByID(models.ConnectionInterface, string) (models.Message, error)
}

type MessageFinder struct {
	repo     MessagesRepoInterface
	database models.DatabaseInterface
}

func NewMessageFinder(repo MessagesRepoInterface, database models.DatabaseInterface) MessageFinder {
	return MessageFinder{
		repo:     repo,
		database: database,
	}
}

func (finder MessageFinder) Find(messageID string) (Message, error) {
	message, err := finder.repo.FindByID(finder.database.Connection(), messageID)
	if err != nil {
		return Message{}, err
	}

	return Message{Status: message.Status}, nil
}
