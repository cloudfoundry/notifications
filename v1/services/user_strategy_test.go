package services_test

import (
	"reflect"
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/v2/queue"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserStrategy", func() {
	var (
		strategy        services.UserStrategy
		v1Enqueuer      *mocks.Enqueuer
		v2Enqueuer      *mocks.V2Enqueuer
		conn            *mocks.Connection
		requestReceived time.Time
	)

	BeforeEach(func() {
		requestReceived, _ = time.Parse(time.RFC3339Nano, "2015-06-08T14:37:35.181067085-07:00")
		conn = mocks.NewConnection()
		v1Enqueuer = mocks.NewEnqueuer()
		v2Enqueuer = mocks.NewV2Enqueuer()
		strategy = services.NewUserStrategy(v1Enqueuer, v2Enqueuer)
	})

	Describe("Dispatch", func() {
		Context("when the job is not v2", func() {
			It("calls v1Enqueuer.Enqueue with the correct arguments for a user", func() {
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
					TemplateID: "some-template-id",
					UAAHost:    "uaa",
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

				Expect(reflect.ValueOf(v1Enqueuer.EnqueueCall.Receives.Connection).Pointer()).To(Equal(reflect.ValueOf(conn).Pointer()))
				Expect(v1Enqueuer.EnqueueCall.Receives.Users).To(Equal(users))
				Expect(v1Enqueuer.EnqueueCall.Receives.Options).To(Equal(services.Options{
					ReplyTo:           "reply-to@example.com",
					Subject:           "this is the subject",
					To:                "dr@strangelove.com",
					KindID:            "forgot_waterbottle",
					KindDescription:   "Water Bottle Reminder",
					SourceDescription: "The Water Bottle System",
					Text:              "Please make sure to leave your bottle in a place that is safe and dry",
					TemplateID:        "some-template-id",
					HTML: services.HTML{
						BodyContent:    "<p>The water bottle needs to be safe and dry</p>",
						BodyAttributes: "some-html-body-attributes",
						Head:           "<head></head>",
						Doctype:        "<html>",
					},
					Endorsement: services.UserEndorsement,
				}))
				Expect(v1Enqueuer.EnqueueCall.Receives.Space).To(Equal(cf.CloudControllerSpace{}))
				Expect(v1Enqueuer.EnqueueCall.Receives.Org).To(Equal(cf.CloudControllerOrganization{}))
				Expect(v1Enqueuer.EnqueueCall.Receives.Client).To(Equal("mister-client"))
				Expect(v1Enqueuer.EnqueueCall.Receives.Scope).To(Equal(""))
				Expect(v1Enqueuer.EnqueueCall.Receives.UAAHost).To(Equal("uaa"))
				Expect(v1Enqueuer.EnqueueCall.Receives.VCAPRequestID).To(Equal("some-vcap-request-id"))
				Expect(v1Enqueuer.EnqueueCall.Receives.RequestReceived).To(Equal(requestReceived))
			})
		})
		Context("when the job is v2", func() {
			It("calls v2Enqueuer.Enqueue with the correct arguments for a user", func() {
				_, err := strategy.Dispatch(services.Dispatch{
					JobType:    "v2",
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
					TemplateID: "some-template-id",
					UAAHost:    "uaa",
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

				users := []queue.User{{GUID: "user-123"}}
				Expect(reflect.ValueOf(v2Enqueuer.EnqueueCall.Receives.Connection).Pointer()).To(Equal(reflect.ValueOf(conn).Pointer()))
				Expect(v2Enqueuer.EnqueueCall.Receives.Users).To(Equal(users))
				Expect(v2Enqueuer.EnqueueCall.Receives.Options).To(Equal(queue.Options{
					ReplyTo:           "reply-to@example.com",
					Subject:           "this is the subject",
					To:                "dr@strangelove.com",
					KindID:            "forgot_waterbottle",
					KindDescription:   "Water Bottle Reminder",
					SourceDescription: "The Water Bottle System",
					Text:              "Please make sure to leave your bottle in a place that is safe and dry",
					TemplateID:        "some-template-id",
					HTML: queue.HTML{
						BodyContent:    "<p>The water bottle needs to be safe and dry</p>",
						BodyAttributes: "some-html-body-attributes",
						Head:           "<head></head>",
						Doctype:        "<html>",
					},
					Endorsement: services.UserEndorsement,
				}))
				Expect(v2Enqueuer.EnqueueCall.Receives.UAAHost).To(Equal("uaa"))
			})
		})
	})
})
