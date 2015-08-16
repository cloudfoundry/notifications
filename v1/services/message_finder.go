package services

import "github.com/cloudfoundry-incubator/notifications/v1/models"

type Message struct {
	Status string
}

type messagesRepoFinder interface {
	FindByID(models.ConnectionInterface, string) (models.Message, error)
}

type MessageFinder struct {
	repo messagesRepoFinder
}

func NewMessageFinder(repo messagesRepoFinder) MessageFinder {
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
