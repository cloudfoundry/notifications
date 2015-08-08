package services_test

import (
	"reflect"
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserStrategy", func() {
	var (
		strategy        services.UserStrategy
		enqueuer        *fakes.Enqueuer
		conn            *fakes.Connection
		requestReceived time.Time
	)

	BeforeEach(func() {
		requestReceived, _ = time.Parse(time.RFC3339Nano, "2015-06-08T14:37:35.181067085-07:00")
		conn = fakes.NewConnection()
		enqueuer = fakes.NewEnqueuer()
		strategy = services.NewUserStrategy(enqueuer)
	})

	Describe("Dispatch", func() {
		It("calls enqueuer.Enqueue with the correct arguments for a user", func() {
			_, err := strategy.Dispatch(services.Dispatch{
				GUID:       "user-123",
				Connection: conn,
				Message: services.DispatchMessage{
					To:      "dr@strangelove.com",
					ReplyTo: "reply-to@example.com",
					Subject: "this is the subject",
					Text:    "Please make sure to leave your bottle in a place that is safe and dry",
					HTML: services.HTML{
						BodyContent:    "<p>The water bottle needs to be safe and dry</p>",
						BodyAttributes: "some-html-body-attributes",
						Head:           "<head></head>",
						Doctype:        "<html>",
					},
				},
				UAAHost: "uaa",
				Kind: services.DispatchKind{
					ID:          "forgot_waterbottle",
					Description: "Water Bottle Reminder",
				},
				Client: services.DispatchClient{
					ID:          "mister-client",
					Description: "The Water Bottle System",
				},
				VCAPRequest: services.DispatchVCAPRequest{
					ID:          "some-vcap-request-id",
					ReceiptTime: requestReceived,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			users := []services.User{{GUID: "user-123"}}

			Expect(reflect.ValueOf(enqueuer.EnqueueCall.Args.Connection).Pointer()).To(Equal(reflect.ValueOf(conn).Pointer()))
			Expect(enqueuer.EnqueueCall.Args.Users).To(Equal(users))
			Expect(enqueuer.EnqueueCall.Args.Options).To(Equal(services.Options{
				ReplyTo:           "reply-to@example.com",
				Subject:           "this is the subject",
				To:                "dr@strangelove.com",
				KindID:            "forgot_waterbottle",
				KindDescription:   "Water Bottle Reminder",
				SourceDescription: "The Water Bottle System",
				Text:              "Please make sure to leave your bottle in a place that is safe and dry",
				HTML: services.HTML{
					BodyContent:    "<p>The water bottle needs to be safe and dry</p>",
					BodyAttributes: "some-html-body-attributes",
					Head:           "<head></head>",
					Doctype:        "<html>",
				},
				Endorsement: services.UserEndorsement,
			}))
			Expect(enqueuer.EnqueueCall.Args.Space).To(Equal(cf.CloudControllerSpace{}))
			Expect(enqueuer.EnqueueCall.Args.Org).To(Equal(cf.CloudControllerOrganization{}))
			Expect(enqueuer.EnqueueCall.Args.Client).To(Equal("mister-client"))
			Expect(enqueuer.EnqueueCall.Args.Scope).To(Equal(""))
			Expect(enqueuer.EnqueueCall.Args.UAAHost).To(Equal("uaa"))
			Expect(enqueuer.EnqueueCall.Args.VCAPRequestID).To(Equal("some-vcap-request-id"))
			Expect(enqueuer.EnqueueCall.Args.RequestReceived).To(Equal(requestReceived))
		})
	})
})
