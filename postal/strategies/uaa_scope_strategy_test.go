package strategies_test

import (
	"encoding/json"
	"errors"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UAA Scope Strategy", func() {
	var strategy strategies.UAAScopeStrategy
	var options postal.Options
	var tokenLoader *fakes.TokenLoader
	var userLoader *fakes.UserLoader
	var templatesLoader *fakes.TemplatesLoader
	var mailer *fakes.Mailer
	var clientID string
	var receiptsRepo *fakes.ReceiptsRepo
	var conn *fakes.DBConn
	var findsUserGUIDs *fakes.FindsUserGUIDs
	var users map[string]uaa.User

	scope := "great.scope"

	BeforeEach(func() {
		clientID = "mister-client"
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

		receiptsRepo = fakes.NewReceiptsRepo()

		mailer = fakes.NewMailer()

		findsUserGUIDs = fakes.NewFindsUserGUIDs()
		findsUserGUIDs.GUIDsWithScopes[scope] = []string{"user-311"}

		users = map[string]uaa.User{
			"user-311": uaa.User{
				ID:     "user-311",
				Emails: []string{"user-311@example.com"},
			},
		}
		userLoader = fakes.NewUserLoader()
		userLoader.Users = users

		templatesLoader = fakes.NewTemplatesLoader()

		strategy = strategies.NewUAAScopeStrategy(tokenLoader, userLoader, findsUserGUIDs, templatesLoader, mailer, receiptsRepo)
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

			It("records a receipt for each user", func() {
				_, err := strategy.Dispatch(clientID, scope, options, conn)
				if err != nil {
					panic(err)
				}

				Expect(receiptsRepo.CreateUserGUIDs).To(Equal([]string{"user-311"}))
				Expect(receiptsRepo.ClientID).To(Equal(clientID))
				Expect(receiptsRepo.KindID).To(Equal(options.KindID))
			})

			It("call mailer.Deliver with the correct arguments for an UAA Scope", func() {
				templates := postal.Templates{
					Subject: "default-missing-subject",
					Text:    "default-scope-text",
					HTML:    "default-scope-html",
				}

				templatesLoader.Templates = templates

				_, err := strategy.Dispatch(clientID, "great.scope", options, conn)
				if err != nil {
					panic(err)
				}

				Expect(templatesLoader.ContentSuffix).To(Equal(models.UAAScopeBodyTemplateName))
				Expect(mailer.DeliverArguments).To(ContainElement(conn))
				Expect(mailer.DeliverArguments).To(ContainElement(templates))
				Expect(mailer.DeliverArguments).To(ContainElement(users))
				Expect(mailer.DeliverArguments).To(ContainElement(options))
				Expect(mailer.DeliverArguments).To(ContainElement(cf.CloudControllerOrganization{}))
				Expect(mailer.DeliverArguments).To(ContainElement(cf.CloudControllerSpace{}))
				Expect(mailer.DeliverArguments).To(ContainElement(scope))
				Expect(mailer.DeliverArguments).To(ContainElement(clientID))
				Expect(userLoader.LoadedGUIDs).To(Equal([]string{"user-311"}))
			})
		})

		Context("failure cases", func() {
			Context("when token loader fails to return a token", func() {
				It("returns an error", func() {
					tokenLoader.LoadError = errors.New("BOOM!")
					_, err := strategy.Dispatch(clientID, "great.scope", options, conn)

					Expect(err).To(Equal(errors.New("BOOM!")))
				})
			})

			Context("when userLoader fails to load users", func() {
				It("returns the error", func() {
					userLoader.LoadError = errors.New("BOOM!")
					_, err := strategy.Dispatch(clientID, "great.scope", options, conn)

					Expect(err).To(Equal(errors.New("BOOM!")))
				})
			})

			Context("when templateLoader fails to load templates", func() {
				It("returns the error", func() {
					templatesLoader.LoadError = errors.New("BOOM!")

					_, err := strategy.Dispatch(clientID, "great.scope", options, conn)

					Expect(err).To(BeAssignableToTypeOf(postal.TemplateLoadError("")))
				})
			})

			Context("when create receipts call returns an err", func() {
				It("returns an error", func() {
					receiptsRepo.CreateReceiptsError = true

					_, err := strategy.Dispatch(clientID, "great.scope", options, conn)
					Expect(err).To(HaveOccurred())
				})
			})

			Context("when finds user GUIDs returns an error", func() {
				It("returns an error", func() {
					findsUserGUIDs.UserGUIDsBelongingToScopeError = errors.New("BOOM!")

					_, err := strategy.Dispatch(clientID, "great.scope", options, conn)
					Expect(err).To(HaveOccurred())
				})
			})

			Context("when an invalid scope is passed", func() {
				It("returns an error", func() {
					defaultScopes := []string{"cloud_controller.read", "cloud_controller.write", "openid", "approvals.me",
						"cloud_controller_service_permissions.read", "scim.me", "uaa.user", "password.write", "scim.userids", "oauth.approvals"}
					for _, scope := range defaultScopes {
						_, err := strategy.Dispatch(clientID, scope, options, conn)
						Expect(err).To(HaveOccurred())
						Expect(err).To(MatchError(strategies.DefaultScopeError{}))
					}
				})
			})
		})
	})

	Describe("Trim", func() {
		Describe("TrimFields", func() {
			It("trims the specified fields from the response object", func() {
				responses, err := json.Marshal([]strategies.Response{
					{
						Status:         "delivered",
						Recipient:      "user-311",
						NotificationID: "123-456",
					},
				})

				trimmedResponses := strategy.Trim(responses)

				var result []map[string]string
				err = json.Unmarshal(trimmedResponses, &result)
				if err != nil {
					panic(err)
				}

				Expect(result).To(ContainElement(map[string]string{"status": "delivered",
					"recipient":       "user-311",
					"notification_id": "123-456",
				}))
			})
		})
	})
})
