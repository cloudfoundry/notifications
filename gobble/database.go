package gobble

import (
	"database/sql"

	sql_migrate "github.com/rubenv/sql-migrate"
	"gopkg.in/gorp.v1"

	_ "github.com/go-sql-driver/mysql"
)

type DatabaseInterface interface {
	Migrate(string)
}

type ConnectionInterface interface {
	Insert(...interface{}) error
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

	Initializer{}.InitializeDBMap(conn)

	return &DB{Connection: conn}
}

type Initializer struct{}

func (Initializer) InitializeDBMap(dbMap *gorp.DbMap) {
	dbMap.AddTableWithName(Job{}, "jobs").SetKeys(true, "ID").SetVersionCol("Version")
}

func (db DB) Migrate(migrationsPath string) {
	sql_migrate.SetTable("gobble_model_migrations")

	migrations := &sql_migrate.FileMigrationSource{
		Dir: migrationsPath,
	}

	_, err := sql_migrate.Exec(db.Connection.Db, "mysql", migrations, sql_migrate.Up)
	if err != nil {
		panic(err)
	}
}
