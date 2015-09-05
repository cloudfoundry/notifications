package queue_test

import (
	"bytes"
	"errors"
	"log"
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/models"
	"github.com/cloudfoundry-incubator/notifications/v2/queue"
	"github.com/nu7hatch/gouuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JobEnqueuer", func() {
	var (
		enqueuer      queue.JobEnqueuer
		logger        *log.Logger
		buffer        *bytes.Buffer
		gobbleQueue   *mocks.Queue
		conn          *mocks.Connection
		transaction   *mocks.Transaction
		messagesRepo  *mocks.MessagesRepository
		guidGenerator *mocks.GUIDGenerator
		space         cf.CloudControllerSpace
		org           cf.CloudControllerOrganization
		reqReceived   time.Time
	)

	BeforeEach(func() {
		buffer = bytes.NewBuffer([]byte{})
		logger = log.New(buffer, "", 0)
		gobbleQueue = mocks.NewQueue()

		transaction = mocks.NewTransaction()
		conn = mocks.NewConnection()
		conn.TransactionCall.Returns.Transaction = transaction

		guidGenerator = mocks.NewGUIDGenerator()
		guid1 := uuid.UUID([16]byte{0xDE, 0xAD, 0xBE, 0xEF, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55})
		guid2 := uuid.UUID([16]byte{0xDE, 0xAD, 0xBE, 0xEF, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x00, 0x11, 0x22, 0x33, 0x44, 0x56})
		guid3 := uuid.UUID([16]byte{0xDE, 0xAD, 0xBE, 0xEF, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x00, 0x11, 0x22, 0x33, 0x44, 0x57})
		guid4 := uuid.UUID([16]byte{0xDE, 0xAD, 0xBE, 0xEF, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x00, 0x11, 0x22, 0x33, 0x44, 0x58})
		guidGenerator.GenerateCall.Returns.GUIDs = []*uuid.UUID{&guid1, &guid2, &guid3, &guid4}

		messagesRepo = mocks.NewMessagesRepository()
		enqueuer = queue.NewJobEnqueuer(gobbleQueue, guidGenerator.Generate, messagesRepo)
		space = cf.CloudControllerSpace{Name: "the-space"}
		org = cf.CloudControllerOrganization{Name: "the-org"}
		reqReceived, _ = time.Parse(time.RFC3339Nano, "2015-06-08T14:40:12.207187819-07:00")
	})

	Describe("Enqueue", func() {
		It("returns the correct types of responses for users", func() {
			users := []queue.User{{GUID: "user-1"}, {Email: "user-2@example.com"}, {GUID: "user-3"}, {GUID: "user-4"}}
			responses := enqueuer.Enqueue(conn, users, queue.Options{KindID: "the-kind"}, space, org, "the-client", "my-uaa-host", "my.scope", "some-request-id", reqReceived, "some-campaign")

			Expect(responses).To(HaveLen(4))
			Expect(responses).To(ConsistOf([]queue.Response{
				{
					Status:         "queued",
					Recipient:      "user-1",
					NotificationID: "deadbeef-aabb-ccdd-eeff-001122334455",
					VCAPRequestID:  "some-request-id",
				},
				{
					Status:         "queued",
					Recipient:      "user-2@example.com",
					NotificationID: "deadbeef-aabb-ccdd-eeff-001122334456",
					VCAPRequestID:  "some-request-id",
				},
				{
					Status:         "queued",
					Recipient:      "user-3",
					NotificationID: "deadbeef-aabb-ccdd-eeff-001122334457",
					VCAPRequestID:  "some-request-id",
				},
				{
					Status:         "queued",
					Recipient:      "user-4",
					NotificationID: "deadbeef-aabb-ccdd-eeff-001122334458",
					VCAPRequestID:  "some-request-id",
				},
			}))
		})

		It("enqueues jobs with the deliveries", func() {
			users := []queue.User{{GUID: "user-1"}, {GUID: "user-2"}, {GUID: "user-3"}, {GUID: "user-4"}}
			enqueuer.Enqueue(conn, users, queue.Options{}, space, org, "the-client", "my-uaa-host", "my.scope", "some-request-id", reqReceived, "some-campaign")

			var deliveries []queue.Delivery
			for _, job := range gobbleQueue.EnqueueCall.Receives.Jobs {
				var delivery queue.Delivery
				err := job.Unmarshal(&delivery)
				if err != nil {
					panic(err)
				}
				deliveries = append(deliveries, delivery)
			}

			Expect(deliveries).To(HaveLen(4))
			Expect(deliveries).To(ConsistOf([]queue.Delivery{
				{
					JobType:         "v2",
					Options:         queue.Options{},
					UserGUID:        "user-1",
					Space:           space,
					Organization:    org,
					ClientID:        "the-client",
					MessageID:       "deadbeef-aabb-ccdd-eeff-001122334455",
					UAAHost:         "my-uaa-host",
					Scope:           "my.scope",
					VCAPRequestID:   "some-request-id",
					RequestReceived: reqReceived,
					CampaignID:      "some-campaign",
				},
				{
					JobType:         "v2",
					Options:         queue.Options{},
					UserGUID:        "user-2",
					Space:           space,
					Organization:    org,
					ClientID:        "the-client",
					MessageID:       "deadbeef-aabb-ccdd-eeff-001122334456",
					UAAHost:         "my-uaa-host",
					Scope:           "my.scope",
					VCAPRequestID:   "some-request-id",
					RequestReceived: reqReceived,
					CampaignID:      "some-campaign",
				},
				{
					JobType:         "v2",
					Options:         queue.Options{},
					UserGUID:        "user-3",
					Space:           space,
					Organization:    org,
					ClientID:        "the-client",
					MessageID:       "deadbeef-aabb-ccdd-eeff-001122334457",
					UAAHost:         "my-uaa-host",
					Scope:           "my.scope",
					VCAPRequestID:   "some-request-id",
					RequestReceived: reqReceived,
					CampaignID:      "some-campaign",
				},
				{
					JobType:         "v2",
					Options:         queue.Options{},
					UserGUID:        "user-4",
					Space:           space,
					Organization:    org,
					ClientID:        "the-client",
					MessageID:       "deadbeef-aabb-ccdd-eeff-001122334458",
					UAAHost:         "my-uaa-host",
					Scope:           "my.scope",
					VCAPRequestID:   "some-request-id",
					RequestReceived: reqReceived,
					CampaignID:      "some-campaign",
				},
			}))
		})

		It("Inserts a StatusQueued for each of the jobs", func() {
			users := []queue.User{{GUID: "user-1"}, {GUID: "user-2"}, {GUID: "user-3"}, {GUID: "user-4"}}
			enqueuer.Enqueue(conn, users, queue.Options{}, space, org, "the-client", "my-uaa-host", "my.scope", "some-request-id", reqReceived, "some-campaign")

			var messages []models.Message
			for _, call := range messagesRepo.InsertCalls {
				messages = append(messages, call.Receives.Message)
			}

			Expect(messages).To(HaveLen(4))
			Expect(messages).To(ConsistOf([]models.Message{
				{
					ID:         "deadbeef-aabb-ccdd-eeff-001122334455",
					Status:     queue.StatusQueued,
					CampaignID: "some-campaign",
				},
				{
					ID:         "deadbeef-aabb-ccdd-eeff-001122334456",
					Status:     queue.StatusQueued,
					CampaignID: "some-campaign",
				},
				{
					ID:         "deadbeef-aabb-ccdd-eeff-001122334457",
					Status:     queue.StatusQueued,
					CampaignID: "some-campaign",
				},
				{
					ID:         "deadbeef-aabb-ccdd-eeff-001122334458",
					Status:     queue.StatusQueued,
					CampaignID: "some-campaign",
				},
			}))
		})

		Context("using a transaction", func() {
			It("commits the transaction when everything goes well", func() {
				users := []queue.User{{GUID: "user-1"}, {GUID: "user-2"}, {GUID: "user-3"}, {GUID: "user-4"}}
				responses := enqueuer.Enqueue(conn, users, queue.Options{}, space, org, "the-client", "my-uaa-host", "my.scope", "some-request-id", reqReceived, "some-campaign")

				Expect(transaction.BeginCall.WasCalled).To(BeTrue())
				Expect(transaction.CommitCall.WasCalled).To(BeTrue())
				Expect(transaction.RollbackCall.WasCalled).To(BeFalse())

				Expect(responses).ToNot(BeEmpty())
			})

			It("rolls back the transaction when there is an error in message repo upserting", func() {
				messagesRepo.InsertCall.Returns.Error = errors.New("BOOM!")
				users := []queue.User{{GUID: "user-1"}}
				enqueuer.Enqueue(conn, users, queue.Options{}, space, org, "the-client", "my-uaa-host", "my.scope", "some-request-id", reqReceived, "some-campaign")

				Expect(transaction.BeginCall.WasCalled).To(BeTrue())
				Expect(transaction.CommitCall.WasCalled).To(BeFalse())
				Expect(transaction.RollbackCall.WasCalled).To(BeTrue())
			})

			It("returns an empty slice of Response if transaction fails", func() {
				transaction.CommitCall.Returns.Error = errors.New("the commit blew up")

				users := []queue.User{{GUID: "user-1"}, {GUID: "user-2"}, {GUID: "user-3"}, {GUID: "user-4"}}
				responses := enqueuer.Enqueue(conn, users, queue.Options{}, space, org, "the-client", "my-uaa-host", "my.scope", "some-request-id", reqReceived, "some-campaign")

				Expect(transaction.BeginCall.WasCalled).To(BeTrue())
				Expect(transaction.CommitCall.WasCalled).To(BeTrue())
				Expect(transaction.RollbackCall.WasCalled).To(BeFalse())

				Expect(responses).To(Equal([]queue.Response{}))
			})
		})
	})
})
