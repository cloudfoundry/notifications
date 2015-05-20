package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/coopernurse/gorp"
)

type Database struct {
	Conn                *Connection
	ConnectionWasCalled bool
	SeedWasCalled       bool
	MigrateWasCalled    bool
	MigrationsPath      string
}

func NewDatabase() *Database {
	return &Database{
		Conn: NewConnection(),
	}
}

func (fake *Database) Connection() models.ConnectionInterface {
	fake.ConnectionWasCalled = true
	return fake.Conn
}

func (fake Database) TraceOn(prefix string, logger gorp.GorpLogger) {}

func (fake *Database) Seed() {
	fake.SeedWasCalled = true
}

func (fake *Database) Migrate(migrationsPath string) {
	fake.MigrationsPath = migrationsPath
	fake.MigrateWasCalled = true
}

func (*Database) Setup() {
}

type SQLDatabase struct {
	MaxOpenConnections int
}

func (fake *SQLDatabase) SetMaxOpenConns(n int) {
	fake.MaxOpenConnections = n
}
