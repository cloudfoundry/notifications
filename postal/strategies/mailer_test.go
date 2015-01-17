package strategies_test

import (
	"bytes"
	"errors"
	"log"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mailer", func() {
	var mailer strategies.Mailer
	var logger *log.Logger
	var buffer *bytes.Buffer
	var queue *fakes.Queue
	var conn *fakes.DBConn
	var space cf.CloudControllerSpace
	var org cf.CloudControllerOrganization
	var messagesRepo *fakes.MessagesRepo

	BeforeEach(func() {
		buffer = bytes.NewBuffer([]byte{})
		logger = log.New(buffer, "", 0)
		queue = fakes.NewQueue()
		conn = fakes.NewDBConn()
		messagesRepo = fakes.NewMessagesRepo()
		mailer = strategies.NewMailer(queue, fakes.NewIncrementingGUIDGenerator().Generate, messagesRepo)
		space = cf.CloudControllerSpace{Name: "the-space"}
		org = cf.CloudControllerOrganization{Name: "the-org"}
	})

	Describe("Deliver", func() {
		It("returns the correct types of responses for users", func() {
			users := map[string]uaa.User{
				"user-1": {ID: "user-1", Emails: []string{"user-1@example.com"}},
				"user-2": {},
				"user-3": {ID: "user-3"},
				"user-4": {ID: "user-4", Emails: []string{"user-4"}},
			}
			responses := mailer.Deliver(conn, postal.Templates{}, users, postal.Options{KindID: "the-kind"}, space, org, "the-client", "my.scope")

			Expect(len(responses)).To(Equal(4))

			guidsSeen := map[string]string{}

			for _, r := range responses {
				Expect(r.Status).To(Equal("queued"))
				Expect(users).To(HaveKey(r.Recipient))
				Expect(guidsSeen).NotTo(HaveKey(r.NotificationID))
				guidsSeen[r.NotificationID] = r.NotificationID
			}
		})

		It("enqueues jobs with the deliveries", func() {
			users := map[string]uaa.User{
				"user-1": {ID: "user-1", Emails: []string{"user-1@example.com"}},
				"user-2": {},
				"user-3": {ID: "user-3"},
				"user-4": {ID: "user-4", Emails: []string{"user-4"}},
			}
			mailer.Deliver(conn, postal.Templates{}, users, postal.Options{}, space, org, "the-client", "my.scope")

			for _ = range users {
				job := <-queue.Reserve("me")
				var delivery postal.Delivery
				err := job.Unmarshal(&delivery)
				if err != nil {
					panic(err)
				}

				user := users[delivery.UserGUID]
				Expect(delivery).To(Equal(postal.Delivery{
					User: user,
					Options: postal.Options{
						ReplyTo:           "",
						Subject:           "",
						KindDescription:   "",
						SourceDescription: "",
						Text:              "",
						HTML:              postal.HTML{},
						KindID:            "",
					},
					UserGUID:     delivery.UserGUID,
					Space:        space,
					Organization: org,
					ClientID:     "the-client",
					Templates:    postal.Templates{Subject: "", Text: "", HTML: ""},
					MessageID:    delivery.MessageID,
					Scope:        "my.scope",
				}))
			}
		})

		It("Upserts a StatusQueued for each of the jobs", func() {
			users := map[string]uaa.User{
				"user-1": {ID: "user-1", Emails: []string{"user-1@example.com"}},
				"user-2": {},
				"user-3": {ID: "user-3"},
				"user-4": {ID: "user-4", Emails: []string{"user-4"}},
			}
			mailer.Deliver(conn, postal.Templates{}, users, postal.Options{}, space, org, "the-client", "my.scope")

			for _ = range users {
				job := <-queue.Reserve("me")
				var delivery postal.Delivery
				err := job.Unmarshal(&delivery)
				if err != nil {
					panic(err)
				}

				message, err := messagesRepo.FindByID(conn, delivery.MessageID)
				Expect(err).ToNot(HaveOccurred())
				Expect(message.Status).To(Equal(postal.StatusQueued))
			}

		})

		Context("using a transaction", func() {
			It("commits the transaction when everything goes well", func() {
				users := map[string]uaa.User{
					"user-1": {ID: "user-1", Emails: []string{"user-1@example.com"}},
					"user-2": {},
					"user-3": {ID: "user-3"},
					"user-4": {ID: "user-4", Emails: []string{"user-4"}},
				}
				responses := mailer.Deliver(conn, postal.Templates{}, users, postal.Options{}, space, org, "the-client", "my.scope")

				Expect(conn.BeginWasCalled).To(BeTrue())
				Expect(conn.CommitWasCalled).To(BeTrue())
				Expect(conn.RollbackWasCalled).To(BeFalse())
				Expect(responses).ToNot(BeEmpty())
			})

			It("rolls back the transaction when there is an error in queueing", func() {
				queue.EnqueueError = errors.New("BOOM!")

				users := map[string]uaa.User{
					"user-1": {ID: "user-1", Emails: []string{"user-1@example.com"}},
				}
				mailer.Deliver(conn, postal.Templates{}, users, postal.Options{}, space, org, "the-client", "my.scope")

				Expect(conn.BeginWasCalled).To(BeTrue())
				Expect(conn.CommitWasCalled).To(BeFalse())
				Expect(conn.RollbackWasCalled).To(BeTrue())
			})

			It("rolls back the transaction when there is an error in message repo upserting", func() {
				messagesRepo.UpsertError = errors.New("BOOM!")

				users := map[string]uaa.User{
					"user-1": {ID: "user-1", Emails: []string{"user-1@example.com"}},
				}
				mailer.Deliver(conn, postal.Templates{}, users, postal.Options{}, space, org, "the-client", "my.scope")

				Expect(conn.BeginWasCalled).To(BeTrue())
				Expect(conn.CommitWasCalled).To(BeFalse())
				Expect(conn.RollbackWasCalled).To(BeTrue())
			})

			It("returns an empty []Response{} if transaction fails", func() {
				conn.CommitError = "the commit blew up"
				users := map[string]uaa.User{
					"user-1": {ID: "user-1", Emails: []string{"user-1@example.com"}},
					"user-2": {},
					"user-3": {ID: "user-3"},
					"user-4": {ID: "user-4", Emails: []string{"user-4"}},
				}
				responses := mailer.Deliver(conn, postal.Templates{}, users, postal.Options{}, space, org, "the-client", "my.scope")

				Expect(conn.BeginWasCalled).To(BeTrue())
				Expect(conn.CommitWasCalled).To(BeTrue())
				Expect(conn.RollbackWasCalled).To(BeFalse())
				Expect(responses).To(Equal([]strategies.Response{}))
			})

		})
	})
})
