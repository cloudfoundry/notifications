package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
)

type MessageFinder struct {
	Messages map[string]services.Message

	FindCall struct {
		Arguments []interface{}
		Error     error
	}
}

func NewMessageFinder() *MessageFinder {
	return &MessageFinder{
		Messages: map[string]services.Message{},
	}
}

func (finder *MessageFinder) Find(database models.DatabaseInterface, messageID string) (services.Message, error) {
	finder.FindCall.Arguments = []interface{}{database, messageID}
	return finder.Messages[messageID], finder.FindCall.Error
}
