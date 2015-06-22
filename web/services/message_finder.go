package services

import "github.com/cloudfoundry-incubator/notifications/models"

type Message struct {
	Status string
}

type MessagesRepoInterface interface {
	FindByID(models.ConnectionInterface, string) (models.Message, error)
}

type MessageFinderInterface interface {
	Find(models.DatabaseInterface, string) (Message, error)
}

type MessageFinder struct {
	repo MessagesRepoInterface
}

func NewMessageFinder(repo MessagesRepoInterface) MessageFinder {
	return MessageFinder{
		repo: repo,
	}
}

func (finder MessageFinder) Find(database models.DatabaseInterface, messageID string) (Message, error) {
	message, err := finder.repo.FindByID(database.Connection(), messageID)
	if err != nil {
		return Message{}, err
	}

	return Message{Status: message.Status}, nil
}
