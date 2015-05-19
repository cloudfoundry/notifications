package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MessageFinder.Find", func() {
	var (
		finder       services.MessageFinder
		messagesRepo *fakes.MessagesRepo
		messageID    string
		database     *fakes.Database
	)

	BeforeEach(func() {
		messagesRepo = fakes.NewMessagesRepo()
		messageID = "a-message-id"
		database = fakes.NewDatabase()

		finder = services.NewMessageFinder(messagesRepo)
	})

	Context("when a message exists with the given id", func() {
		It("returns the right Message struct", func() {
			messagesRepo.Messages[messageID] = models.Message{Status: postal.StatusDelivered}

			message, err := finder.Find(database, messageID)

			Expect(err).NotTo(HaveOccurred())
			Expect(message.Status).To(Equal(postal.StatusDelivered))

			Expect(database.ConnectionWasCalled).To(BeTrue())
		})
	})

	Context("when the underlying repo returns an error", func() {
		It("bubbles up the error", func() {
			messagesRepo.FindByIDError = errors.New("generic repo error (it could be anything!)")

			_, err := finder.Find(database, messageID)
			Expect(err).To(MatchError(messagesRepo.FindByIDError))
		})
	})

})
