package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type UnsubscribesRepo struct {
    Unsubscribes map[string]models.Unsubscribe
}

func NewUnsubscribesRepo() *UnsubscribesRepo {
    return &UnsubscribesRepo{
        Unsubscribes: map[string]models.Unsubscribe{},
    }
}

func (fake *UnsubscribesRepo) Create(conn models.ConnectionInterface, unsubscribe models.Unsubscribe) (models.Unsubscribe, error) {
    key := unsubscribe.ClientID + unsubscribe.KindID + unsubscribe.UserID
    if _, ok := fake.Unsubscribes[key]; ok {
        return unsubscribe, models.ErrDuplicateRecord{}
    }
    fake.Unsubscribes[key] = unsubscribe
    return unsubscribe, nil
}

func (fake *UnsubscribesRepo) Upsert(conn models.ConnectionInterface, unsubscribe models.Unsubscribe) (models.Unsubscribe, error) {
    key := unsubscribe.ClientID + unsubscribe.KindID + unsubscribe.UserID
    fake.Unsubscribes[key] = unsubscribe
    return unsubscribe, nil
}

func (fake *UnsubscribesRepo) Find(conn models.ConnectionInterface, clientID string, kindID string, userID string) (models.Unsubscribe, error) {
    key := clientID + kindID + userID
    if unsubscribe, ok := fake.Unsubscribes[key]; ok {
        return unsubscribe, nil
    }
    return models.Unsubscribe{}, models.ErrRecordNotFound{}
}

func (fake *UnsubscribesRepo) Destroy(conn models.ConnectionInterface, unsubscribe models.Unsubscribe) (int, error) {
    key := unsubscribe.ClientID + unsubscribe.KindID + unsubscribe.UserID
    delete(fake.Unsubscribes, key)
    return 0, nil
}
