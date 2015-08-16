package fakes

import (
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/v1/models"
)

type UnsubscribesRepo struct {
	store map[string]bool
}

func NewUnsubscribesRepo() *UnsubscribesRepo {
	return &UnsubscribesRepo{
		store: make(map[string]bool),
	}
}

func (fake *UnsubscribesRepo) Set(conn models.ConnectionInterface, userID, clientID, kindID string, unsubscribe bool) error {
	key := fmt.Sprintf("%s|%s|%s", userID, clientID, kindID)
	fake.store[key] = unsubscribe
	return nil
}

func (fake *UnsubscribesRepo) Get(conn models.ConnectionInterface, userID, clientID, kindID string) (bool, error) {
	key := fmt.Sprintf("%s|%s|%s", userID, clientID, kindID)
	return fake.store[key], nil
}
