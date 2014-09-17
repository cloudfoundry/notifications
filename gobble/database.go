package gobble

import (
    "database/sql"
    "errors"
    "fmt"
    "net/url"
    "os"
    "strings"
    "sync"

    "bitbucket.org/liamstask/goose/lib/goose"
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

    databaseURL := loadDatabaseURL()
    db, err := sql.Open("mysql", databaseURL)
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

    migrationsDir := os.Getenv("GOBBLE_MIGRATIONS_DIR")
    migrate(databaseURL, migrationsDir)
    conn.AddTableWithName(Job{}, "jobs").SetKeys(true, "ID")

    _database = &DB{
        Connection: conn,
    }

    return _database
}

func loadDatabaseURL() string {
    databaseURL := os.Getenv("DATABASE_URL")
    databaseURL = strings.TrimPrefix(databaseURL, "http://")
    databaseURL = strings.TrimPrefix(databaseURL, "https://")
    databaseURL = strings.TrimPrefix(databaseURL, "tcp://")
    databaseURL = strings.TrimPrefix(databaseURL, "mysql://")
    databaseURL = strings.TrimPrefix(databaseURL, "mysql2://")
    parsedURL, err := url.Parse("tcp://" + databaseURL)
    if err != nil {
        panic(errors.New(fmt.Sprintf("Could not parse DATABASE_URL %q, it does not fit format %q", os.Getenv("DATABASE_URL"), "tcp://user:pass@host/dname")))
    }

    password, _ := parsedURL.User.Password()
    return fmt.Sprintf("%s:%s@%s(%s)%s?parseTime=true", parsedURL.User.Username(), password, parsedURL.Scheme, parsedURL.Host, parsedURL.Path)
}

func migrate(databaseURL, migrationsDir string) {
    dbDriver := goose.DBDriver{
        Name:    "mysql",
        OpenStr: databaseURL,
        Import:  "github.com/go-sql-driver/mysql",
        Dialect: &goose.MySqlDialect{},
    }

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
