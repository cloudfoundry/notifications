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

	BeforeEach(func() {
		buffer = bytes.NewBuffer([]byte{})
		logger = log.New(buffer, "", 0)
		queue = fakes.NewQueue()
		conn = fakes.NewDBConn()
		mailer = strategies.NewMailer(queue, fakes.GUIDGenerator)
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
			responses := mailer.Deliver(conn, postal.Templates{}, users, postal.Options{KindID: "the-kind"}, space, org, "the-client")

			Expect(len(responses)).To(Equal(4))
			Expect(responses).To(ContainElement(strategies.Response{
				Status:         "queued",
				Recipient:      "user-1",
				NotificationID: "deadbeef-aabb-ccdd-eeff-001122334455",
				Email:          "user-1@example.com",
			}))

			Expect(responses).To(ContainElement(strategies.Response{
				Status:         "queued",
				Recipient:      "user-2",
				NotificationID: "deadbeef-aabb-ccdd-eeff-001122334455",
				Email:          "",
			}))

			Expect(responses).To(ContainElement(strategies.Response{
				Status:         "queued",
				Recipient:      "user-3",
				NotificationID: "deadbeef-aabb-ccdd-eeff-001122334455",
				Email:          "",
			}))

			Expect(responses).To(ContainElement(strategies.Response{
				Status:         "queued",
				Recipient:      "user-4",
				NotificationID: "deadbeef-aabb-ccdd-eeff-001122334455",
				Email:          "user-4",
			}))
		})

		It("enqueues jobs with the deliveries", func() {
			users := map[string]uaa.User{
				"user-1": {ID: "user-1", Emails: []string{"user-1@example.com"}},
				"user-2": {},
				"user-3": {ID: "user-3"},
				"user-4": {ID: "user-4", Emails: []string{"user-4"}},
			}
			mailer.Deliver(conn, postal.Templates{}, users, postal.Options{}, space, org, "the-client")

			for userGUID, user := range users {
				job := <-queue.Reserve("me")
				var delivery postal.Delivery
				err := job.Unmarshal(&delivery)
				if err != nil {
					panic(err)
				}
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
					UserGUID:     userGUID,
					Space:        space,
					Organization: org,
					ClientID:     "the-client",
					Templates:    postal.Templates{Subject: "", Text: "", HTML: ""},
					MessageID:    "deadbeef-aabb-ccdd-eeff-001122334455",
				}))
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
				mailer.Deliver(conn, postal.Templates{}, users, postal.Options{}, space, org, "the-client")

				Expect(conn.BeginWasCalled).To(BeTrue())
				Expect(conn.CommitWasCalled).To(BeTrue())
				Expect(conn.RollbackWasCalled).To(BeFalse())
			})

			It("rolls back the transaction when there is an error", func() {
				queue.EnqueueError = errors.New("BOOM!")

				users := map[string]uaa.User{
					"user-1": {ID: "user-1", Emails: []string{"user-1@example.com"}},
				}
				mailer.Deliver(conn, postal.Templates{}, users, postal.Options{}, space, org, "the-client")

				Expect(conn.BeginWasCalled).To(BeTrue())
				Expect(conn.CommitWasCalled).To(BeFalse())
				Expect(conn.RollbackWasCalled).To(BeTrue())
			})
		})
	})
})
