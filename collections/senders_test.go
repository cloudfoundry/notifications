package collections_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"

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

	Describe("Add", func() {
		BeforeEach(func() {
			sendersRepository.InsertCall.ReturnSender = models.Sender{
				ID:       "some-sender-id",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}
		})

		It("adds a sender to the collection", func() {
			sender, err := sendersCollection.Add(conn, collections.Sender{
				Name:     "some-sender",
				ClientID: "some-client-id",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(sender).To(Equal(collections.Sender{
				ID:       "some-sender-id",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}))

			Expect(sendersRepository.InsertCall.Conn).To(Equal(conn))
			Expect(sendersRepository.InsertCall.Sender).To(Equal(models.Sender{
				Name:     "some-sender",
				ClientID: "some-client-id",
			}))
		})

		It("will idempotently add duplicates", func() {
			sendersRepository.InsertCall.ReturnSender = models.Sender{}
			sendersRepository.InsertCall.Err = models.DuplicateRecordError{}
			sendersRepository.GetByClientIDAndNameCall.ReturnSender = models.Sender{
				ID:       "some-sender-id",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}

			sender, err := sendersCollection.Add(conn, collections.Sender{
				Name:     "some-sender",
				ClientID: "some-client-id",
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(sender).To(Equal(collections.Sender{
				ID:       "some-sender-id",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}))
			Expect(sendersRepository.GetByClientIDAndNameCall.Conn).To(Equal(conn))
			Expect(sendersRepository.GetByClientIDAndNameCall.ClientID).To(Equal("some-client-id"))
			Expect(sendersRepository.GetByClientIDAndNameCall.Name).To(Equal("some-sender"))
		})

		Context("failure cases", func() {
			It("handles unexpected database errors", func() {
				sendersRepository.InsertCall.ReturnSender = models.Sender{}
				sendersRepository.InsertCall.Err = errors.New("BOOM!")

				_, err := sendersCollection.Add(conn, collections.Sender{
					Name:     "some-sender",
					ClientID: "some-client-id",
				})
				Expect(err).To(MatchError(collections.PersistenceError{
					Err: errors.New("BOOM!"),
				}))
			})

			It("validates that a sender has a name", func() {
				_, err := sendersCollection.Add(conn, collections.Sender{
					ClientID: "some-client-id",
				})
				Expect(err).To(MatchError(collections.NewValidationError("missing sender name")))
			})

			It("validates that a sender has a client id", func() {
				_, err := sendersCollection.Add(conn, collections.Sender{
					Name: "some-sender",
				})
				Expect(err).To(MatchError(collections.NewValidationError("missing sender client_id")))
			})

			It("returns a persistence error when the sender cannot be found by client id and name", func() {
				sendersRepository.InsertCall.ReturnSender = models.Sender{}
				sendersRepository.InsertCall.Err = models.DuplicateRecordError{}

				sendersRepository.GetByClientIDAndNameCall.ReturnSender = models.Sender{}
				sendersRepository.GetByClientIDAndNameCall.Err = errors.New("BOOM!")

				_, err := sendersCollection.Add(conn, collections.Sender{
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
			sendersRepository.GetCall.ReturnSender = models.Sender{
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

			Expect(sendersRepository.GetCall.Conn).To(Equal(conn))
			Expect(sendersRepository.GetCall.SenderID).To(Equal("some-sender-id"))
		})

		Context("failure cases", func() {
			It("validates that a sender id was specified", func() {
				_, err := sendersCollection.Get(conn, "", "some-client-id")
				Expect(err).To(MatchError(collections.NewValidationError("missing sender id")))
			})

			It("validates that a client id was specified", func() {
				_, err := sendersCollection.Get(conn, "some-sender-id", "")
				Expect(err).To(MatchError(collections.NewValidationError("missing client id")))
			})

			It("generates a not found error when the sender does not exist", func() {
				sendersRepository.GetCall.Err = models.RecordNotFoundError("sender not found")

				_, err := sendersCollection.Get(conn, "some-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.NotFoundError{
					Message: "sender not found",
					Err:     models.RecordNotFoundError("sender not found"),
				}))
			})

			It("generates a not found error when the sender belongs to a different client", func() {
				_, err := sendersCollection.Get(conn, "some-sender-id", "mismatch-client-id")
				Expect(err).To(MatchError(collections.NewNotFoundError("sender not found")))
			})

			It("handles unexpected database errors", func() {
				sendersRepository.GetCall.ReturnSender = models.Sender{}
				sendersRepository.GetCall.Err = errors.New("BOOM!")

				_, err := sendersCollection.Get(conn, "some-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.PersistenceError{
					Err: errors.New("BOOM!"),
				}))
			})
		})
	})
})
