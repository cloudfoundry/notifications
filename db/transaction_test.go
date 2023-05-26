package db_test

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/testing/helpers"
	"github.com/cloudfoundry-incubator/notifications/v1/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Transaction", func() {
	var (
		transaction db.TransactionInterface
		conn        db.ConnectionInterface
	)

	BeforeEach(func() {
		db := db.NewDatabase(sqlDB, db.Config{})
		helpers.TruncateTables(db)
		conn = db.Connection()
		transaction = conn.Transaction()
	})

	Describe("Begin/Commit", func() {
		It("commits the transaction to the database", func() {
			err := transaction.Begin()
			Expect(err).NotTo(HaveOccurred())

			repo := models.NewClientsRepo()
			_, err = repo.Upsert(transaction, models.Client{
				ID:          "my-client",
				Description: "My Client",
			})
			Expect(err).NotTo(HaveOccurred())

			err = transaction.Commit()
			Expect(err).NotTo(HaveOccurred())

			client, err := repo.Find(conn, "my-client")
			Expect(err).NotTo(HaveOccurred())

			Expect(client.ID).To(Equal("my-client"))
			Expect(client.Description).To(Equal("My Client"))
		})
	})

	Describe("Begin/Rollback", func() {
		It("rolls back the transaction from the database", func() {
			err := transaction.Begin()
			Expect(err).NotTo(HaveOccurred())

			repo := models.NewClientsRepo()
			_, err = repo.Upsert(transaction, models.Client{
				ID:          "my-client",
				Description: "My Client",
			})
			Expect(err).NotTo(HaveOccurred())

			err = transaction.Rollback()
			Expect(err).NotTo(HaveOccurred())

			_, err = repo.Find(conn, "my-client")
			Expect(err).To(BeAssignableToTypeOf(models.NotFoundError{}))
		})
	})
})
