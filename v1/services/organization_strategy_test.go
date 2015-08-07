package services_test

import (
	"errors"
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Organization Strategy", func() {
	var (
		strategy           services.OrganizationStrategy
		tokenLoader        *fakes.TokenLoader
		organizationLoader *fakes.OrganizationLoader
		enqueuer           *fakes.Enqueuer
		conn               *fakes.Connection
		findsUserGUIDs     *fakes.FindsUserGUIDs
		requestReceived    time.Time
	)

	BeforeEach(func() {
		requestReceived, _ = time.Parse(time.RFC3339Nano, "2015-06-08T14:38:03.180764129-07:00")
		conn = fakes.NewConnection()
		tokenHeader := map[string]interface{}{
			"alg": "FAST",
		}
		tokenClaims := map[string]interface{}{
			"client_id": "mister-client",
			"exp":       int64(3404281214),
			"iss":       "testzone1",
			"scope":     []string{"notifications.write"},
		}
		tokenLoader = fakes.NewTokenLoader()
		tokenLoader.Token = fakes.BuildToken(tokenHeader, tokenClaims)
		enqueuer = fakes.NewEnqueuer()
		findsUserGUIDs = fakes.NewFindsUserGUIDs()
		findsUserGUIDs.OrganizationGuids["org-001"] = []string{"user-123", "user-456"}
		organizationLoader = fakes.NewOrganizationLoader()
		organizationLoader.Organization = cf.CloudControllerOrganization{
			Name: "my-org",
			GUID: "org-001",
		}
		strategy = services.NewOrganizationStrategy(tokenLoader, organizationLoader, findsUserGUIDs, enqueuer)
	})

	Describe("Dispatch", func() {
		Context("when the request is valid", func() {
			It("call enqueuer.Enqueue with the correct arguments for an organization", func() {
				_, err := strategy.Dispatch(services.Dispatch{
					GUID:       "org-001",
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
					UAAHost: "testzone1",
				})
				Expect(err).NotTo(HaveOccurred())

				users := []services.User{
					{GUID: "user-123"},
					{GUID: "user-456"},
				}

				Expect(enqueuer.EnqueueCall.Args.Connection).To(Equal(conn))
				Expect(enqueuer.EnqueueCall.Args.Users).To(Equal(users))
				Expect(enqueuer.EnqueueCall.Args.Options).To(Equal(services.Options{
					ReplyTo:           "reply-to@example.com",
					Subject:           "this is the subject",
					To:                "dr@strangelove.com",
					KindID:            "forgot_password",
					KindDescription:   "Password reminder",
					SourceDescription: "Login system",
					Text:              "Please reset your password by clicking on this link...",
					HTML: services.HTML{
						BodyContent:    "<p>Welcome to the system, now get off my lawn.</p>",
						BodyAttributes: "some-html-body-attributes",
						Head:           "<head></head>",
						Doctype:        "<html>",
					},
					Endorsement: services.OrganizationEndorsement,
				}))
				Expect(enqueuer.EnqueueCall.Args.Space).To(Equal(cf.CloudControllerSpace{}))
				Expect(enqueuer.EnqueueCall.Args.Org).To(Equal(cf.CloudControllerOrganization{
					Name: "my-org",
					GUID: "org-001",
				}))
				Expect(enqueuer.EnqueueCall.Args.Client).To(Equal("mister-client"))
				Expect(enqueuer.EnqueueCall.Args.Scope).To(Equal(""))
				Expect(enqueuer.EnqueueCall.Args.VCAPRequestID).To(Equal("some-vcap-request-id"))
				Expect(enqueuer.EnqueueCall.Args.RequestReceived).To(Equal(requestReceived))
				Expect(enqueuer.EnqueueCall.Args.UAAHost).To(Equal("testzone1"))
				Expect(tokenLoader.LoadArgument).To(Equal("testzone1"))
			})

			Context("when the org role field is set", func() {
				It("calls enqueuer.Enqueue with the correct arguments", func() {
					_, err := strategy.Dispatch(services.Dispatch{
						GUID:       "org-001",
						Role:       "OrgManager",
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
					})
					Expect(err).NotTo(HaveOccurred())

					Expect(enqueuer.EnqueueCall.Args.Options).To(Equal(services.Options{
						ReplyTo:           "reply-to@example.com",
						Subject:           "this is the subject",
						To:                "dr@strangelove.com",
						KindID:            "forgot_password",
						KindDescription:   "Password reminder",
						SourceDescription: "Login system",
						Text:              "Please reset your password by clicking on this link...",
						Role:              "OrgManager",
						HTML: services.HTML{
							BodyContent:    "<p>Welcome to the system, now get off my lawn.</p>",
							BodyAttributes: "some-html-body-attributes",
							Head:           "<head></head>",
							Doctype:        "<html>",
						},
						Endorsement: services.OrganizationRoleEndorsement,
					}))
				})
			})
		})

		Context("failure cases", func() {
			Context("when token loader fails to return a token", func() {
				It("returns an error", func() {
					tokenLoader.LoadError = errors.New("BOOM!")

					_, err := strategy.Dispatch(services.Dispatch{})
					Expect(err).To(Equal(errors.New("BOOM!")))
				})
			})

			Context("when organizationLoader fails to load an organization", func() {
				It("returns the error", func() {
					organizationLoader.LoadError = errors.New("BOOM!")

					_, err := strategy.Dispatch(services.Dispatch{})
					Expect(err).To(Equal(errors.New("BOOM!")))
				})
			})

			Context("when finds user GUIDs returns an error", func() {
				It("returns an error", func() {
					findsUserGUIDs.UserGUIDsBelongingToOrganizationError = errors.New("BOOM!")

					_, err := strategy.Dispatch(services.Dispatch{})
					Expect(err).To(Equal(errors.New("BOOM!")))
				})
			})
		})
	})
})
