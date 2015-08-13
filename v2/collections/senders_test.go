package collections_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SendersCollection", func() {
	var (
		sendersCollection collections.SendersCollection
		sendersRepository *fakes.SendersRepository
		conn              *fakes.Connection
	)

	BeforeEach(func() {
		sendersRepository = fakes.NewSendersRepository()

		sendersCollection = collections.NewSendersCollection(sendersRepository)
		conn = fakes.NewConnection()
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
				sendersRepository.GetCall.Returns.Err = models.RecordNotFoundError("sender not found")

				_, err := sendersCollection.Get(conn, "some-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.NotFoundError{models.RecordNotFoundError("sender not found")}))
			})

			It("generates a not found error when the sender belongs to a different client", func() {
				_, err := sendersCollection.Get(conn, "some-sender-id", "mismatch-client-id")
				Expect(err).To(MatchError(collections.NotFoundError{errors.New("sender not found")}))
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
