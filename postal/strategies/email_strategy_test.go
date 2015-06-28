package strategies_test

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EmailStrategy", func() {
	var emailStrategy strategies.EmailStrategy

	Describe("Dispatch", func() {
		var (
			enqueuer        *fakes.Enqueuer
			conn            *fakes.Connection
			requestReceived time.Time
		)

		BeforeEach(func() {
			enqueuer = fakes.NewEnqueuer()
			emailStrategy = strategies.NewEmailStrategy(enqueuer)
			conn = fakes.NewConnection()
			requestReceived, _ = time.Parse(time.RFC3339Nano, "2015-06-08T14:37:35.181067085-07:00")
		})

		It("calls Enqueue on it's enqueuer with proper arguments", func() {
			emailStrategy.Dispatch(strategies.Dispatch{
				Connection: conn,
				Client: strategies.Client{
					ID:          "some-client-id",
					Description: "description of a client",
				},
				Kind: strategies.Kind{
					ID:          "some-kind-id",
					Description: "description of a kind",
				},
				Message: strategies.Message{
					ReplyTo: "reply-to@example.com",
					Subject: "this is the subject",
					To:      "dr@strangelove.com",
					Text:    "email text",
					HTML: strategies.HTML{
						BodyContent:    "some html body content",
						BodyAttributes: "some html body attributes",
						Head:           "the html head tag",
						Doctype:        "the html doctype",
					},
				},
				VCAPRequest: strategies.VCAPRequest{
					ID:          "some-vcap-request-id",
					ReceiptTime: requestReceived,
				},
			})

			users := []strategies.User{{Email: "dr@strangelove.com"}}

			Expect(enqueuer.EnqueueCall.Args.Connection).To(Equal(conn))
			Expect(enqueuer.EnqueueCall.Args.Users).To(Equal(users))
			Expect(enqueuer.EnqueueCall.Args.Options).To(Equal(strategies.Options{
				ReplyTo:           "reply-to@example.com",
				Subject:           "this is the subject",
				KindDescription:   "description of a kind",
				SourceDescription: "description of a client",
				Text:              "email text",
				HTML: strategies.HTML{
					BodyContent:    "some html body content",
					BodyAttributes: "some html body attributes",
					Head:           "the html head tag",
					Doctype:        "the html doctype",
				},
				KindID:      "some-kind-id",
				To:          "dr@strangelove.com",
				Role:        "",
				Endorsement: strategies.EmailEndorsement,
			}))
			Expect(enqueuer.EnqueueCall.Args.Space).To(Equal(cf.CloudControllerSpace{}))
			Expect(enqueuer.EnqueueCall.Args.Org).To(Equal(cf.CloudControllerOrganization{}))
			Expect(enqueuer.EnqueueCall.Args.Client).To(Equal("some-client-id"))
			Expect(enqueuer.EnqueueCall.Args.Scope).To(Equal(""))
			Expect(enqueuer.EnqueueCall.Args.VCAPRequestID).To(Equal("some-vcap-request-id"))
			Expect(enqueuer.EnqueueCall.Args.RequestReceived).To(Equal(requestReceived))
		})
	})
})
