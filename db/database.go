package db

import (
	"database/sql"

	"gopkg.in/gorp.v1"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	connection *Connection
	config     Config
}

type DatabaseInterface interface {
	Connection() ConnectionInterface
	TraceOn(string, gorp.GorpLogger)
	RawConnection() *sql.DB
}

func NewDatabase(db *sql.DB, config Config) *DB {
	connection := &Connection{
		DbMap: &gorp.DbMap{
			Db: db,
			Dialect: gorp.MySQLDialect{
				Engine:   "InnoDB",
				Encoding: "UTF8",
			},
		},
	}

	database := &DB{
		config:     config,
		connection: connection,
	}

	return database
}

func (database *DB) Connection() ConnectionInterface {
	return database.connection
}

type TableMapInterface interface {
	AddTableWithName(i interface{}, name string) *gorp.TableMap
}

func (database *DB) TableMap() TableMapInterface {
	return database.connection
}

func (database *DB) TraceOn(prefix string, logger gorp.GorpLogger) {
	database.connection.TraceOn(prefix, logger)
}

func (database *DB) RawConnection() *sql.DB {
	return database.connection.Db
}
