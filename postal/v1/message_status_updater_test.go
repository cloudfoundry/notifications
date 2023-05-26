package v1_test

import (
	"bytes"
	"errors"

	"github.com/cloudfoundry-incubator/notifications/postal/v1"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/pivotal-golang/lager"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MessageStatusUpdater", func() {
	var (
		updater      v1.MessageStatusUpdater
		messagesRepo *mocks.MessagesRepo
		logger       lager.Logger
		buffer       *bytes.Buffer
		conn         *mocks.Connection
	)

	BeforeEach(func() {
		conn = mocks.NewConnection()
		messagesRepo = mocks.NewMessagesRepo()
		messagesRepo.UpsertCall.Returns.Messages = []models.Message{
			{
				ID:     "some-message-id",
				Status: "message-status",
			},
		}

		buffer = bytes.NewBuffer([]byte{})
		logger = lager.NewLogger("notifications")
		logger.RegisterSink(lager.NewWriterSink(buffer, lager.INFO))

		updater = v1.NewMessageStatusUpdater(messagesRepo)
	})

	It("updates the status of the message", func() {
		updater.Update(conn, "some-message-id", "message-status", "campaign-id", logger)

		Expect(messagesRepo.UpsertCall.Receives.Connection).To(Equal(conn))
		Expect(messagesRepo.UpsertCall.Receives.Messages[0]).To(Equal(models.Message{
			ID:     "some-message-id",
			Status: "message-status",
		}))
	})

	Context("failure cases", func() {
		It("logs the error when the repository fails to upsert", func() {
			messagesRepo.UpsertCall.Returns.Error = errors.New("failed to upsert")

			updater.Update(conn, "some-message-id", "message-status", "campaign-id", logger)

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
