package v2_test

import (
	"bytes"
	"errors"

	"github.com/cloudfoundry-incubator/notifications/postal/v2"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/models"
	"github.com/pivotal-golang/lager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("V2MessageStatusUpdater", func() {
	var (
		updater      v2.V2MessageStatusUpdater
		messagesRepo *mocks.MessagesRepository
		logger       lager.Logger
		buffer       *bytes.Buffer
		conn         *mocks.Connection
	)

	BeforeEach(func() {
		conn = mocks.NewConnection()
		messagesRepo = mocks.NewMessagesRepository()

		buffer = bytes.NewBuffer([]byte{})
		logger = lager.NewLogger("notifications")
		logger.RegisterSink(lager.NewWriterSink(buffer, lager.INFO))

		updater = v2.NewV2MessageStatusUpdater(messagesRepo)
	})

	It("updates the status of the message", func() {
		updater.Update(conn, "some-message-id", "message-status", "campaign-id", logger)

		Expect(messagesRepo.UpdateCall.Receives.Connection).To(Equal(conn))
		Expect(messagesRepo.UpdateCall.Receives.Message).To(Equal(models.Message{
			ID:         "some-message-id",
			Status:     "message-status",
			CampaignID: "campaign-id",
		}))
	})

	Context("failure cases", func() {
		It("logs the error when the repository fails to update", func() {
			messagesRepo.UpdateCall.Returns.Error = errors.New("failed to update")

			updater.Update(conn, "some-message-id", "message-status", "campaign-id", logger)

			lines, err := parseLogLines(buffer.Bytes())
			Expect(err).NotTo(HaveOccurred())

			Expect(lines).To(HaveLen(1))
			line := lines[0]

			Expect(line).To(Equal(logLine{
				Source:   "notifications",
				Message:  "notifications.message-updater.failed-message-status-update",
				LogLevel: int(lager.ERROR),
				Data: map[string]interface{}{
					"session": "1",
					"error":   "failed to update",
					"status":  "message-status",
				},
			}))
		})
	})
})
