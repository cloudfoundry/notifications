package postal_test

import (
	"bytes"
	"errors"

	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/pivotal-golang/lager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MessageStatusUpdater", func() {
	var (
		updater      postal.MessageStatusUpdater
		messagesRepo *fakes.MessagesRepository
		logger       lager.Logger
		buffer       *bytes.Buffer
		conn         *fakes.Connection
	)

	BeforeEach(func() {
		conn = fakes.NewConnection()
		messagesRepo = fakes.NewMessagesRepository()

		buffer = bytes.NewBuffer([]byte{})
		logger = lager.NewLogger("notifications")
		logger.RegisterSink(lager.NewWriterSink(buffer, lager.INFO))

		updater = postal.NewMessageStatusUpdater(messagesRepo)
	})

	It("updates the status of the message", func() {
		updater.Update(conn, "some-message-id", "message-status", logger)

		Expect(messagesRepo.UpsertCall.Receives.Connection).To(Equal(conn))
		Expect(messagesRepo.UpsertCall.Receives.Message).To(Equal(models.Message{
			ID:     "some-message-id",
			Status: "message-status",
		}))
	})

	Context("failure cases", func() {
		It("logs the error when the repository fails to upsert", func() {
			messagesRepo.UpsertCall.Returns.Error = errors.New("failed to upsert")

			updater.Update(conn, "some-message-id", "message-status", logger)

			lines, err := parseLogLines(buffer.Bytes())
			Expect(err).NotTo(HaveOccurred())

			Expect(lines).To(HaveLen(1))
			line := lines[0]

			Expect(line).To(Equal(logLine{
				Source:   "notifications",
				Message:  "notifications.message-updater.failed-message-status-upsert",
				LogLevel: int(lager.ERROR),
				Data: map[string]interface{}{
					"session": "1",
					"error":   "failed to upsert",
					"status":  "message-status",
				},
			}))
		})
	})
})
