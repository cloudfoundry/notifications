package application

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/util"
	v1models "github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/go-sql-driver/mysql"
)

type DBProvider struct {
	sqlDB *sql.DB
	env   Environment
}

func NewDBProvider(env Environment) *DBProvider {
	var err error
	databaseURL := env.DatabaseURL
	if env.DatabaseCACertFile != "" {
		registerTLSConfig(env)
		databaseURL += "&tls=custom"
		if !env.DatabaseEnableIdentityVerification {
			databaseURL += "&trustServerCertificate=true&verifyServerCertificate=true&disableSslHostnameVerification=true"
		}
	}

	sqlDB, err := sql.Open("mysql", databaseURL)
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

func registerTLSConfig(env Environment) {
	ca, err := ioutil.ReadFile(env.DatabaseCACertFile)
	if err != nil {
		panic(err)
	}
	rootCertPool := x509.NewCertPool()
	if ok := rootCertPool.AppendCertsFromPEM(ca); !ok {
		panic("Failed to append PEM when creating database connection")
	}
	mysql.RegisterTLSConfig("custom", &tls.Config{
		RootCAs:    rootCertPool,
		ServerName: env.DatabaseCommonName,
	})
}
