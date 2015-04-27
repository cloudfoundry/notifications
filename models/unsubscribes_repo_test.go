package models_test

import (
	"github.com/cloudfoundry-incubator/notifications/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UnsubscribesRepo", func() {
	var repo models.UnsubscribesRepo
	var conn *models.Connection

	BeforeEach(func() {
		TruncateTables()
		repo = models.NewUnsubscribesRepo()

		db := models.NewDatabase(sqlDB, models.Config{})
		db.Setup()
		conn = db.Connection().(*models.Connection)
	})

	Describe("Get/Set", func() {
		It("returns false for unsubscribes that have not been set", func() {
			isUnsubscribed, err := repo.Get(conn, "user-id", "client-id", "kind-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(isUnsubscribed).To(BeFalse())
		})

		It("returns true for unsubscribes that have been set", func() {
			err := repo.Set(conn, "user-id", "client-id", "kind-id", true)
			Expect(err).NotTo(HaveOccurred())

			isUnsubscribed, err := repo.Get(conn, "user-id", "client-id", "kind-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(isUnsubscribed).To(BeTrue())
		})

		It("returns false for unsubscribes that have been explicitly unsubscribed", func() {
			err := repo.Set(conn, "user-id", "client-id", "kind-id", false)
			Expect(err).NotTo(HaveOccurred())

			isUnsubscribed, err := repo.Get(conn, "user-id", "client-id", "kind-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(isUnsubscribed).To(BeFalse())
		})

		It("returns false for unsubscribes that have been unset", func() {
			err := repo.Set(conn, "user-id", "client-id", "kind-id", true)
			Expect(err).NotTo(HaveOccurred())

			isUnsubscribed, err := repo.Get(conn, "user-id", "client-id", "kind-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(isUnsubscribed).To(BeTrue())

			err = repo.Set(conn, "user-id", "client-id", "kind-id", false)
			Expect(err).NotTo(HaveOccurred())

			isUnsubscribed, err = repo.Get(conn, "user-id", "client-id", "kind-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(isUnsubscribed).To(BeFalse())
		})
	})

	Describe("FindAllByUserID", func() {
		It("finds all unsubscribes for a user", func() {
			err := repo.Set(conn, "correct-user", "raptors", "hungry", true)
			Expect(err).NotTo(HaveOccurred())

			err = repo.Set(conn, "correct-user", "raptors", "sleepy", true)
			Expect(err).NotTo(HaveOccurred())

			err = repo.Set(conn, "other-user", "dogs", "barking", true)
			Expect(err).NotTo(HaveOccurred())

			unsubscribes, err := repo.FindAllByUserID(conn, "correct-user")
			if err != nil {
				panic(err)
			}

			Expect(unsubscribes).To(HaveLen(2))
		})
	})
})
