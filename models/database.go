package models

import (
	"database/sql"

	"sync"

	"github.com/coopernurse/gorp"
	sql_migrate "github.com/rubenv/sql-migrate"

	_ "github.com/go-sql-driver/mysql"
)

var _database *DB
var mutex sync.Mutex

type DB struct {
	connection *Connection
}

type DatabaseInterface interface {
	Connection() ConnectionInterface
	TraceOn(string, gorp.GorpLogger)
}

func NewDatabase(databaseURL, migrationsPath string) *DB {
	if _database != nil {
		return _database
	}

	mutex.Lock()
	defer mutex.Unlock()
	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	connection := &Connection{
		DbMap: &gorp.DbMap{
			Db: db,
			Dialect: gorp.MySQLDialect{
				Engine:   "InnoDB",
				Encoding: "UTF8",
			},
		},
	}

	_database = &DB{
		connection: connection,
	}

	_database.migrate(migrationsPath)

	return _database
}

func (database DB) migrate(migrationsPath string) {
	sql_migrate.SetTable("notifications_model_migrations")

	migrations := &sql_migrate.FileMigrationSource{
		Dir: migrationsPath,
	}

	_, err := sql_migrate.Exec(database.connection.Db, "mysql", migrations, sql_migrate.Up)
	if err != nil {
		panic(err)
	}

	database.connection.AddTableWithName(Client{}, "clients").SetKeys(true, "Primary").ColMap("ID").SetUnique(true)
	database.connection.AddTableWithName(Kind{}, "kinds").SetKeys(true, "Primary").SetUniqueTogether("id", "client_id")
	database.connection.AddTableWithName(Receipt{}, "receipts").SetKeys(true, "Primary").SetUniqueTogether("user_guid", "client_id", "kind_id")
	database.connection.AddTableWithName(Unsubscribe{}, "unsubscribes").SetKeys(true, "Primary").SetUniqueTogether("user_id", "client_id", "kind_id")
	database.connection.AddTableWithName(GlobalUnsubscribe{}, "global_unsubscribes").SetKeys(true, "Primary").ColMap("UserID").SetUnique(true)
	database.connection.AddTableWithName(Template{}, "templates").SetKeys(true, "Primary").ColMap("Name").SetUnique(true)
}

func (database *DB) Connection() ConnectionInterface {
	return database.connection
}

func (database *DB) TraceOn(prefix string, logger gorp.GorpLogger) {
	database.connection.TraceOn(prefix, logger)
}
