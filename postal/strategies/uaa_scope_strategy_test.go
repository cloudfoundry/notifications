package strategies_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UAA Scope Strategy", func() {
	var (
		strategy       strategies.UAAScopeStrategy
		options        postal.Options
		tokenLoader    *fakes.TokenLoader
		mailer         *fakes.Mailer
		clientID       string
		conn           *fakes.DBConn
		findsUserGUIDs *fakes.FindsUserGUIDs
		scope          string
		vcapRequestID  string
	)

	BeforeEach(func() {
		scope = "great.scope"
		clientID = "mister-client"
		vcapRequestID = "some-request-id"
		conn = fakes.NewDBConn()

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
		findsUserGUIDs.GUIDsWithScopes[scope] = []string{"user-311"}

		strategy = strategies.NewUAAScopeStrategy(tokenLoader, findsUserGUIDs, mailer)
	})

	Describe("Dispatch", func() {
		Context("when the request is valid", func() {
			BeforeEach(func() {
				options = postal.Options{
					KindID:            "forgot_waterbottle",
					KindDescription:   "Water Bottle Reminder",
					SourceDescription: "The Water Bottle System",
					Text:              "Please make sure to leave your bottle in a place that is safe and dry",
					HTML:              postal.HTML{BodyContent: "<p>The water bottle needs to be safe and dry</p>"},
				}
			})

			It("call mailer.Deliver with the correct arguments for an UAA Scope", func() {
				Expect(options.Endorsement).To(BeEmpty())

				_, err := strategy.Dispatch(clientID, "great.scope", vcapRequestID, options, conn)
				if err != nil {
					panic(err)
				}

				options.Endorsement = strategies.ScopeEndorsement
				users := []strategies.User{{GUID: "user-311"}}

				Expect(mailer.DeliverArguments).To(Equal(map[string]interface{}{
					"connection":      conn,
					"users":           users,
					"options":         options,
					"space":           cf.CloudControllerSpace{},
					"org":             cf.CloudControllerOrganization{},
					"client":          clientID,
					"scope":           scope,
					"vcap-request-id": vcapRequestID,
				}))
			})
		})

		Context("failure cases", func() {
			Context("when token loader fails to return a token", func() {
				It("returns an error", func() {
					tokenLoader.LoadError = errors.New("BOOM!")
					_, err := strategy.Dispatch(clientID, "great.scope", vcapRequestID, options, conn)

					Expect(err).To(Equal(errors.New("BOOM!")))
				})
			})

			Context("when finds user GUIDs returns an error", func() {
				It("returns an error", func() {
					findsUserGUIDs.UserGUIDsBelongingToScopeError = errors.New("BOOM!")

					_, err := strategy.Dispatch(clientID, "great.scope", vcapRequestID, options, conn)
					Expect(err).To(HaveOccurred())
				})
			})

			Context("when an invalid scope is passed", func() {
				It("returns an error", func() {
					defaultScopes := []string{"cloud_controller.read", "cloud_controller.write", "openid", "approvals.me",
						"cloud_controller_service_permissions.read", "scim.me", "uaa.user", "password.write", "scim.userids", "oauth.approvals"}
					for _, scope := range defaultScopes {
						_, err := strategy.Dispatch(clientID, scope, vcapRequestID, options, conn)
						Expect(err).To(HaveOccurred())
						Expect(err).To(MatchError(strategies.DefaultScopeError{}))
					}
				})
			})
		})
	})
})
