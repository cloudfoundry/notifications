package models_test

import (
	"database/sql"
	"testing"

	"github.com/cloudfoundry-incubator/notifications/application"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestModelsSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "v2/models")
}

var sqlDB *sql.DB

var _ = BeforeEach(func() {
	env, err := application.NewEnvironment()
	Expect(err).NotTo(HaveOccurred())

	sqlDB, err = sql.Open("mysql", env.DatabaseURL)
	Expect(err).NotTo(HaveOccurred())
})
