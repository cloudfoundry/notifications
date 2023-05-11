package db_test

import (
	"github.com/cloudfoundry-incubator/notifications/db"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Database", func() {
	var database *db.DB

	BeforeEach(func() {
		database = db.NewDatabase(sqlDB, db.Config{})
	})

	Describe("Connection", func() {
		It("returns a Connection", func() {
			connection := database.Connection()
			Expect(connection).To(BeAssignableToTypeOf(&db.Connection{}))
		})
	})
})
