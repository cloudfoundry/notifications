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
    Connection *gorp.DbMap
}

type ConnectionInterface interface {
    Delete(...interface{}) (int64, error)
    Insert(...interface{}) error
    Select(interface{}, string, ...interface{}) ([]interface{}, error)
    SelectOne(interface{}, string, ...interface{}) error
    Update(...interface{}) (int64, error)
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

    _database = &DB{
        Connection: &gorp.DbMap{
            Db: db,
            Dialect: gorp.MySQLDialect{
                Engine:   "InnoDB",
                Encoding: "UTF8",
            },
        },
    }

    _database.migrate()

    return _database
}

func (database DB) migrate() {
    database.Connection.AddTableWithName(Client{}, "clients").SetKeys(true, "Primary").ColMap("ID").SetUnique(true)
    database.Connection.AddTableWithName(Kind{}, "kinds").SetKeys(true, "Primary").SetUniqueTogether("id", "client_id")
    database.Connection.AddTableWithName(Receipt{}, "receipts").SetKeys(true, "Primary").SetUniqueTogether("user_guid", "client_id", "kind_id")

    err := database.Connection.CreateTablesIfNotExists()
    if err != nil {
        panic(err)
    }

}
