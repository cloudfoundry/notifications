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

	"errors"

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
		MaxQueueLength:  d.env.GobbleMaxQueueLength,
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

	tlsConfig := &tls.Config{
		RootCAs:            rootCertPool,
		ServerName:         env.DatabaseCommonName,
		InsecureSkipVerify: false,
	}

	if !env.DatabaseEnableIdentityVerification {
		tlsConfig.InsecureSkipVerify = true
		tlsConfig.VerifyPeerCertificate = func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			return VerifyCertificatesIgnoreHostname(rawCerts, rootCertPool)
		}
	}

	mysql.RegisterTLSConfig("custom", tlsConfig)
}

func VerifyCertificatesIgnoreHostname(rawCerts [][]byte, caCertPool *x509.CertPool) error {
	certs := make([]*x509.Certificate, len(rawCerts))

	for i, asn1Data := range rawCerts {
		cert, err := x509.ParseCertificate(asn1Data)
		if err != nil {
			return errors.New("tls: failed to parse certificate from server: " + err.Error())
		}

		certs[i] = cert
	}

	opts := x509.VerifyOptions{
		Roots:         caCertPool,
		CurrentTime:   time.Now(),
		Intermediates: x509.NewCertPool(),
	}

	for i, cert := range certs {
		if i == 0 {
			continue
		}

		opts.Intermediates.AddCert(cert)
	}

	_, err := certs[0].Verify(opts)
	return err
}
