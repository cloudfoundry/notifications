package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/go-gorp/gorp"
)

type Database struct {
	Conn                *Connection
	ConnectionWasCalled bool
	SeedWasCalled       bool
	MigrateWasCalled    bool
	MigrationsPath      string
	TracePrefix         string
	TraceLogger         gorp.GorpLogger
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

func (fake *Database) TraceOn(prefix string, logger gorp.GorpLogger) {
	fake.TracePrefix = prefix
	fake.TraceLogger = logger
}

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
