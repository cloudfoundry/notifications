package fakes

import (
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/coopernurse/gorp"
)

type Database struct {
    Conn *FakeDBConn
}

func NewDatabase() *Database {
    return &Database{
        Conn: &FakeDBConn{},
    }
}

func (fake Database) Connection() models.ConnectionInterface {
    return fake.Conn
}

func (fake Database) TraceOn(prefix string, logger gorp.GorpLogger) {}
