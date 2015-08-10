package models_test

import (
	"github.com/cloudfoundry-incubator/notifications/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Database", func() {
	var db *models.DB

	BeforeEach(func() {
		TruncateTables()
		db = models.NewDatabase(sqlDB, models.Config{})
	})

	Describe("Connection", func() {
		It("returns a Connection", func() {
			connection := db.Connection()
			Expect(connection).To(BeAssignableToTypeOf(&models.Connection{}))
		})
	})
})
