package collections_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SendersCollection", func() {
	var (
		sendersCollection collections.SendersCollection
		sendersRepository *mocks.SendersRepository
		conn              *mocks.Connection
	)

	BeforeEach(func() {
		sendersRepository = mocks.NewSendersRepository()

		sendersCollection = collections.NewSendersCollection(sendersRepository)
		conn = mocks.NewConnection()
	})

	Describe("Set", func() {
		BeforeEach(func() {
			sendersRepository.InsertCall.Returns.Sender = models.Sender{
				ID:       "some-sender-id",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}
		})

		It("adds a sender to the collection", func() {
			sender, err := sendersCollection.Set(conn, collections.Sender{
				Name:     "some-sender",
				ClientID: "some-client-id",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(sender).To(Equal(collections.Sender{
				ID:       "some-sender-id",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}))

			Expect(sendersRepository.InsertCall.Receives.Conn).To(Equal(conn))
			Expect(sendersRepository.InsertCall.Receives.Sender).To(Equal(models.Sender{
				Name:     "some-sender",
				ClientID: "some-client-id",
			}))
		})

		It("will idempotently add duplicates", func() {
			sendersRepository.InsertCall.Returns.Sender = models.Sender{}
			sendersRepository.InsertCall.Returns.Err = models.DuplicateRecordError{}
			sendersRepository.GetByClientIDAndNameCall.Returns.Sender = models.Sender{
				ID:       "some-sender-id",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}

			sender, err := sendersCollection.Set(conn, collections.Sender{
				Name:     "some-sender",
				ClientID: "some-client-id",
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(sender).To(Equal(collections.Sender{
				ID:       "some-sender-id",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}))
			Expect(sendersRepository.GetByClientIDAndNameCall.Receives.Conn).To(Equal(conn))
			Expect(sendersRepository.GetByClientIDAndNameCall.Receives.ClientID).To(Equal("some-client-id"))
			Expect(sendersRepository.GetByClientIDAndNameCall.Receives.Name).To(Equal("some-sender"))
		})

		Context("failure cases", func() {
			It("handles unexpected database errors", func() {
				sendersRepository.InsertCall.Returns.Sender = models.Sender{}
				sendersRepository.InsertCall.Returns.Err = errors.New("BOOM!")

				_, err := sendersCollection.Set(conn, collections.Sender{
					Name:     "some-sender",
					ClientID: "some-client-id",
				})
				Expect(err).To(MatchError(collections.PersistenceError{
					Err: errors.New("BOOM!"),
				}))
			})

			It("returns a persistence error when the sender cannot be found by client id and name", func() {
				sendersRepository.InsertCall.Returns.Sender = models.Sender{}
				sendersRepository.InsertCall.Returns.Err = models.DuplicateRecordError{}

				sendersRepository.GetByClientIDAndNameCall.Returns.Sender = models.Sender{}
				sendersRepository.GetByClientIDAndNameCall.Returns.Err = errors.New("BOOM!")

				_, err := sendersCollection.Set(conn, collections.Sender{
					Name:     "some-sender",
					ClientID: "some-client-id",
				})
				Expect(err).To(MatchError(collections.PersistenceError{
					Err: errors.New("BOOM!"),
				}))
			})
		})
	})

	Describe("List", func() {
		It("will list all senders in the collection", func() {
			sendersRepository.ListCall.Returns.Senders = []models.Sender{
				{
					ID:       "some-sender-id",
					Name:     "some-sender",
					ClientID: "some-client-id",
				},
			}

			senders, err := sendersCollection.List(conn, "some-client-id")
			Expect(err).NotTo(HaveOccurred())

			Expect(len(senders)).To(Equal(1))
			Expect(senders[0].ID).To(Equal("some-sender-id"))

			Expect(sendersRepository.ListCall.Receives.Conn).To(Equal(conn))
			Expect(sendersRepository.ListCall.Receives.ClientID).To(Equal("some-client-id"))
		})

		Context("failure cases", func() {
			It("handles unexpected database errors", func() {
				sendersRepository.ListCall.Returns.Err = errors.New("BOOM!")

				_, err := sendersCollection.List(conn, "some-client-id")
				Expect(err).To(MatchError(collections.PersistenceError{
					Err: errors.New("BOOM!"),
				}))
			})
		})
	})

	Describe("Get", func() {
		BeforeEach(func() {
			sendersRepository.GetCall.Returns.Sender = models.Sender{
				ID:       "some-sender-id",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}
		})

		It("will retrieve a sender from the collection", func() {
			sender, err := sendersCollection.Get(conn, "some-sender-id", "some-client-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(sender).To(Equal(collections.Sender{
				ID:       "some-sender-id",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}))

			Expect(sendersRepository.GetCall.Receives.Conn).To(Equal(conn))
			Expect(sendersRepository.GetCall.Receives.SenderID).To(Equal("some-sender-id"))
		})

		Context("failure cases", func() {
			It("generates a not found error when the sender does not exist", func() {
				recordNotFoundError := models.NewRecordNotFoundError("sender not found")
				sendersRepository.GetCall.Returns.Err = recordNotFoundError

				_, err := sendersCollection.Get(conn, "some-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.NotFoundError{recordNotFoundError}))
			})

			It("generates a not found error when the sender belongs to a different client", func() {
				_, err := sendersCollection.Get(conn, "some-sender-id", "mismatch-client-id")
				Expect(err).To(MatchError(collections.NotFoundError{errors.New(`Sender with id "some-sender-id" could not be found`)}))
			})

			It("handles unexpected database errors", func() {
				sendersRepository.GetCall.Returns.Sender = models.Sender{}
				sendersRepository.GetCall.Returns.Err = errors.New("BOOM!")

				_, err := sendersCollection.Get(conn, "some-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.PersistenceError{
					Err: errors.New("BOOM!"),
				}))
			})
		})
	})
})
