package fakes

import "github.com/cloudfoundry-incubator/notifications/web/services"

type MessageFinder struct {
	Messages  map[string]services.Message
	FindError error
}

func NewMessageFinder() *MessageFinder {
	return &MessageFinder{
		Messages: map[string]services.Message{},
	}
}

func (finder MessageFinder) Find(messageID string) (services.Message, error) {
	return finder.Messages[messageID], finder.FindError
}
