package gobble_test

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/gobble"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var sqlDB *sql.DB

func TestGobbleSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gobble")
}

var _ = BeforeSuite(func() {
	var err error
	sqlDB, err = instantiateDBConnection()
	Expect(err).NotTo(HaveOccurred())

	env, err := application.NewEnvironment()
	Expect(err).NotTo(HaveOccurred())

	gobble.NewDatabase(sqlDB).Migrate(env.GobbleMigrationsPath)
})

var _ = BeforeEach(func() {
	var err error
	sqlDB, err = instantiateDBConnection()
	Expect(err).NotTo(HaveOccurred())
})

func instantiateDBConnection() (*sql.DB, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	databaseURL = strings.TrimPrefix(databaseURL, "http://")
	databaseURL = strings.TrimPrefix(databaseURL, "https://")
	databaseURL = strings.TrimPrefix(databaseURL, "tcp://")
	databaseURL = strings.TrimPrefix(databaseURL, "mysql://")
	databaseURL = strings.TrimPrefix(databaseURL, "mysql2://")
	parsedURL, err := url.Parse("tcp://" + databaseURL)
	if err != nil {
		Fail(fmt.Sprintf("Could not parse DATABASE_URL %q, it does not fit format %q", os.Getenv("DATABASE_URL"), "tcp://user:pass@host/dname"))
	}

	password, _ := parsedURL.User.Password()
	databaseURL = fmt.Sprintf("%s:%s@%s(%s)%s?parseTime=true", parsedURL.User.Username(), password, parsedURL.Scheme, parsedURL.Host, parsedURL.Path)

	return sql.Open("mysql", databaseURL)
}

func TruncateTables() {
	database := gobble.NewDatabase(sqlDB)
	database.Connection.TruncateTables()
}
