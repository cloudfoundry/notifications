package strategies_test

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EmailStrategy", func() {
	var emailStrategy strategies.EmailStrategy

	Describe("Dispatch", func() {
		var (
			mailer          *fakes.Mailer
			conn            *fakes.Connection
			requestReceived time.Time
		)

		BeforeEach(func() {
			mailer = fakes.NewMailer()
			emailStrategy = strategies.NewEmailStrategy(mailer)
			conn = fakes.NewConnection()
			requestReceived, _ = time.Parse(time.RFC3339Nano, "2015-06-08T14:37:35.181067085-07:00")
		})

		It("calls Deliver on it's mailer with proper arguments", func() {
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

			Expect(mailer.DeliverCall.Args.Connection).To(Equal(conn))
			Expect(mailer.DeliverCall.Args.Users).To(Equal(users))
			Expect(mailer.DeliverCall.Args.Options).To(Equal(postal.Options{
				ReplyTo:           "reply-to@example.com",
				Subject:           "this is the subject",
				KindDescription:   "description of a kind",
				SourceDescription: "description of a client",
				Text:              "email text",
				HTML: postal.HTML{
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
			Expect(mailer.DeliverCall.Args.Space).To(Equal(cf.CloudControllerSpace{}))
			Expect(mailer.DeliverCall.Args.Org).To(Equal(cf.CloudControllerOrganization{}))
			Expect(mailer.DeliverCall.Args.Client).To(Equal("some-client-id"))
			Expect(mailer.DeliverCall.Args.Scope).To(Equal(""))
			Expect(mailer.DeliverCall.Args.VCAPRequestID).To(Equal("some-vcap-request-id"))
			Expect(mailer.DeliverCall.Args.RequestReceived).To(Equal(requestReceived))
		})
	})
})
