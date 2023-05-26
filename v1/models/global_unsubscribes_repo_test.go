package models_test

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/testing/helpers"
	"github.com/cloudfoundry-incubator/notifications/v1/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GlobalUnsubscribesRepo", func() {
	var repo models.GlobalUnsubscribesRepo
	var conn *db.Connection

	Describe("Set/Get", func() {
		BeforeEach(func() {
			database := db.NewDatabase(sqlDB, db.Config{})
			helpers.TruncateTables(database)
			conn = database.Connection().(*db.Connection)
			repo = models.NewGlobalUnsubscribesRepo()
		})

		It("sets the global unsubscribe field for a user, allowing it to be retrieved later", func() {
			err := repo.Set(conn, "my-user", true)
			if err != nil {
				panic(err)
			}

			unsubscribed, err := repo.Get(conn, "my-user")
			if err != nil {
				panic(err)
			}

			Expect(unsubscribed).To(BeTrue())

			err = repo.Set(conn, "my-user", false)
			if err != nil {
				panic(err)
			}

			unsubscribed, err = repo.Get(conn, "my-user")
			if err != nil {
				panic(err)
			}

			Expect(unsubscribed).To(BeFalse())
		})
	})
})
