package strategies_test

import (
	"bytes"
	"errors"
	"log"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mailer", func() {
	var mailer strategies.Mailer
	var logger *log.Logger
	var buffer *bytes.Buffer
	var queue *fakes.Queue
	var conn *fakes.Connection
	var space cf.CloudControllerSpace
	var org cf.CloudControllerOrganization
	var messagesRepo *fakes.MessagesRepo

	BeforeEach(func() {
		buffer = bytes.NewBuffer([]byte{})
		logger = log.New(buffer, "", 0)
		queue = fakes.NewQueue()
		conn = fakes.NewConnection()
		messagesRepo = fakes.NewMessagesRepo()
		mailer = strategies.NewMailer(queue, fakes.NewIncrementingGUIDGenerator().Generate, messagesRepo)
		space = cf.CloudControllerSpace{Name: "the-space"}
		org = cf.CloudControllerOrganization{Name: "the-org"}
	})

	Describe("Deliver", func() {
		It("returns the correct types of responses for users", func() {
			users := []strategies.User{{GUID: "user-1"}, {Email: "user-2@example.com"}, {GUID: "user-3"}, {GUID: "user-4"}}
			responses := mailer.Deliver(conn, users, postal.Options{KindID: "the-kind"}, space, org, "the-client", "my.scope", "some-request-id")

			Expect(responses).To(HaveLen(4))
			Expect(responses).To(ConsistOf([]strategies.Response{
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
			users := []strategies.User{{GUID: "user-1"}, {GUID: "user-2"}, {GUID: "user-3"}, {GUID: "user-4"}}
			mailer.Deliver(conn, users, postal.Options{}, space, org, "the-client", "my.scope", "some-request-id")

			var deliveries []postal.Delivery
			for _ = range users {
				job := <-queue.Reserve("me")
				var delivery postal.Delivery
				err := job.Unmarshal(&delivery)
				if err != nil {
					panic(err)
				}
				deliveries = append(deliveries, delivery)
			}

			Expect(deliveries).To(HaveLen(4))
			Expect(deliveries).To(ConsistOf([]postal.Delivery{
				{
					Options:       postal.Options{},
					UserGUID:      "user-1",
					Space:         space,
					Organization:  org,
					ClientID:      "the-client",
					MessageID:     "deadbeef-aabb-ccdd-eeff-001122334455",
					Scope:         "my.scope",
					VCAPRequestID: "some-request-id",
				},
				{
					Options:       postal.Options{},
					UserGUID:      "user-2",
					Space:         space,
					Organization:  org,
					ClientID:      "the-client",
					MessageID:     "deadbeef-aabb-ccdd-eeff-001122334456",
					Scope:         "my.scope",
					VCAPRequestID: "some-request-id",
				},
				{
					Options:       postal.Options{},
					UserGUID:      "user-3",
					Space:         space,
					Organization:  org,
					ClientID:      "the-client",
					MessageID:     "deadbeef-aabb-ccdd-eeff-001122334457",
					Scope:         "my.scope",
					VCAPRequestID: "some-request-id",
				},
				{
					Options:       postal.Options{},
					UserGUID:      "user-4",
					Space:         space,
					Organization:  org,
					ClientID:      "the-client",
					MessageID:     "deadbeef-aabb-ccdd-eeff-001122334458",
					Scope:         "my.scope",
					VCAPRequestID: "some-request-id",
				},
			}))
		})

		It("Upserts a StatusQueued for each of the jobs", func() {
			users := []strategies.User{{GUID: "user-1"}, {GUID: "user-2"}, {GUID: "user-3"}, {GUID: "user-4"}}
			mailer.Deliver(conn, users, postal.Options{}, space, org, "the-client", "my.scope", "some-request-id")

			var statuses []string
			for _ = range users {
				job := <-queue.Reserve("me")
				var delivery postal.Delivery
				err := job.Unmarshal(&delivery)
				if err != nil {
					panic(err)
				}

				message, err := messagesRepo.FindByID(conn, delivery.MessageID)
				if err != nil {
					panic(err)
				}

				statuses = append(statuses, message.Status)
			}

			Expect(statuses).To(HaveLen(4))
			Expect(statuses).To(ConsistOf([]string{postal.StatusQueued, postal.StatusQueued, postal.StatusQueued, postal.StatusQueued}))
		})

		Context("using a transaction", func() {
			It("commits the transaction when everything goes well", func() {
				users := []strategies.User{{GUID: "user-1"}, {GUID: "user-2"}, {GUID: "user-3"}, {GUID: "user-4"}}
				responses := mailer.Deliver(conn, users, postal.Options{}, space, org, "the-client", "my.scope", "some-request-id")

				Expect(conn.BeginWasCalled).To(BeTrue())
				Expect(conn.CommitWasCalled).To(BeTrue())
				Expect(conn.RollbackWasCalled).To(BeFalse())
				Expect(responses).ToNot(BeEmpty())
			})

			It("rolls back the transaction when there is an error in message repo upserting", func() {
				messagesRepo.UpsertError = errors.New("BOOM!")
				users := []strategies.User{{GUID: "user-1"}}
				mailer.Deliver(conn, users, postal.Options{}, space, org, "the-client", "my.scope", "some-request-id")

				Expect(conn.BeginWasCalled).To(BeTrue())
				Expect(conn.CommitWasCalled).To(BeFalse())
				Expect(conn.RollbackWasCalled).To(BeTrue())
			})

			It("returns an empty []Response{} if transaction fails", func() {
				conn.CommitError = "the commit blew up"
				users := []strategies.User{{GUID: "user-1"}, {GUID: "user-2"}, {GUID: "user-3"}, {GUID: "user-4"}}
				responses := mailer.Deliver(conn, users, postal.Options{}, space, org, "the-client", "my.scope", "some-request-id")

				Expect(conn.BeginWasCalled).To(BeTrue())
				Expect(conn.CommitWasCalled).To(BeTrue())
				Expect(conn.RollbackWasCalled).To(BeFalse())
				Expect(responses).To(Equal([]strategies.Response{}))
			})
		})
	})
})
