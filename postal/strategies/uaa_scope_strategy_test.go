package strategies_test

import (
	"errors"
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UAA Scope Strategy", func() {
	var (
		strategy        strategies.UAAScopeStrategy
		tokenLoader     *fakes.TokenLoader
		mailer          *fakes.Mailer
		conn            *fakes.Connection
		findsUserGUIDs  *fakes.FindsUserGUIDs
		requestReceived time.Time
	)

	BeforeEach(func() {
		requestReceived, _ = time.Parse(time.RFC3339Nano, "2015-06-08T14:37:35.181067085-07:00")
		conn = fakes.NewConnection()
		tokenHeader := map[string]interface{}{
			"alg": "FAST",
		}
		tokenClaims := map[string]interface{}{
			"client_id": "mister-client",
			"exp":       int64(3404281214),
			"scope":     []string{"notifications.write"},
		}
		tokenLoader = fakes.NewTokenLoader()
		tokenLoader.Token = fakes.BuildToken(tokenHeader, tokenClaims)
		mailer = fakes.NewMailer()
		findsUserGUIDs = fakes.NewFindsUserGUIDs()
		findsUserGUIDs.GUIDsWithScopes["great.scope"] = []string{"user-311"}
		strategy = strategies.NewUAAScopeStrategy(tokenLoader, findsUserGUIDs, mailer)
	})

	Describe("Dispatch", func() {
		Context("when the request is valid", func() {
			It("call mailer.Deliver with the correct arguments for an UAA Scope", func() {
				_, err := strategy.Dispatch(strategies.Dispatch{
					GUID:       "great.scope",
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

				users := []strategies.User{{GUID: "user-311"}}

				Expect(mailer.DeliverCall.Args.Connection).To(Equal(conn))
				Expect(mailer.DeliverCall.Args.Users).To(Equal(users))
				Expect(mailer.DeliverCall.Args.Options).To(Equal(postal.Options{
					ReplyTo:           "reply-to@example.com",
					Subject:           "this is the subject",
					To:                "dr@strangelove.com",
					KindID:            "forgot_waterbottle",
					KindDescription:   "Water Bottle Reminder",
					SourceDescription: "The Water Bottle System",
					Text:              "Please make sure to leave your bottle in a place that is safe and dry",
					HTML: postal.HTML{
						BodyContent:    "<p>The water bottle needs to be safe and dry</p>",
						BodyAttributes: "some-html-body-attributes",
						Head:           "<head></head>",
						Doctype:        "<html>",
					},
					Endorsement: strategies.ScopeEndorsement,
				}))
				Expect(mailer.DeliverCall.Args.Space).To(Equal(cf.CloudControllerSpace{}))
				Expect(mailer.DeliverCall.Args.Org).To(Equal(cf.CloudControllerOrganization{}))
				Expect(mailer.DeliverCall.Args.Client).To(Equal("mister-client"))
				Expect(mailer.DeliverCall.Args.Scope).To(Equal("great.scope"))
				Expect(mailer.DeliverCall.Args.VCAPRequestID).To(Equal("some-vcap-request-id"))
				Expect(mailer.DeliverCall.Args.RequestReceived).To(Equal(requestReceived))
			})
		})

		Context("failure cases", func() {
			Context("when token loader fails to return a token", func() {
				It("returns an error", func() {
					tokenLoader.LoadError = errors.New("BOOM!")

					_, err := strategy.Dispatch(strategies.Dispatch{})
					Expect(err).To(Equal(errors.New("BOOM!")))
				})
			})

			Context("when finds user GUIDs returns an error", func() {
				It("returns an error", func() {
					findsUserGUIDs.UserGUIDsBelongingToScopeError = errors.New("BOOM!")

					_, err := strategy.Dispatch(strategies.Dispatch{})
					Expect(err).To(HaveOccurred())
				})
			})

			Context("when an invalid scope is passed", func() {
				It("returns an error", func() {
					defaultScopes := []string{
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

					for _, scope := range defaultScopes {
						_, err := strategy.Dispatch(strategies.Dispatch{
							GUID: scope,
						})
						Expect(err).To(HaveOccurred())
						Expect(err).To(MatchError(strategies.DefaultScopeError{}))
					}
				})
			})
		})
	})
})
