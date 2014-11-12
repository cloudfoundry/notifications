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

var _ = Describe("Organization Strategy", func() {
	var strategy strategies.OrganizationStrategy
	var options postal.Options
	var tokenLoader *fakes.TokenLoader
	var userLoader *fakes.UserLoader
	var organizationLoader *fakes.OrganizationLoader
	var templatesLoader *fakes.TemplatesLoader
	var mailer *fakes.Mailer
	var clientID string
	var receiptsRepo *fakes.ReceiptsRepo
	var conn *fakes.DBConn
	var findsUserGUIDs *fakes.FindsUserGUIDs
	var users map[string]uaa.User

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
		findsUserGUIDs.OrganizationGuids["org-001"] = []string{"user-123", "user-456"}

		users = map[string]uaa.User{
			"user-123": uaa.User{
				ID:     "user-123",
				Emails: []string{"user-123@example.com"},
			},
			"user-456": uaa.User{
				ID:     "user-456",
				Emails: []string{"user-456@example.com"},
			},
		}
		userLoader = fakes.NewUserLoader()
		userLoader.Users = users

		templatesLoader = fakes.NewTemplatesLoader()

		organizationLoader = fakes.NewOrganizationLoader()
		organizationLoader.Organization = cf.CloudControllerOrganization{
			Name: "my-org",
			GUID: "org-001",
		}

		strategy = strategies.NewOrganizationStrategy(tokenLoader, userLoader, organizationLoader, findsUserGUIDs, templatesLoader, mailer, receiptsRepo)
	})

	Describe("Dispatch", func() {
		Context("when the request is valid", func() {
			BeforeEach(func() {
				options = postal.Options{
					KindID:            "forgot_password",
					KindDescription:   "Password reminder",
					SourceDescription: "Login system",
					Text:              "Please reset your password by clicking on this link...",
					HTML:              postal.HTML{BodyContent: "<p>Please reset your password by clicking on this link...</p>"},
				}
			})

			It("records a receipt for each user", func() {
				_, err := strategy.Dispatch(clientID, "org-001", options, conn)
				if err != nil {
					panic(err)
				}

				Expect(receiptsRepo.CreateUserGUIDs).To(Equal([]string{"user-123", "user-456"}))
				Expect(receiptsRepo.ClientID).To(Equal(clientID))
				Expect(receiptsRepo.KindID).To(Equal(options.KindID))
			})

			It("call mailer.Deliver with the correct arguments for an organization", func() {
				templates := postal.Templates{
					Subject: "default-missing-subject",
					Text:    "default-organization-text",
					HTML:    "default-organization-html",
				}

				templatesLoader.Templates = templates

				_, err := strategy.Dispatch(clientID, "org-001", options, conn)
				if err != nil {
					panic(err)
				}

				Expect(templatesLoader.ContentSuffix).To(Equal(models.OrganizationBodyTemplateName))
				Expect(mailer.DeliverArguments).To(ContainElement(conn))
				Expect(mailer.DeliverArguments).To(ContainElement(templates))
				Expect(mailer.DeliverArguments).To(ContainElement(users))
				Expect(mailer.DeliverArguments).To(ContainElement(options))
				Expect(mailer.DeliverArguments).To(ContainElement(cf.CloudControllerOrganization{
					Name: "my-org",
					GUID: "org-001",
				}))
				Expect(mailer.DeliverArguments).To(ContainElement(cf.CloudControllerSpace{}))
				Expect(mailer.DeliverArguments).To(ContainElement(clientID))
				Expect(userLoader.LoadedGUIDs).To(Equal([]string{"user-123", "user-456"}))
			})
		})

		Context("failure cases", func() {
			Context("when token loader fails to return a token", func() {
				It("returns an error", func() {
					tokenLoader.LoadError = errors.New("BOOM!")
					_, err := strategy.Dispatch(clientID, "org-001", options, conn)

					Expect(err).To(Equal(errors.New("BOOM!")))
				})
			})

			Context("when organizationLoader fails to load an organization", func() {
				It("returns the error", func() {
					organizationLoader.LoadError = errors.New("BOOM!")
					_, err := strategy.Dispatch(clientID, "org-009", options, conn)

					Expect(err).To(Equal(errors.New("BOOM!")))
				})
			})

			Context("when userLoader fails to load users", func() {
				It("returns the error", func() {
					userLoader.LoadError = errors.New("BOOM!")
					_, err := strategy.Dispatch(clientID, "org-001", options, conn)

					Expect(err).To(Equal(errors.New("BOOM!")))
				})
			})

			Context("when templateLoader fails to load templates", func() {
				It("returns the error", func() {
					templatesLoader.LoadError = errors.New("BOOM!")

					_, err := strategy.Dispatch(clientID, "org-001", options, conn)

					Expect(err).To(BeAssignableToTypeOf(postal.TemplateLoadError("")))
				})
			})

			Context("when create receipts call returns an err", func() {
				It("returns an error", func() {
					receiptsRepo.CreateReceiptsError = true

					_, err := strategy.Dispatch(clientID, "org-001", options, conn)
					Expect(err).ToNot(BeNil())
				})
			})

			Context("when finds user GUIDs returns an error", func() {
				It("returns an error", func() {
					findsUserGUIDs.UserGUIDsBelongingToOrganizationError = errors.New("BOOM!")

					_, err := strategy.Dispatch(clientID, "org-001", options, conn)
					Expect(err).ToNot(BeNil())
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
						Recipient:      "user-123",
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
					"recipient":       "user-123",
					"notification_id": "123-456",
				}))
			})
		})
	})
})
