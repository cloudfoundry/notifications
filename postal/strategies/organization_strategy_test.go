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

var _ = Describe("Organization Strategy", func() {
	var (
		strategy           strategies.OrganizationStrategy
		options            postal.Options
		tokenLoader        *fakes.TokenLoader
		organizationLoader *fakes.OrganizationLoader
		mailer             *fakes.Mailer
		clientID           string
		conn               *fakes.DBConn
		findsUserGUIDs     *fakes.FindsUserGUIDs
		vcapRequestID      string
	)

	BeforeEach(func() {
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
		findsUserGUIDs.OrganizationGuids["org-001"] = []string{"user-123", "user-456"}

		organizationLoader = fakes.NewOrganizationLoader()
		organizationLoader.Organization = cf.CloudControllerOrganization{
			Name: "my-org",
			GUID: "org-001",
		}

		strategy = strategies.NewOrganizationStrategy(tokenLoader, organizationLoader, findsUserGUIDs, mailer)
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

			It("call mailer.Deliver with the correct arguments for an organization", func() {
				Expect(options.Endorsement).To(BeEmpty())

				_, err := strategy.Dispatch(clientID, "org-001", vcapRequestID, options, conn)
				if err != nil {
					panic(err)
				}

				options.Endorsement = strategies.OrganizationEndorsement
				users := []strategies.User{{GUID: "user-123"}, {GUID: "user-456"}}

				Expect(mailer.DeliverArguments).To(Equal(map[string]interface{}{
					"connection": conn,
					"users":      users,
					"options":    options,
					"space":      cf.CloudControllerSpace{},
					"org": cf.CloudControllerOrganization{
						Name: "my-org",
						GUID: "org-001",
					},
					"client":          clientID,
					"scope":           "",
					"vcap-request-id": vcapRequestID,
				}))
			})

			Context("when the org role field is set", func() {
				It("calls mailer.Deliver with the correct arguments", func() {
					options.Role = "OrgManager"

					Expect(options.Endorsement).To(BeEmpty())

					_, err := strategy.Dispatch(clientID, "org-001", vcapRequestID, options, conn)
					if err != nil {
						panic(err)
					}

					options.Endorsement = strategies.OrganizationRoleEndorsement

					Expect(mailer.DeliverArguments).To(ContainElement(options))
				})
			})
		})

		Context("failure cases", func() {
			Context("when token loader fails to return a token", func() {
				It("returns an error", func() {
					tokenLoader.LoadError = errors.New("BOOM!")
					_, err := strategy.Dispatch(clientID, "org-001", vcapRequestID, options, conn)

					Expect(err).To(Equal(errors.New("BOOM!")))
				})
			})

			Context("when organizationLoader fails to load an organization", func() {
				It("returns the error", func() {
					organizationLoader.LoadError = errors.New("BOOM!")
					_, err := strategy.Dispatch(clientID, "org-009", vcapRequestID, options, conn)

					Expect(err).To(Equal(errors.New("BOOM!")))
				})
			})

			Context("when finds user GUIDs returns an error", func() {
				It("returns an error", func() {
					findsUserGUIDs.UserGUIDsBelongingToOrganizationError = errors.New("BOOM!")

					_, err := strategy.Dispatch(clientID, "org-001", vcapRequestID, options, conn)
					Expect(err).ToNot(BeNil())
				})
			})
		})
	})
})
