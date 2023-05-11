package services_test

import (
	"errors"
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/testing/helpers"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UAA Scope Strategy", func() {
	var (
		strategy        services.UAAScopeStrategy
		tokenLoader     *mocks.TokenLoader
		enqueuer        *mocks.Enqueuer
		conn            *mocks.Connection
		findsUserIDs    *mocks.FindsUserIDs
		requestReceived time.Time
		defaultScopes   []string
		token           string
	)

	BeforeEach(func() {
		defaultScopes = []string{
			"cloud_controller.read",
			"cloud_controller.write",
			"openid",
			"approvals.me",
			"cloud_controller_service_permissions.read",
			"scim.me",
			"uaa.user",
			"password.write",
			"scim.userids",
			"oauth.approvals",
		}

		requestReceived, _ = time.Parse(time.RFC3339Nano, "2015-06-08T14:37:35.181067085-07:00")
		conn = mocks.NewConnection()

		tokenHeader := map[string]interface{}{
			"alg": "RS256",
		}
		tokenClaims := map[string]interface{}{
			"client_id": "mister-client",
			"exp":       int64(3404281214),
			"scope":     []string{"notifications.write"},
		}
		token = helpers.BuildToken(tokenHeader, tokenClaims)

		tokenLoader = mocks.NewTokenLoader()
		tokenLoader.LoadCall.Returns.Token = token
		enqueuer = mocks.NewEnqueuer()

		findsUserIDs = mocks.NewFindsUserIDs()
		findsUserIDs.UserIDsBelongingToScopeCall.Returns.UserIDs = []string{"user-311"}

		strategy = services.NewUAAScopeStrategy(tokenLoader, findsUserIDs, enqueuer, defaultScopes)
	})

	Describe("Dispatch", func() {
		Context("when the JobType is unspecified", func() {
			Context("when the request is valid", func() {
				It("should call enqueuer.Enqueue with the correct arguments for an UAA Scope", func() {
					_, err := strategy.Dispatch(services.Dispatch{
						GUID:       "great.scope",
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
						UAAHost: "uaa",
					})
					Expect(err).NotTo(HaveOccurred())

					users := []services.User{{GUID: "user-311"}}

					Expect(enqueuer.EnqueueCall.Receives.Connection).To(Equal(conn))
					Expect(enqueuer.EnqueueCall.Receives.Users).To(Equal(users))
					Expect(enqueuer.EnqueueCall.Receives.Options).To(Equal(services.Options{
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
						Endorsement: services.ScopeEndorsement,
					}))
					Expect(enqueuer.EnqueueCall.Receives.Space).To(Equal(cf.CloudControllerSpace{}))
					Expect(enqueuer.EnqueueCall.Receives.Org).To(Equal(cf.CloudControllerOrganization{}))
					Expect(enqueuer.EnqueueCall.Receives.Client).To(Equal("mister-client"))
					Expect(enqueuer.EnqueueCall.Receives.Scope).To(Equal("great.scope"))
					Expect(enqueuer.EnqueueCall.Receives.VCAPRequestID).To(Equal("some-vcap-request-id"))
					Expect(enqueuer.EnqueueCall.Receives.RequestReceived).To(Equal(requestReceived))
					Expect(enqueuer.EnqueueCall.Receives.UAAHost).To(Equal("uaa"))

					Expect(findsUserIDs.UserIDsBelongingToScopeCall.Receives.Scope).To(Equal("great.scope"))
					Expect(findsUserIDs.UserIDsBelongingToScopeCall.Receives.Token).To(Equal(token))
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

			Context("when finds user IDs returns an error", func() {
				It("returns an error", func() {
					findsUserIDs.UserIDsBelongingToScopeCall.Returns.Error = errors.New("BOOM!")

					_, err := strategy.Dispatch(services.Dispatch{})
					Expect(err).To(HaveOccurred())
				})
			})

			Context("when an default scope is passed", func() {
				It("returns an error", func() {
					for _, scope := range defaultScopes {
						_, err := strategy.Dispatch(services.Dispatch{
							GUID: scope,
						})
						Expect(err).To(HaveOccurred())
						Expect(err).To(MatchError(services.DefaultScopeError{}))
					}
				})
			})
		})
	})
})
