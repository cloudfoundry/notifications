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
    return _database
}
