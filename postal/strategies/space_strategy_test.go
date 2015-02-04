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

var _ = Describe("Space Strategy", func() {
	var strategy strategies.SpaceStrategy
	var options postal.Options
	var tokenLoader *fakes.TokenLoader
	var spaceLoader *fakes.SpaceLoader
	var organizationLoader *fakes.OrganizationLoader
	var mailer *fakes.Mailer
	var clientID string
	var conn *fakes.DBConn
	var findsUserGUIDs *fakes.FindsUserGUIDs

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

		mailer = fakes.NewMailer()

		findsUserGUIDs = fakes.NewFindsUserGUIDs()
		findsUserGUIDs.SpaceGuids["space-001"] = []string{"user-123", "user-456"}

		spaceLoader = fakes.NewSpaceLoader()
		spaceLoader.Space = cf.CloudControllerSpace{
			Name:             "production",
			GUID:             "space-001",
			OrganizationGUID: "org-001",
		}

		organizationLoader = fakes.NewOrganizationLoader()
		organizationLoader.Organization = cf.CloudControllerOrganization{
			Name: "the-org",
			GUID: "org-001",
		}

		strategy = strategies.NewSpaceStrategy(tokenLoader, spaceLoader, organizationLoader, findsUserGUIDs, mailer)
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

			It("calls mailer.Deliver with the correct arguments for a space", func() {
				Expect(options.Endorsement).To(BeEmpty())

				_, err := strategy.Dispatch(clientID, "space-001", options, conn)
				if err != nil {
					panic(err)
				}

				users := []strategies.User{{GUID: "user-123"}, {GUID: "user-456"}}

				options.Endorsement = strategies.SpaceEndorsement
				Expect(mailer.DeliverArguments).To(Equal(map[string]interface{}{
					"connection": conn,
					"users":      users,
					"options":    options,
					"space": cf.CloudControllerSpace{
						GUID:             "space-001",
						Name:             "production",
						OrganizationGUID: "org-001",
					},
					"org": cf.CloudControllerOrganization{
						Name: "the-org",
						GUID: "org-001",
					},
					"client": clientID,
					"scope":  "",
				}))
			})
		})

		Context("failure cases", func() {
			Context("when token loader fails to return a token", func() {
				It("returns an error", func() {
					tokenLoader.LoadError = errors.New("BOOM!")
					_, err := strategy.Dispatch(clientID, "space-001", options, conn)

					Expect(err).To(Equal(errors.New("BOOM!")))
				})
			})

			Context("when spaceLoader fails to load a space", func() {
				It("returns an error", func() {
					spaceLoader.LoadError = errors.New("BOOM!")
					_, err := strategy.Dispatch(clientID, "space-000", options, conn)

					Expect(err).To(Equal(errors.New("BOOM!")))
				})
			})

			Context("when findsUserGUIDs returns an err", func() {
				It("returns an error", func() {
					findsUserGUIDs.UserGUIDsBelongingToSpaceError = errors.New("BOOM!")

					_, err := strategy.Dispatch(clientID, "space-001", options, conn)
					Expect(err).To(Equal(findsUserGUIDs.UserGUIDsBelongingToSpaceError))
				})
			})
		})
	})
})
