package models

import (
    "database/sql"

    "sync"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/coopernurse/gorp"

    _ "github.com/go-sql-driver/mysql"
)

var _database *DB
var mutex sync.Mutex

type DB struct {
    connection *Connection
}

func Database() *DB {
    if _database != nil {
        return _database
    }

    mutex.Lock()
    defer mutex.Unlock()
    env := config.NewEnvironment()
    db, err := sql.Open("mysql", env.DatabaseURL)
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

    _database.migrate()

    return _database
}

func (database DB) migrate() {
    database.connection.AddTableWithName(Client{}, "clients").SetKeys(true, "Primary").ColMap("ID").SetUnique(true)
    database.connection.AddTableWithName(Kind{}, "kinds").SetKeys(true, "Primary").SetUniqueTogether("id", "client_id")
    database.connection.AddTableWithName(Receipt{}, "receipts").SetKeys(true, "Primary").SetUniqueTogether("user_guid", "client_id", "kind_id")
    database.connection.AddTableWithName(Unsubscribe{}, "unsubscribes").SetKeys(true, "Primary").SetUniqueTogether("user_id", "client_id", "kind_id")

    err := database.connection.CreateTablesIfNotExists()
    if err != nil {
        panic(err)
    }

}

func (database *DB) Connection() *Connection {
    return database.connection
}
