package services_test

import (
	"errors"
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Enqueuer", func() {
	var (
		enqueuer          services.Enqueuer
		queue             *mocks.Queue
		gobbleInitializer *mocks.GobbleInitializer
		conn              *mocks.Connection
		transaction       *mocks.Transaction
		space             cf.CloudControllerSpace
		org               cf.CloudControllerOrganization
		reqReceived       time.Time
		messagesRepo      *mocks.MessagesRepo
	)

	BeforeEach(func() {
		queue = mocks.NewQueue()

		transaction = mocks.NewTransaction()
		conn = mocks.NewConnection()

		conn.TransactionCall.Returns.Transaction = transaction
		transaction.Connection = conn

		gobbleInitializer = mocks.NewGobbleInitializer()

		space = cf.CloudControllerSpace{Name: "the-space"}
		org = cf.CloudControllerOrganization{Name: "the-org"}
		reqReceived, _ = time.Parse(time.RFC3339Nano, "2015-06-08T14:40:12.207187819-07:00")

		messagesRepo = mocks.NewMessagesRepo()
		messagesRepo.UpsertCall.Returns.Messages = []models.Message{
			{
				ID:     "first-random-guid",
				Status: services.StatusQueued,
			},
			{
				ID:     "second-random-guid",
				Status: services.StatusQueued,
			},
			{
				ID:     "third-random-guid",
				Status: services.StatusQueued,
			},
			{
				ID:     "fourth-random-guid",
				Status: services.StatusQueued,
			},
		}

		enqueuer = services.NewEnqueuer(queue, messagesRepo, gobbleInitializer)
	})

	Describe("Enqueue", func() {
		It("returns the correct types of responses for users", func() {
			users := []services.User{{GUID: "user-1"}, {Email: "user-2@example.com"}, {GUID: "user-3"}, {GUID: "user-4"}}
			responses, err := enqueuer.Enqueue(conn, users, services.Options{KindID: "the-kind"}, space, org, "the-client", "my-uaa-host", "my.scope", "some-request-id", reqReceived)

			Expect(err).ToNot(HaveOccurred())
			Expect(responses).To(HaveLen(4))
			Expect(responses).To(ConsistOf([]services.Response{
				{
					Status:         "queued",
					Recipient:      "user-1",
					NotificationID: "first-random-guid",
					VCAPRequestID:  "some-request-id",
				},
				{
					Status:         "queued",
					Recipient:      "user-2@example.com",
					NotificationID: "second-random-guid",
					VCAPRequestID:  "some-request-id",
				},
				{
					Status:         "queued",
					Recipient:      "user-3",
					NotificationID: "third-random-guid",
					VCAPRequestID:  "some-request-id",
				},
				{
					Status:         "queued",
					Recipient:      "user-4",
					NotificationID: "fourth-random-guid",
					VCAPRequestID:  "some-request-id",
				},
			}))
		})

		It("enqueues jobs with the deliveries", func() {
			users := []services.User{
				{GUID: "user-1"},
				{GUID: "user-2"},
				{GUID: "user-3"},
				{GUID: "user-4"},
			}
			enqueuer.Enqueue(conn, users, services.Options{}, space, org, "the-client", "my-uaa-host", "my.scope", "some-request-id", reqReceived)

			var deliveries []services.Delivery
			for _, job := range queue.EnqueueCall.Receives.Jobs {
				var delivery services.Delivery
				err := job.Unmarshal(&delivery)
				if err != nil {
					panic(err)
				}
				deliveries = append(deliveries, delivery)
			}

			Expect(deliveries).To(HaveLen(4))
			Expect(deliveries).To(ConsistOf([]services.Delivery{
				{
					Options:         services.Options{},
					UserGUID:        "user-1",
					Space:           space,
					Organization:    org,
					ClientID:        "the-client",
					MessageID:       "first-random-guid",
					UAAHost:         "my-uaa-host",
					Scope:           "my.scope",
					VCAPRequestID:   "some-request-id",
					RequestReceived: reqReceived,
				},
				{
					Options:         services.Options{},
					UserGUID:        "user-2",
					Space:           space,
					Organization:    org,
					ClientID:        "the-client",
					MessageID:       "second-random-guid",
					UAAHost:         "my-uaa-host",
					Scope:           "my.scope",
					VCAPRequestID:   "some-request-id",
					RequestReceived: reqReceived,
				},
				{
					Options:         services.Options{},
					UserGUID:        "user-3",
					Space:           space,
					Organization:    org,
					ClientID:        "the-client",
					MessageID:       "third-random-guid",
					UAAHost:         "my-uaa-host",
					Scope:           "my.scope",
					VCAPRequestID:   "some-request-id",
					RequestReceived: reqReceived,
				},
				{
					Options:         services.Options{},
					UserGUID:        "user-4",
					Space:           space,
					Organization:    org,
					ClientID:        "the-client",
					MessageID:       "fourth-random-guid",
					UAAHost:         "my-uaa-host",
					Scope:           "my.scope",
					VCAPRequestID:   "some-request-id",
					RequestReceived: reqReceived,
				},
			}))
		})

		It("upserts a StatusQueued for each of the jobs", func() {
			users := []services.User{{GUID: "user-1"}, {GUID: "user-2"}, {GUID: "user-3"}, {GUID: "user-4"}}
			enqueuer.Enqueue(conn, users, services.Options{}, space, org, "the-client", "my-uaa-host", "my.scope", "some-request-id", reqReceived)

			messages := messagesRepo.UpsertCall.Receives.Messages
			Expect(messages).To(HaveLen(4))
			Expect(messages).To(Equal([]models.Message{
				{Status: services.StatusQueued},
				{Status: services.StatusQueued},
				{Status: services.StatusQueued},
				{Status: services.StatusQueued},
			}))
		})

		Context("using a transaction", func() {
			var users []services.User

			BeforeEach(func() {
				users = []services.User{
					{GUID: "user-1"},
					{GUID: "user-2"},
					{GUID: "user-3"},
					{GUID: "user-4"},
				}
			})

			It("initializes the DbMap", func() {
				enqueuer.Enqueue(conn, users, services.Options{}, space, org, "the-client", "my-uaa-host", "my.scope", "some-request-id", reqReceived)

				isSamePtr := (gobbleInitializer.InitializeDBMapCall.Receives.DbMap == transaction.GetDbMapCall.Returns.DbMap)
				Expect(isSamePtr).To(BeTrue())
				Expect(transaction.GetDbMapCall.WasCalled).To(BeTrue())
			})

			It("commits the transaction when everything goes well", func() {
				responses, err := enqueuer.Enqueue(conn, users, services.Options{}, space, org, "the-client", "my-uaa-host", "my.scope", "some-request-id", reqReceived)

				Expect(err).ToNot(HaveOccurred())
				Expect(transaction.BeginCall.WasCalled).To(BeTrue())
				Expect(transaction.CommitCall.WasCalled).To(BeTrue())
				Expect(transaction.RollbackCall.WasCalled).To(BeFalse())

				Expect(responses).ToNot(BeEmpty())
				Expect(err).ToNot(HaveOccurred())
			})

			It("rolls back the transaction when there is an error in message repo upserting", func() {
				messagesRepo.UpsertCall.Returns.Error = errors.New("BOOM!")
				_, err := enqueuer.Enqueue(conn, users, services.Options{}, space, org, "the-client", "my-uaa-host", "my.scope", "some-request-id", reqReceived)

				Expect(transaction.BeginCall.WasCalled).To(BeTrue())
				Expect(transaction.CommitCall.WasCalled).To(BeFalse())
				Expect(transaction.RollbackCall.WasCalled).To(BeTrue())
				Expect(err).To(HaveOccurred())
			})

			It("rolls back the transaction when there is an error in enqueuing", func() {
				queue.EnqueueCall.Returns.Error = errors.New("BOOM!")
				_, err := enqueuer.Enqueue(conn, users, services.Options{}, space, org, "the-client", "my-uaa-host", "my.scope", "some-request-id", reqReceived)

				Expect(transaction.BeginCall.WasCalled).To(BeTrue())
				Expect(transaction.CommitCall.WasCalled).To(BeFalse())
				Expect(transaction.RollbackCall.WasCalled).To(BeTrue())
				Expect(err).To(HaveOccurred())
			})

			It("uses the same transaction for the queue as it did for the messages repo", func() {
				enqueuer.Enqueue(conn, users, services.Options{}, space, org, "the-client", "my-uaa-host", "my.scope", "some-request-id", reqReceived)

				Expect(messagesRepo.UpsertCall.Receives.Connection).To(Equal(transaction))
				Expect(queue.EnqueueCall.Receives.Connection).To(Equal(transaction))
			})

			It("does not commit the transaction until the jobs have been queued", func() {
				queue.EnqueueCall.Hook = func() {
					Expect(transaction.CommitCall.WasCalled).To(BeFalse())
				}

				enqueuer.Enqueue(conn, users, services.Options{}, space, org, "the-client", "my-uaa-host", "my.scope", "some-request-id", reqReceived)
			})

			It("returns an empty slice of Response if transaction fails", func() {
				transaction.CommitCall.Returns.Error = errors.New("the commit blew up")
				responses, err := enqueuer.Enqueue(conn, users, services.Options{}, space, org, "the-client", "my-uaa-host", "my.scope", "some-request-id", reqReceived)

				Expect(transaction.BeginCall.WasCalled).To(BeTrue())
				Expect(transaction.CommitCall.WasCalled).To(BeTrue())
				Expect(transaction.RollbackCall.WasCalled).To(BeFalse())

				Expect(responses).To(Equal([]services.Response{}))
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
