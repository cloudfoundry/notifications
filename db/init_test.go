package db_test

import (
	"database/sql"
	"testing"

	"github.com/cloudfoundry-incubator/notifications/application"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDBSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "db")
}

var sqlDB *sql.DB

var _ = BeforeEach(func() {
	env, err := application.NewEnvironment()
	Expect(err).NotTo(HaveOccurred())

	sqlDB, err = sql.Open("mysql", env.DatabaseURL)
	Expect(err).NotTo(HaveOccurred())
})
