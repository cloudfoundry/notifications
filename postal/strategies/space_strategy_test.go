package strategies_test

import (
	"errors"
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Space Strategy", func() {
	var (
		strategy           strategies.SpaceStrategy
		tokenLoader        *fakes.TokenLoader
		spaceLoader        *fakes.SpaceLoader
		organizationLoader *fakes.OrganizationLoader
		enqueuer           *fakes.Enqueuer
		conn               *fakes.Connection
		findsUserGUIDs     *fakes.FindsUserGUIDs
		requestReceived    time.Time
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
		enqueuer = fakes.NewEnqueuer()
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
		strategy = strategies.NewSpaceStrategy(tokenLoader, spaceLoader, organizationLoader, findsUserGUIDs, enqueuer)
	})

	Describe("Dispatch", func() {
		Context("when the request is valid", func() {
			It("calls enqueuer.Enqueue with the correct arguments for a space", func() {
				_, err := strategy.Dispatch(strategies.Dispatch{
					GUID:       "space-001",
					Connection: conn,
					Message: strategies.Message{
						To:      "dr@strangelove.com",
						ReplyTo: "reply-to@example.com",
						Subject: "this is the subject",
						Text:    "Please reset your password by clicking on this link...",
						HTML: strategies.HTML{
							BodyContent:    "<p>Welcome to the system, now get off my lawn.</p>",
							BodyAttributes: "some-html-body-attributes",
							Head:           "<head></head>",
							Doctype:        "<html>",
						},
					},
					Kind: strategies.Kind{
						ID:          "forgot_password",
						Description: "Password reminder",
					},
					Client: strategies.Client{
						ID:          "mister-client",
						Description: "Login system",
					},
					VCAPRequest: strategies.VCAPRequest{
						ID:          "some-vcap-request-id",
						ReceiptTime: requestReceived,
					},
				})
				Expect(err).NotTo(HaveOccurred())

				users := []strategies.User{{GUID: "user-123"}, {GUID: "user-456"}}

				Expect(enqueuer.EnqueueCall.Args.Connection).To(Equal(conn))
				Expect(enqueuer.EnqueueCall.Args.Users).To(Equal(users))
				Expect(enqueuer.EnqueueCall.Args.Options).To(Equal(strategies.Options{
					ReplyTo:           "reply-to@example.com",
					Subject:           "this is the subject",
					To:                "dr@strangelove.com",
					KindID:            "forgot_password",
					KindDescription:   "Password reminder",
					SourceDescription: "Login system",
					Text:              "Please reset your password by clicking on this link...",
					HTML: strategies.HTML{
						BodyContent:    "<p>Welcome to the system, now get off my lawn.</p>",
						BodyAttributes: "some-html-body-attributes",
						Head:           "<head></head>",
						Doctype:        "<html>",
					},
					Endorsement: strategies.SpaceEndorsement,
				}))
				Expect(enqueuer.EnqueueCall.Args.Space).To(Equal(cf.CloudControllerSpace{
					GUID:             "space-001",
					Name:             "production",
					OrganizationGUID: "org-001",
				}))
				Expect(enqueuer.EnqueueCall.Args.Org).To(Equal(cf.CloudControllerOrganization{
					Name: "the-org",
					GUID: "org-001",
				}))
				Expect(enqueuer.EnqueueCall.Args.Client).To(Equal("mister-client"))
				Expect(enqueuer.EnqueueCall.Args.Scope).To(Equal(""))
				Expect(enqueuer.EnqueueCall.Args.VCAPRequestID).To(Equal("some-vcap-request-id"))
				Expect(enqueuer.EnqueueCall.Args.RequestReceived).To(Equal(requestReceived))
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

			Context("when spaceLoader fails to load a space", func() {
				It("returns an error", func() {
					spaceLoader.LoadError = errors.New("BOOM!")

					_, err := strategy.Dispatch(strategies.Dispatch{})
					Expect(err).To(Equal(errors.New("BOOM!")))
				})
			})

			Context("when findsUserGUIDs returns an err", func() {
				It("returns an error", func() {
					findsUserGUIDs.UserGUIDsBelongingToSpaceError = errors.New("BOOM!")

					_, err := strategy.Dispatch(strategies.Dispatch{})
					Expect(err).To(Equal(errors.New("BOOM!")))
				})
			})
		})
	})
})
