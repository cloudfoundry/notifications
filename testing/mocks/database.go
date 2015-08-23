package mocks

import (
	"database/sql"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/go-gorp/gorp"
)

type Database struct {
	Conn                *Connection
	ConnectionWasCalled bool
	TracePrefix         string
	TraceLogger         gorp.GorpLogger
	rawDB               *sql.DB
}

func NewDatabase() *Database {
	return &Database{
		Conn:  NewConnection(),
		rawDB: &sql.DB{},
	}
}

func (fake *Database) Connection() db.ConnectionInterface {
	fake.ConnectionWasCalled = true
	return fake.Conn
}

func (fake *Database) TraceOn(prefix string, logger gorp.GorpLogger) {
	fake.TracePrefix = prefix
	fake.TraceLogger = logger
}

func (fake *Database) RawConnection() *sql.DB {
	return fake.rawDB
}

type SQLDatabase struct {
	MaxOpenConnections int
}

func (fake *SQLDatabase) SetMaxOpenConns(n int) {
	fake.MaxOpenConnections = n
}
