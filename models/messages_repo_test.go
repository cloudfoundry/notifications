package models_test

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MessagesRepo", func() {
	var repo models.MessagesRepo
	var conn models.ConnectionInterface
	var message models.Message

	BeforeEach(func() {
		TruncateTables()
		repo = models.NewMessagesRepo()
		env := application.NewEnvironment()
		conn = models.NewDatabase(models.Config{
			DatabaseURL:    env.DatabaseURL,
			MigrationsPath: env.ModelMigrationsDir,
		}).Connection()
		message = models.Message{
			ID:     "message-id-123",
			Status: postal.StatusDelivered,
		}

	})

	Describe("FindByID", func() {
		It("finds messages created in the database", func() {
			message, err := repo.Create(conn, message)
			if err != nil {
				panic(err)
			}

			messageFound, err := repo.FindByID(conn, message.ID)
			Expect(err).ToNot(HaveOccurred())

			Expect(messageFound).To(Equal(message))
		})

		Context("When the message does not exists", func() {
			It("FindByID returns a models.RecordNotFoundError", func() {
				_, err := repo.FindByID(conn, "missing-id")
				Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))
			})
		})
	})

	Describe("Upsert", func() {
		Context("when no record exists yet with the message id", func() {
			It("inserts a new record", func() {
				message.UpdatedAt = time.Now().Add(100 * time.Hour)
				_, err := repo.Upsert(conn, message)
				if err != nil {
					panic(err)
				}

				messageFound, err := repo.FindByID(conn, message.ID)
				Expect(err).ToNot(HaveOccurred())

				Expect(messageFound.ID).To(Equal(message.ID))
				Expect(messageFound.Status).To(Equal(message.Status))
				Expect(messageFound.UpdatedAt).To(BeTemporally("~", time.Now(), time.Minute))
			})
		})
		Context("when a record already exists with the message id", func() {
			It("updates the existing record", func() {
				_, err := repo.Create(conn, message)
				if err != nil {
					panic(err)
				}

				updatedMessage := message
				updatedMessage.Status = postal.StatusFailed

				updatedMessage.UpdatedAt = time.Now().Add(100 * time.Hour)
				updatedMessage, err = repo.Upsert(conn, updatedMessage)
				if err != nil {
					panic(err)
				}

				messageFound, err := repo.FindByID(conn, message.ID)
				Expect(err).ToNot(HaveOccurred())

				Expect(messageFound.UpdatedAt).To(BeTemporally("~", time.Now(), time.Minute))
				Expect(messageFound.ID).To(Equal(updatedMessage.ID))
				Expect(messageFound.Status).To(Equal(updatedMessage.Status))
			})
		})
	})
})
