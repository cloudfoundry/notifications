package services_test

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EmailStrategy", func() {
	var emailStrategy services.EmailStrategy

	Describe("Dispatch", func() {
		var (
			v1Enqueuer      *mocks.Enqueuer
			v2Enqueuer      *mocks.V2Enqueuer
			conn            *mocks.Connection
			requestReceived time.Time
		)

		BeforeEach(func() {
			v1Enqueuer = mocks.NewEnqueuer()
			v2Enqueuer = mocks.NewV2Enqueuer()
			emailStrategy = services.NewEmailStrategy(v1Enqueuer, v2Enqueuer)
			conn = mocks.NewConnection()
			requestReceived, _ = time.Parse(time.RFC3339Nano, "2015-06-08T14:37:35.181067085-07:00")
		})

		Context("when the dispatch JobType is unspecified", func() {
			It("calls Enqueue on it's enqueuer with proper arguments", func() {
				emailStrategy.Dispatch(services.Dispatch{
					Connection: conn,
					Client: services.DispatchClient{
						ID:          "some-client-id",
						Description: "description of a client",
					},
					Kind: services.DispatchKind{
						ID:          "some-kind-id",
						Description: "description of a kind",
					},
					TemplateID: "some-template-id",
					Message: services.DispatchMessage{
						ReplyTo: "reply-to@example.com",
						Subject: "this is the subject",
						To:      "dr@strangelove.com",
						Text:    "email text",
						HTML: services.HTML{
							BodyContent:    "some html body content",
							BodyAttributes: "some html body attributes",
							Head:           "the html head tag",
							Doctype:        "the html doctype",
						},
					},
					VCAPRequest: services.DispatchVCAPRequest{
						ID:          "some-vcap-request-id",
						ReceiptTime: requestReceived,
					},
					UAAHost: "uaahost",
				})

				users := []services.User{{Email: "dr@strangelove.com"}}

				Expect(v1Enqueuer.EnqueueCall.Receives.Connection).To(Equal(conn))
				Expect(v1Enqueuer.EnqueueCall.Receives.Users).To(Equal(users))
				Expect(v1Enqueuer.EnqueueCall.Receives.Options).To(Equal(services.Options{
					ReplyTo:           "reply-to@example.com",
					Subject:           "this is the subject",
					KindDescription:   "description of a kind",
					SourceDescription: "description of a client",
					Text:              "email text",
					TemplateID:        "some-template-id",
					HTML: services.HTML{
						BodyContent:    "some html body content",
						BodyAttributes: "some html body attributes",
						Head:           "the html head tag",
						Doctype:        "the html doctype",
					},
					KindID:      "some-kind-id",
					To:          "dr@strangelove.com",
					Role:        "",
					Endorsement: services.EmailEndorsement,
				}))
				Expect(v1Enqueuer.EnqueueCall.Receives.Space).To(Equal(cf.CloudControllerSpace{}))
				Expect(v1Enqueuer.EnqueueCall.Receives.Org).To(Equal(cf.CloudControllerOrganization{}))
				Expect(v1Enqueuer.EnqueueCall.Receives.Client).To(Equal("some-client-id"))
				Expect(v1Enqueuer.EnqueueCall.Receives.Scope).To(Equal(""))
				Expect(v1Enqueuer.EnqueueCall.Receives.VCAPRequestID).To(Equal("some-vcap-request-id"))
				Expect(v1Enqueuer.EnqueueCall.Receives.RequestReceived).To(Equal(requestReceived))
				Expect(v1Enqueuer.EnqueueCall.Receives.UAAHost).To(Equal("uaahost"))

				Expect(v2Enqueuer.EnqueueCall.WasCalled).To(BeFalse())
			})
		})

		Context("when the dispatch JobType is v2", func() {
			It("calls Enqueue on its enqueuer with proper arguments", func() {
				emailStrategy.Dispatch(services.Dispatch{
					JobType:    "v2",
					Connection: conn,
					Client: services.DispatchClient{
						ID:          "some-client-id",
						Description: "description of a client",
					},
					TemplateID: "some-template-id",
					Message: services.DispatchMessage{
						ReplyTo: "reply-to@example.com",
						Subject: "this is the subject",
						To:      "dr@strangelove.com",
						Text:    "email text",
						HTML: services.HTML{
							BodyContent:    "some html body content",
							BodyAttributes: "some html body attributes",
							Head:           "the html head tag",
							Doctype:        "the html doctype",
						},
					},
					VCAPRequest: services.DispatchVCAPRequest{
						ID:          "some-vcap-request-id",
						ReceiptTime: requestReceived,
					},
					UAAHost: "uaahost",
				})

				Expect(v2Enqueuer.EnqueueCall.Receives.Connection).To(Equal(conn))
				Expect(v2Enqueuer.EnqueueCall.Receives.Users[0].Email).To(Equal("dr@strangelove.com"))
				Expect(v2Enqueuer.EnqueueCall.Receives.UAAHost).To(Equal("uaahost"))

				Expect(v1Enqueuer.EnqueueCall.WasCalled).To(BeFalse())
			})
		})
	})
})
