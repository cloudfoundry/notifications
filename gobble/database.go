package gobble

import (
	"database/sql"

	"bitbucket.org/liamstask/goose/lib/goose"
	"github.com/go-gorp/gorp"

	_ "github.com/go-sql-driver/mysql"
)

type DatabaseInterface interface {
	Migrate(string)
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

	conn.AddTableWithName(Job{}, "jobs").SetKeys(true, "ID").SetVersionCol("Version")

	return &DB{Connection: conn}
}

func (db DB) Migrate(migrationsDir string) {
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
