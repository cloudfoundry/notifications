package fakes

import "github.com/cloudfoundry-incubator/notifications/db"

type GlobalUnsubscribesRepo struct {
	unsubscribes []string
	SetError     error
}

func NewGlobalUnsubscribesRepo() *GlobalUnsubscribesRepo {
	return &GlobalUnsubscribesRepo{
		unsubscribes: make([]string, 0),
	}
}

func (fake *GlobalUnsubscribesRepo) Set(conn db.ConnectionInterface, userID string, globalUnsubscribe bool) error {
	if fake.SetError != nil {
		return fake.SetError
	}

	if globalUnsubscribe {
		fake.unsubscribes = append(fake.unsubscribes, userID)
	} else {
		for index, id := range fake.unsubscribes {
			if id == userID {
				fake.unsubscribes = append(fake.unsubscribes[:index], fake.unsubscribes[index+1:]...)
			}
			return nil
		}
	}

	return nil
}

func (fake *GlobalUnsubscribesRepo) Get(conn db.ConnectionInterface, userID string) (bool, error) {
	for _, id := range fake.unsubscribes {
		if id == userID {
			return true, nil
		}
	}

	return false, nil
}
