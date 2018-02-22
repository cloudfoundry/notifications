package application

import (
	"crypto/rand"
	"database/sql"
	"log"
	"os"
	"path"
	"time"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/util"
	v1models "github.com/cloudfoundry-incubator/notifications/v1/models"
)

type DBProvider struct {
	sqlDB *sql.DB
	env   Environment
}

func NewDBProvider(env Environment) *DBProvider {
	var err error
	sqlDB, err := sql.Open("mysql", env.DatabaseURL)
	if err != nil {
		panic(err)
	}

	if err := sqlDB.Ping(); err != nil {
		panic(err)
	}

	sqlDB.SetMaxOpenConns(env.DBMaxOpenConns)
	return &DBProvider{
		sqlDB: sqlDB,
		env:   env,
	}
}

func (d *DBProvider) GobbleDatabase() gobble.DatabaseInterface {
	return gobble.NewDatabase(d.sqlDB)
}

func (d *DBProvider) Queue() gobble.QueueInterface {
	return gobble.NewQueue(d.GobbleDatabase(), util.NewClock(), gobble.Config{
		WaitMaxDuration: time.Duration(d.env.GobbleWaitMaxDuration) * time.Millisecond,
	})
}

func (d *DBProvider) Database() db.DatabaseInterface {
	database := v1models.NewDatabase(d.sqlDB, v1models.Config{
		DefaultTemplatePath: path.Join(d.env.RootPath, "templates", "default.json"),
	})

	if d.env.DBLoggingEnabled {
		database.TraceOn("[DB]", log.New(os.Stdout, "", 0))
	}

	return database
}

func (d *DBProvider) MessagesRepo() v1models.MessagesRepo {
	return v1models.NewMessagesRepo(util.NewIDGenerator(rand.Reader).Generate)
}
