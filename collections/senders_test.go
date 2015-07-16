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
		sendersRepository.InsertCall.ReturnSender = models.Sender{
			ID:       "some-sender-id",
			Name:     "some-sender",
			ClientID: "some-client-id",
		}

		sendersCollection = collections.NewSendersCollection(sendersRepository)
		conn = fakes.NewConnection()
	})

	Describe("Add", func() {
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
				Expect(err).To(MatchError(collections.ValidationError{
					Err: errors.New("missing sender name"),
				}))
			})

			It("validates that a sender has a client id", func() {
				_, err := sendersCollection.Add(conn, collections.Sender{
					Name: "some-sender",
				})
				Expect(err).To(MatchError(collections.ValidationError{
					Err: errors.New("missing sender client_id"),
				}))
			})
		})
	})
})
