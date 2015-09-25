package services_test

import (
	"errors"
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/testing/helpers"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/v2/queue"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Space Strategy", func() {
	var (
		strategy           services.SpaceStrategy
		tokenLoader        *mocks.TokenLoader
		spaceLoader        *mocks.SpaceLoader
		organizationLoader *mocks.OrganizationLoader
		v1Enqueuer         *mocks.Enqueuer
		v2Enqueuer         *mocks.V2Enqueuer
		conn               *mocks.Connection
		findsUserIDs       *mocks.FindsUserIDs
		requestReceived    time.Time
		token              string
	)

	BeforeEach(func() {
		requestReceived, _ = time.Parse(time.RFC3339Nano, "2015-06-08T14:37:35.181067085-07:00")
		conn = mocks.NewConnection()
		tokenHeader := map[string]interface{}{
			"alg": "FAST",
		}
		tokenClaims := map[string]interface{}{
			"client_id": "mister-client",
			"exp":       int64(3404281214),
			"iss":       "uaa",
			"scope":     []string{"notifications.write"},
		}
		token = helpers.BuildToken(tokenHeader, tokenClaims)

		tokenLoader = mocks.NewTokenLoader()
		tokenLoader.LoadCall.Returns.Token = token
		v1Enqueuer = mocks.NewEnqueuer()
		v2Enqueuer = mocks.NewV2Enqueuer()

		findsUserIDs = mocks.NewFindsUserIDs()
		findsUserIDs.UserIDsBelongingToSpaceCall.Returns.UserIDs = []string{"user-123", "user-456"}

		spaceLoader = mocks.NewSpaceLoader()
		spaceLoader.LoadCall.Returns.Space = cf.CloudControllerSpace{
			Name:             "production",
			GUID:             "space-001",
			OrganizationGUID: "org-001",
		}
		organizationLoader = mocks.NewOrganizationLoader()
		organizationLoader.LoadCall.Returns.Organization = cf.CloudControllerOrganization{
			Name: "the-org",
			GUID: "org-001",
		}
		strategy = services.NewSpaceStrategy(tokenLoader, spaceLoader, organizationLoader, findsUserIDs, v1Enqueuer, v2Enqueuer)
	})

	Describe("Dispatch", func() {
		Context("when the request is valid", func() {
			Context("and the dispatch JobType is v1", func() {
				It("calls v1Enqueuer.Enqueue with the correct arguments for a space", func() {
					_, err := strategy.Dispatch(services.Dispatch{
						GUID:       "space-001",
						Connection: conn,
						Message: services.DispatchMessage{
							To:      "dr@strangelove.com",
							ReplyTo: "reply-to@example.com",
							Subject: "this is the subject",
							Text:    "Please reset your password by clicking on this link...",
							HTML: services.HTML{
								BodyContent:    "<p>Welcome to the system, now get off my lawn.</p>",
								BodyAttributes: "some-html-body-attributes",
								Head:           "<head></head>",
								Doctype:        "<html>",
							},
						},
						TemplateID: "some-template-id",
						Kind: services.DispatchKind{
							ID:          "forgot_password",
							Description: "Password reminder",
						},
						Client: services.DispatchClient{
							ID:          "mister-client",
							Description: "Login system",
						},
						VCAPRequest: services.DispatchVCAPRequest{
							ID:          "some-vcap-request-id",
							ReceiptTime: requestReceived,
						},
						UAAHost: "uaa",
					})
					Expect(err).NotTo(HaveOccurred())

					users := []services.User{{GUID: "user-123"}, {GUID: "user-456"}}

					Expect(organizationLoader.LoadCall.Receives.OrganizationGUID).To(Equal("org-001"))
					Expect(organizationLoader.LoadCall.Receives.Token).To(Equal(tokenLoader.LoadCall.Returns.Token))

					Expect(spaceLoader.LoadCall.Receives.SpaceGUID).To(Equal("space-001"))
					Expect(spaceLoader.LoadCall.Receives.Token).To(Equal(tokenLoader.LoadCall.Returns.Token))

					Expect(v1Enqueuer.EnqueueCall.Receives.Connection).To(Equal(conn))
					Expect(v1Enqueuer.EnqueueCall.Receives.Users).To(Equal(users))
					Expect(v1Enqueuer.EnqueueCall.Receives.Options).To(Equal(services.Options{
						ReplyTo:           "reply-to@example.com",
						Subject:           "this is the subject",
						To:                "dr@strangelove.com",
						KindID:            "forgot_password",
						KindDescription:   "Password reminder",
						SourceDescription: "Login system",
						Text:              "Please reset your password by clicking on this link...",
						TemplateID:        "some-template-id",
						HTML: services.HTML{
							BodyContent:    "<p>Welcome to the system, now get off my lawn.</p>",
							BodyAttributes: "some-html-body-attributes",
							Head:           "<head></head>",
							Doctype:        "<html>",
						},
						Endorsement: services.SpaceEndorsement,
					}))
					Expect(v1Enqueuer.EnqueueCall.Receives.Space).To(Equal(cf.CloudControllerSpace{
						GUID:             "space-001",
						Name:             "production",
						OrganizationGUID: "org-001",
					}))
					Expect(v1Enqueuer.EnqueueCall.Receives.Org).To(Equal(cf.CloudControllerOrganization{
						Name: "the-org",
						GUID: "org-001",
					}))
					Expect(v1Enqueuer.EnqueueCall.Receives.Client).To(Equal("mister-client"))
					Expect(v1Enqueuer.EnqueueCall.Receives.Scope).To(Equal(""))
					Expect(v1Enqueuer.EnqueueCall.Receives.VCAPRequestID).To(Equal("some-vcap-request-id"))
					Expect(v1Enqueuer.EnqueueCall.Receives.RequestReceived).To(Equal(requestReceived))
					Expect(v1Enqueuer.EnqueueCall.Receives.UAAHost).To(Equal("uaa"))

					Expect(tokenLoader.LoadCall.Receives.UAAHost).To(Equal("uaa"))

					Expect(findsUserIDs.UserIDsBelongingToSpaceCall.Receives.SpaceGUID).To(Equal("space-001"))
					Expect(findsUserIDs.UserIDsBelongingToSpaceCall.Receives.Token).To(Equal(token))
				})
			})
			Context("and the dispatch JobType is v2", func() {
				It("calls v2Enqueuer.Enqueue with the correct arguments for a space", func() {
					_, err := strategy.Dispatch(services.Dispatch{
						JobType:    "v2",
						GUID:       "space-001",
						Connection: conn,
						Message: services.DispatchMessage{
							To:      "dr@strangelove.com",
							ReplyTo: "reply-to@example.com",
							Subject: "this is the subject",
							Text:    "Please reset your password by clicking on this link...",
							HTML: services.HTML{
								BodyContent:    "<p>Welcome to the system, now get off my lawn.</p>",
								BodyAttributes: "some-html-body-attributes",
								Head:           "<head></head>",
								Doctype:        "<html>",
							},
						},
						TemplateID: "some-template-id",
						CampaignID: "some-campaign-id",
						Kind: services.DispatchKind{
							ID:          "forgot_password",
							Description: "Password reminder",
						},
						Client: services.DispatchClient{
							ID:          "mister-client",
							Description: "Login system",
						},
						VCAPRequest: services.DispatchVCAPRequest{
							ID:          "some-vcap-request-id",
							ReceiptTime: requestReceived,
						},
						UAAHost: "uaa",
					})
					Expect(err).NotTo(HaveOccurred())

					users := []queue.User{{GUID: "user-123"}, {GUID: "user-456"}}
					Expect(v2Enqueuer.EnqueueCall.Receives.Connection).To(Equal(conn))
					Expect(v2Enqueuer.EnqueueCall.Receives.Users).To(Equal(users))
					Expect(v2Enqueuer.EnqueueCall.Receives.Options).To(Equal(queue.Options{
						ReplyTo:           "reply-to@example.com",
						Subject:           "this is the subject",
						To:                "dr@strangelove.com",
						KindID:            "forgot_password",
						KindDescription:   "Password reminder",
						SourceDescription: "Login system",
						Text:              "Please reset your password by clicking on this link...",
						TemplateID:        "some-template-id",
						HTML: queue.HTML{
							BodyContent:    "<p>Welcome to the system, now get off my lawn.</p>",
							BodyAttributes: "some-html-body-attributes",
							Head:           "<head></head>",
							Doctype:        "<html>",
						},
						Endorsement: services.SpaceEndorsement,
					}))
					Expect(v2Enqueuer.EnqueueCall.Receives.UAAHost).To(Equal("uaa"))
					Expect(v2Enqueuer.EnqueueCall.Receives.CampaignID).To(Equal("some-campaign-id"))
				})
			})
		})

		Context("failure cases", func() {
			Context("when token loader fails to return a token", func() {
				It("returns an error", func() {
					tokenLoader.LoadCall.Returns.Error = errors.New("BOOM!")

					_, err := strategy.Dispatch(services.Dispatch{})
					Expect(err).To(Equal(errors.New("BOOM!")))
				})
			})

			Context("when spaceLoader fails to load a space", func() {
				It("returns an error", func() {
					spaceLoader.LoadCall.Returns.Error = errors.New("BOOM!")

					_, err := strategy.Dispatch(services.Dispatch{})
					Expect(err).To(Equal(errors.New("BOOM!")))
				})
			})

			Context("when findsUserIDs returns an err", func() {
				It("returns an error", func() {
					findsUserIDs.UserIDsBelongingToSpaceCall.Returns.Error = errors.New("BOOM!")

					_, err := strategy.Dispatch(services.Dispatch{})
					Expect(err).To(Equal(errors.New("BOOM!")))
				})
			})
		})
	})
})
