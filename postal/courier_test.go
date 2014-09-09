package postal_test

import (
    "bytes"
    "errors"
    "log"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Courier", func() {
    var courier postal.Courier
    var fakeCC *FakeCloudController
    var id int
    var logger *log.Logger
    var fakeUAA FakeUAAClient
    var mailClient FakeMailClient
    var buffer *bytes.Buffer
    var options postal.Options
    var tokenLoader postal.TokenLoader
    var userLoader postal.UserLoader
    var spaceLoader postal.SpaceLoader
    var templateLoader postal.TemplateLoader
    var mailer postal.Mailer
    var fs FakeFileSystem
    var env config.Environment
    var queue *FakeQueue
    var clientID string
    var fakeReceiptsRepo FakeReceiptsRepo
    var conn *FakeDBConn
    var fakeUnsubscribesRepo *FakeUnsubscribesRepo

    BeforeEach(func() {
        clientID = "mister-client"
        conn = &FakeDBConn{}

        fakeCC = NewFakeCloudController()
        fakeCC.UsersBySpaceGuid["space-001"] = []cf.CloudControllerUser{
            cf.CloudControllerUser{Guid: "user-123"},
            cf.CloudControllerUser{Guid: "user-456"},
        }

        fakeCC.Spaces = map[string]cf.CloudControllerSpace{
            "space-001": cf.CloudControllerSpace{
                Name:             "production",
                Guid:             "space-001",
                OrganizationGuid: "org-001",
            },
        }

        fakeCC.Orgs = map[string]cf.CloudControllerOrganization{
            "org-001": cf.CloudControllerOrganization{
                Name: "pivotaltracker",
            },
        }

        tokenHeader := map[string]interface{}{
            "alg": "FAST",
        }

        tokenClaims := map[string]interface{}{
            "client_id": "mister-client",
            "exp":       int64(3404281214),
            "scope":     []string{"notifications.write"},
        }
        fakeUAA = FakeUAAClient{
            ClientToken: uaa.Token{
                Access: BuildToken(tokenHeader, tokenClaims),
            },
            UsersByID: map[string]uaa.User{
                "user-123": uaa.User{
                    ID:     "user-123",
                    Emails: []string{"user-123@example.com"},
                },
                "user-456": uaa.User{
                    ID:     "user-456",
                    Emails: []string{"user-456@example.com"},
                },
            },
        }

        fakeReceiptsRepo = NewFakeReceiptsRepo()
        fakeUnsubscribesRepo = NewFakeUnsubscribesRepo()

        buffer = bytes.NewBuffer([]byte{})
        id = 1234
        logger = log.New(buffer, "", 0)
        mailClient = FakeMailClient{}
        env = config.NewEnvironment()
        fs = NewFakeFileSystem(env)

        queue = NewFakeQueue()
        mailer = postal.NewMailer(queue, FakeGuidGenerator, fakeUnsubscribesRepo)

        tokenLoader = postal.NewTokenLoader(&fakeUAA)
        userLoader = postal.NewUserLoader(&fakeUAA, logger, fakeCC)
        spaceLoader = postal.NewSpaceLoader(fakeCC)
        templateLoader = postal.NewTemplateLoader(&fs)

        courier = postal.NewCourier(tokenLoader, userLoader, spaceLoader, templateLoader, mailer, &fakeReceiptsRepo)
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

            Context("failure cases", func() {
                Context("when Cloud Controller is unavailable to load space users", func() {
                    It("returns a CCDownError error", func() {
                        fakeCC.GetUsersBySpaceGuidError = errors.New("BOOM!")
                        _, err := courier.Dispatch(clientID, postal.SpaceGUID("space-001"), options, conn)

                        Expect(err).To(BeAssignableToTypeOf(postal.CCDownError("")))
                    })
                })

                Context("when Cloud Controller is unavailable to load a space", func() {
                    It("returns a CCDownError error", func() {
                        fakeCC.LoadSpaceError = errors.New("BOOM!")
                        _, err := courier.Dispatch(clientID, postal.SpaceGUID("space-000"), options, conn)

                        Expect(err).To(Equal(errors.New("BOOM!")))
                    })
                })

                Context("when UAA cannot be reached", func() {
                    It("returns a UAADownError", func() {
                        fakeUAA.ErrorForUserByID = uaa.NewFailure(404, []byte("Requested route ('uaa.10.244.0.34.xip.io') does not exist"))
                        _, err := courier.Dispatch(clientID, postal.UserGUID("user-123"), options, conn)

                        Expect(err).To(BeAssignableToTypeOf(postal.UAADownError("")))
                    })
                })

                Context("when UAA fails for unknown reasons", func() {
                    It("returns a UAAGenericError", func() {
                        fakeUAA.ErrorForUserByID = errors.New("BOOM!")
                        _, err := courier.Dispatch(clientID, postal.UserGUID("user-123"), options, conn)

                        Expect(err).To(BeAssignableToTypeOf(postal.UAAGenericError("")))
                    })
                })

                Context("when a template cannot be loaded", func() {
                    It("returns a TemplateLoadError", func() {
                        delete(fs.Files, env.RootPath+"/templates/user_body.text")

                        _, err := courier.Dispatch(clientID, postal.UserGUID("user-123"), options, conn)

                        Expect(err).To(BeAssignableToTypeOf(postal.TemplateLoadError("")))
                    })
                })
            })

            Context("When load Users returns multiple users", func() {
                Context("when create receipts call returns an err", func() {
                    It("returns an error", func() {
                        fakeReceiptsRepo.CreateReceiptsError = true

                        _, err := courier.Dispatch(clientID, postal.SpaceGUID("space-001"), options, conn)
                        Expect(err).ToNot(BeNil())
                    })
                })

                It("records a receipt for each user", func() {
                    _, err := courier.Dispatch(clientID, postal.SpaceGUID("space-001"), options, conn)
                    if err != nil {
                        panic(err)
                    }

                    Expect(fakeReceiptsRepo.CreateUserGUIDs).To(Equal([]string{"user-123", "user-456"}))
                    Expect(fakeReceiptsRepo.ClientID).To(Equal(clientID))
                    Expect(fakeReceiptsRepo.KindID).To(Equal(options.KindID))
                })

                It("logs the UUIDs of all recipients", func() {
                    _, err := courier.Dispatch(clientID, postal.SpaceGUID("space-001"), options, conn)
                    if err != nil {
                        panic(err)
                    }

                    lines := strings.Split(buffer.String(), "\n")

                    Expect(lines).To(ContainElement("CloudController user guid: user-123"))
                    Expect(lines).To(ContainElement("CloudController user guid: user-456"))
                })

                It("returns necessary info in the response for the sent mail", func() {
                    courier = postal.NewCourier(tokenLoader, userLoader, spaceLoader, templateLoader, mailer, &fakeReceiptsRepo)
                    responses, err := courier.Dispatch(clientID, postal.SpaceGUID("space-001"), options, conn)
                    if err != nil {
                        panic(err)
                    }

                    Expect(len(responses)).To(Equal(2))
                    Expect(responses).To(ContainElement(postal.Response{
                        Recipient:      "user-123",
                        Status:         "queued",
                        NotificationID: "deadbeef-aabb-ccdd-eeff-001122334455",
                        Email:          "user-123@example.com",
                    }))

                    Expect(responses).To(ContainElement(postal.Response{
                        Recipient:      "user-456",
                        Status:         "queued",
                        NotificationID: "deadbeef-aabb-ccdd-eeff-001122334455",
                        Email:          "user-456@example.com",
                    }))
                })
            })
        })
    })
})
