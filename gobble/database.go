package gobble

import (
	"database/sql"
	"os"

	"bitbucket.org/liamstask/goose/lib/goose"
	"github.com/coopernurse/gorp"

	_ "github.com/go-sql-driver/mysql"
)

type DatabaseInterface interface {
	Migrate()
}

type DB struct {
	Connection *gorp.DbMap
}

func NewDatabase(db *sql.DB) *DB {
	conn := &gorp.DbMap{
		Db: db,
		Dialect: gorp.MySQLDialect{
			Engine:   "InnoDB",
			Encoding: "UTF8",
		},
	}

	conn.AddTableWithName(Job{}, "jobs").SetKeys(true, "ID")

	return &DB{Connection: conn}
}

func (db DB) Migrate() {
	migrationsDir := os.Getenv("GOBBLE_MIGRATIONS_DIR")
	dbConf := &goose.DBConf{
		MigrationsDir: migrationsDir,
		Env:           "notifications",
		Driver: goose.DBDriver{
			Dialect: &goose.MySqlDialect{},
		},
	}

	target, err := goose.GetMostRecentDBVersion(dbConf.MigrationsDir)
	if err != nil {
		panic(err)
	}

	err = goose.RunMigrationsOnDb(dbConf, dbConf.MigrationsDir, target, db.Connection.Db)
	if err != nil {
		panic(err)
	}
}
