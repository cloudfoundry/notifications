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

var _ = Describe("Everyone Strategy", func() {
	var strategy strategies.EveryoneStrategy
	var options postal.Options
	var tokenLoader *fakes.TokenLoader
	var templatesLoader *fakes.TemplatesLoader
	var allUsers *fakes.AllUsers
	var mailer *fakes.Mailer
	var clientID string
	var receiptsRepo *fakes.ReceiptsRepo
	var conn *fakes.DBConn
	var users map[string]uaa.User

	BeforeEach(func() {
		clientID = "my-client"
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

		templatesLoader = fakes.NewTemplatesLoader()

		users = map[string]uaa.User{
			"user-380": uaa.User{
				ID:     "user-380",
				Emails: []string{"user-380@example.com"},
			},
			"user-319": uaa.User{
				ID:     "user-319",
				Emails: []string{"user-319@example.com"},
			},
		}

		allUsers = fakes.NewAllUsers()
		allUsers.GUIDS = []string{"user-380", "user-319"}

		allUsers.Users = users

		strategy = strategies.NewEveryoneStrategy(allUsers, templatesLoader, mailer, receiptsRepo)
	})

	Describe("Dispatch", func() {
		BeforeEach(func() {
			options = postal.Options{
				KindID:            "welcome_user",
				KindDescription:   "Your Official Welcome",
				SourceDescription: "Welcome system",
				Text:              "Welcome to the system, now get off my lawn.",
				HTML:              postal.HTML{BodyContent: "<p>Welcome to the system, now get off my lawn.</p>"},
			}
		})

		It("records a receipt for each user", func() {
			_, err := strategy.Dispatch(clientID, "", options, conn)
			if err != nil {
				panic(err)
			}

			Expect(receiptsRepo.CreateUserGUIDs).To(Equal([]string{"user-380", "user-319"}))
			Expect(receiptsRepo.ClientID).To(Equal(clientID))
			Expect(receiptsRepo.KindID).To(Equal(options.KindID))
		})

		It("call mailer.Deliver with the correct arguments for an organization", func() {
			templates := postal.Templates{
				Subject: "default-missing-subject",
				Text:    "default-everyone-text",
				HTML:    "default-everyone-html",
			}

			templatesLoader.Templates = templates

			_, err := strategy.Dispatch(clientID, "", options, conn)
			if err != nil {
				panic(err)
			}

			Expect(templatesLoader.ContentSuffix).To(Equal(models.EveryoneBodyTemplateName))
			Expect(templatesLoader.SubjectSuffix).To(Equal(models.SubjectMissingTemplateName))
			Expect(mailer.DeliverArguments).To(ContainElement(conn))
			Expect(mailer.DeliverArguments).To(ContainElement(templates))
			Expect(mailer.DeliverArguments).To(ContainElement(users))
			Expect(mailer.DeliverArguments).To(ContainElement(options))
			Expect(mailer.DeliverArguments).To(ContainElement(cf.CloudControllerOrganization{}))
			Expect(mailer.DeliverArguments).To(ContainElement(cf.CloudControllerSpace{}))
			Expect(mailer.DeliverArguments).To(ContainElement(clientID))
		})
	})

	Context("failure cases", func() {
		Context("when allUsers fails to load users", func() {
			It("returns the error", func() {
				allUsers.LoadError = errors.New("BOOM!")
				_, err := strategy.Dispatch(clientID, "", options, conn)

				Expect(err).To(Equal(errors.New("BOOM!")))
			})
		})

		Context("when templateLoader fails to load templates", func() {
			It("returns the error", func() {
				templatesLoader.LoadError = errors.New("BOOM!")

				_, err := strategy.Dispatch(clientID, "", options, conn)

				Expect(err).To(BeAssignableToTypeOf(postal.TemplateLoadError("")))
			})
		})

		Context("when create receipts call returns an err", func() {
			It("returns an error", func() {
				receiptsRepo.CreateReceiptsError = true

				_, err := strategy.Dispatch(clientID, "", options, conn)
				Expect(err).ToNot(BeNil())
			})
		})
	})

	Describe("Trim", func() {
		Describe("TrimFields", func() {
			It("trims the specified fields from the response object", func() {
				responses, err := json.Marshal([]strategies.Response{
					{
						Status:         "delivered",
						Recipient:      "user-319",
						Email:          "",
						NotificationID: "123-456",
					},
					{
						Status:         "delivered",
						Recipient:      "user-380",
						Email:          "",
						NotificationID: "789-1011",
					},
				})

				trimmedResponses := strategy.Trim(responses)

				var result []map[string]string
				err = json.Unmarshal(trimmedResponses, &result)
				if err != nil {
					panic(err)
				}

				Expect(result).To(Equal([]map[string]string{
					{"status": "delivered",
						"recipient":       "user-319",
						"notification_id": "123-456",
					},
					{"status": "delivered",
						"recipient":       "user-380",
						"notification_id": "789-1011",
					},
				}))
			})
		})
	})
})
