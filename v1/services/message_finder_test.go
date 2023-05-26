package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/postal/common"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MessageFinder.Find", func() {
	var (
		finder       services.MessageFinder
		messagesRepo *mocks.MessagesRepo
		database     *mocks.Database
		conn         *mocks.Connection
	)

	BeforeEach(func() {
		messagesRepo = mocks.NewMessagesRepo()
		conn = mocks.NewConnection()
		database = mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = conn

		finder = services.NewMessageFinder(messagesRepo)
	})

	Context("when a message exists with the given id", func() {
		It("returns the right Message struct", func() {
			messagesRepo.FindByIDCall.Returns.Message = models.Message{Status: common.StatusDelivered}

			message, err := finder.Find(database, "a-message-id")

			Expect(err).NotTo(HaveOccurred())
			Expect(message.Status).To(Equal(common.StatusDelivered))

			Expect(messagesRepo.FindByIDCall.Receives.Connection).To(Equal(conn))
			Expect(messagesRepo.FindByIDCall.Receives.MessageID).To(Equal("a-message-id"))
		})
	})

	Context("when the underlying repo returns an error", func() {
		It("bubbles up the error", func() {
			messagesRepo.FindByIDCall.Returns.Error = errors.New("some error")

			_, err := finder.Find(database, "a-message-id")
			Expect(err).To(MatchError(errors.New("some error")))
		})
	})
})
