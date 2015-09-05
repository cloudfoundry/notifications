package postal_test

import (
	"bytes"
	"errors"
	"log"
	"time"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MessageGC", func() {
	var messageGC postal.MessageGC
	var repo *mocks.MessagesRepo
	var oldMessageID string
	var newMessageID string
	var database *mocks.Database
	var conn db.ConnectionInterface
	var loggerBuffer *bytes.Buffer
	var lifetime time.Duration
	var pollingInterval time.Duration

	BeforeEach(func() {
		loggerBuffer = bytes.NewBuffer([]byte{})
		logger := log.New(loggerBuffer, "", 0)

		conn = mocks.NewConnection()
		database = mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = conn

		repo = mocks.NewMessagesRepo()
		lifetime = 2 * time.Minute
		pollingInterval = 500 * time.Millisecond
		oldMessageID = "that-message"
		newMessageID = "this-message"

		messageGC = postal.NewMessageGC(lifetime, database, repo, pollingInterval, logger)
	})

	Describe("Run", func() {
		It("It calls collect every passed in duration", func() {
			messageGC.Run()

			Eventually(func() int {
				return repo.DeleteBeforeCall.CallCount
			}).Should(BeNumerically(">=", 2))

			call1 := repo.DeleteBeforeCall.InvocationTimes[0]
			call2 := repo.DeleteBeforeCall.InvocationTimes[1]
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
			Expect(err).NotTo(HaveOccurred())

			_, err = repo.Upsert(conn, models.Message{
				ID:        newMessageID,
				UpdatedAt: time.Now(),
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("Deletes message statuses older than the specified time", func() {
			messageGC.Collect()

			Expect(repo.DeleteBeforeCall.Receives.Connection).To(Equal(conn))
			Expect(repo.DeleteBeforeCall.Receives.ThresholdTime).To(BeTemporally("~", time.Now().Add(-2*time.Minute), 10*time.Second))
		})

		Context("When the repo errors unexpectantly", func() {
			It("logs the error", func() {
				repo.DeleteBeforeCall.Returns.Error = errors.New("messages table is totally corrupt")

				messageGC.Collect()

				Expect(loggerBuffer.String()).To(ContainSubstring("messages table is totally corrupt"))
			})
		})

	})
})
