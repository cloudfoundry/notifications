package postal_test

import (
	"bytes"
	"errors"
	"log"
	"time"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/testing/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MessageGC", func() {
	var messageGC postal.MessageGC
	var repo *fakes.MessagesRepo
	var oldMessageID string
	var newMessageID string
	var database *fakes.Database
	var conn db.ConnectionInterface
	var loggerBuffer *bytes.Buffer
	var lifetime time.Duration
	var pollingInterval time.Duration

	BeforeEach(func() {
		loggerBuffer = bytes.NewBuffer([]byte{})
		logger := log.New(loggerBuffer, "", 0)
		database = fakes.NewDatabase()
		conn = database.Connection()
		repo = fakes.NewMessagesRepo()
		lifetime = 2 * time.Minute
		pollingInterval = 500 * time.Millisecond
		messageGC = postal.NewMessageGC(lifetime, database, repo, pollingInterval, logger)
		oldMessageID = "that-message"
		newMessageID = "this-message"
	})

	Describe("Run", func() {
		It("It calls collect every passed in duration", func() {
			messageGC.Run()

			Eventually(func() int {
				return len(repo.DeleteBeforeInvocations)
			}).Should(BeNumerically(">=", 2))

			call1 := repo.DeleteBeforeInvocations[0]
			call2 := repo.DeleteBeforeInvocations[1]
			Expect(call2).To(BeTemporally(">", call1.Add(pollingInterval-50*time.Millisecond)))
			Expect(call2).To(BeTemporally("<", call1.Add(pollingInterval+50*time.Millisecond)))
		})
	})

	Describe("Collect", func() {
		BeforeEach(func() {
			_, err := repo.Upsert(conn, models.Message{
				ID:        oldMessageID,
				UpdatedAt: time.Now().Add(-2 * lifetime),
			})
			if err != nil {
				panic(err)
			}

			_, err = repo.Upsert(conn, models.Message{
				ID:        newMessageID,
				UpdatedAt: time.Now(),
			})
			if err != nil {
				panic(err)
			}
		})

		It("Deletes message statuses older than the specified time", func() {
			_, err := repo.FindByID(conn, oldMessageID)
			if err != nil {
				panic(err)
			}

			messageGC.Collect()

			_, err = repo.FindByID(conn, oldMessageID)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))
		})

		It("Does not delete messages newer than the specified time", func() {
			_, err := repo.FindByID(conn, newMessageID)
			if err != nil {
				panic(err)
			}

			messageGC.Collect()

			_, err = repo.FindByID(conn, newMessageID)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("When the repo errors unexpectantly", func() {
			It("logs the error", func() {
				repo.DeleteBeforeError = errors.New("messages table is totally corrupt or something")

				messageGC.Collect()

				Expect(loggerBuffer.String()).To(ContainSubstring(repo.DeleteBeforeError.Error()))
			})
		})

	})
})
