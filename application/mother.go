package application

import (
	"crypto/rand"
	"database/sql"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/util"
	v1models "github.com/cloudfoundry-incubator/notifications/v1/models"
)

type Mother struct {
	sqlDB *sql.DB
	mutex sync.Mutex
	env   Environment
}

func NewMother(env Environment) *Mother {
	return &Mother{
		env: env,
	}
}

func (m *Mother) GobbleDatabase() gobble.DatabaseInterface {
	return gobble.NewDatabase(m.SQLDatabase())
}

func (m *Mother) Queue() gobble.QueueInterface {
	return gobble.NewQueue(m.GobbleDatabase(), util.NewClock(), gobble.Config{
		WaitMaxDuration: time.Duration(m.env.GobbleWaitMaxDuration) * time.Millisecond,
	})
}

func (m *Mother) Database() db.DatabaseInterface {
	database := v1models.NewDatabase(m.SQLDatabase(), v1models.Config{
		DefaultTemplatePath: path.Join(m.env.RootPath, "templates", "default.json"),
	})

	if m.env.DBLoggingEnabled {
		database.TraceOn("[DB]", log.New(os.Stdout, "", 0))
	}

	return database
}

func (m *Mother) SQLDatabase() *sql.DB {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.sqlDB != nil {
		return m.sqlDB
	}

	var err error
	m.sqlDB, err = sql.Open("mysql", m.env.DatabaseURL)
	if err != nil {
		panic(err)
	}

	if err := m.sqlDB.Ping(); err != nil {
		panic(err)
	}

	m.sqlDB.SetMaxOpenConns(m.env.DBMaxOpenConns)

	return m.sqlDB
}

func (m *Mother) MessagesRepo() v1models.MessagesRepo {
	return v1models.NewMessagesRepo(util.NewIDGenerator(rand.Reader).Generate)
}
