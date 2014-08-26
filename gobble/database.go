package gobble

import (
    "database/sql"
    "sync"

    "bitbucket.org/liamstask/goose/lib/goose"
    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/coopernurse/gorp"

    _ "github.com/go-sql-driver/mysql"
)

var (
    _database *DB
    mutex     sync.Mutex
)

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

    conn := &gorp.DbMap{
        Db: db,
        Dialect: gorp.MySQLDialect{
            Engine:   "InnoDB",
            Encoding: "UTF8",
        },
    }

    migrate(env)
    conn.AddTableWithName(Job{}, "jobs").SetKeys(true, "ID")

    _database = &DB{
        Connection: conn,
    }

    return _database
}

func migrate(env config.Environment) {
    dbDriver := goose.DBDriver{
        Name:    "mysql",
        OpenStr: env.DatabaseURL,
        Import:  "github.com/go-sql-driver/mysql",
        Dialect: &goose.MySqlDialect{},
    }
    migrationsDir := env.RootPath + "/gobble/migrations"

    dbConf := &goose.DBConf{
        MigrationsDir: migrationsDir,
        Env:           "notifications",
        Driver:        dbDriver,
    }

    current, err := goose.GetDBVersion(dbConf)
    if err != nil {
        panic(err)
    }

    target, err := goose.GetMostRecentDBVersion(dbConf.MigrationsDir)
    if err != nil {
        panic(err)
    }

    if current != target {
        err = goose.RunMigrations(dbConf, dbConf.MigrationsDir, target)
        if err != nil {
            panic(err)
        }
    }
}
