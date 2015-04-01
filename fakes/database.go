package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/coopernurse/gorp"
)

type Database struct {
	Conn          *DBConn
	SeedWasCalled bool
}

func NewDatabase() *Database {
	return &Database{
		Conn: NewDBConn(),
	}
}

func (fake Database) Connection() models.ConnectionInterface {
	return fake.Conn
}

func (fake Database) TraceOn(prefix string, logger gorp.GorpLogger) {}

func (fake *Database) Seed() {
	fake.SeedWasCalled = true
}

type SQLDatabase struct {
	MaxOpenConnections int
}

func (fake *SQLDatabase) SetMaxOpenConns(n int) {
	fake.MaxOpenConnections = n
}
