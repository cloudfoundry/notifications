package strategies_test

import (
	"reflect"
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserStrategy", func() {
	var (
		strategy        strategies.UserStrategy
		enqueuer        *fakes.Enqueuer
		conn            *fakes.Connection
		requestReceived time.Time
	)

	BeforeEach(func() {
		requestReceived, _ = time.Parse(time.RFC3339Nano, "2015-06-08T14:37:35.181067085-07:00")
		conn = fakes.NewConnection()
		enqueuer = fakes.NewEnqueuer()
		strategy = strategies.NewUserStrategy(enqueuer)
	})

	Describe("Dispatch", func() {
		It("calls enqueuer.Enqueue with the correct arguments for a user", func() {
			_, err := strategy.Dispatch(strategies.Dispatch{
				GUID:       "user-123",
				Connection: conn,
				Message: strategies.Message{
					To:      "dr@strangelove.com",
					ReplyTo: "reply-to@example.com",
					Subject: "this is the subject",
					Text:    "Please make sure to leave your bottle in a place that is safe and dry",
					HTML: strategies.HTML{
						BodyContent:    "<p>The water bottle needs to be safe and dry</p>",
						BodyAttributes: "some-html-body-attributes",
						Head:           "<head></head>",
						Doctype:        "<html>",
					},
				},
				Kind: strategies.Kind{
					ID:          "forgot_waterbottle",
					Description: "Water Bottle Reminder",
				},
				Client: strategies.Client{
					ID:          "mister-client",
					Description: "The Water Bottle System",
				},
				VCAPRequest: strategies.VCAPRequest{
					ID:          "some-vcap-request-id",
					ReceiptTime: requestReceived,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			users := []strategies.User{{GUID: "user-123"}}

			Expect(reflect.ValueOf(enqueuer.EnqueueCall.Args.Connection).Pointer()).To(Equal(reflect.ValueOf(conn).Pointer()))
			Expect(enqueuer.EnqueueCall.Args.Users).To(Equal(users))
			Expect(enqueuer.EnqueueCall.Args.Options).To(Equal(strategies.Options{
				ReplyTo:           "reply-to@example.com",
				Subject:           "this is the subject",
				To:                "dr@strangelove.com",
				KindID:            "forgot_waterbottle",
				KindDescription:   "Water Bottle Reminder",
				SourceDescription: "The Water Bottle System",
				Text:              "Please make sure to leave your bottle in a place that is safe and dry",
				HTML: strategies.HTML{
					BodyContent:    "<p>The water bottle needs to be safe and dry</p>",
					BodyAttributes: "some-html-body-attributes",
					Head:           "<head></head>",
					Doctype:        "<html>",
				},
				Endorsement: strategies.UserEndorsement,
			}))
			Expect(enqueuer.EnqueueCall.Args.Space).To(Equal(cf.CloudControllerSpace{}))
			Expect(enqueuer.EnqueueCall.Args.Org).To(Equal(cf.CloudControllerOrganization{}))
			Expect(enqueuer.EnqueueCall.Args.Client).To(Equal("mister-client"))
			Expect(enqueuer.EnqueueCall.Args.Scope).To(Equal(""))
			Expect(enqueuer.EnqueueCall.Args.VCAPRequestID).To(Equal("some-vcap-request-id"))
			Expect(enqueuer.EnqueueCall.Args.RequestReceived).To(Equal(requestReceived))
		})
	})
})
