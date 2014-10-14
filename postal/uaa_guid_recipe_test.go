package postal_test

import (
    "bytes"
    "encoding/json"
    "errors"
    "log"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/fakes"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("UAA Recipe", func() {
    var uaaRecipe postal.UAARecipe

    var fakeCC *fakes.FakeCloudController
    var id int
    var logger *log.Logger
    var fakeUAA fakes.FakeUAAClient
    var mailClient fakes.FakeMailClient
    var buffer *bytes.Buffer
    var options postal.Options
    var tokenLoader postal.TokenLoader
    var userLoader postal.UserLoader
    var spaceLoader postal.SpaceLoader
    var templateLoader postal.TemplateLoader
    var mailer *fakes.FakeMailer
    var fs FakeFileSystem
    var env config.Environment
    var queue *fakes.FakeQueue
    var clientID string
    var fakeReceiptsRepo fakes.FakeReceiptsRepo
    var conn *fakes.FakeDBConn
    var fakeUnsubscribesRepo *fakes.FakeUnsubscribesRepo

    BeforeEach(func() {
        clientID = "mister-client"
        conn = &fakes.FakeDBConn{}

        fakeCC = fakes.NewFakeCloudController()
        fakeCC.UsersBySpaceGuid["space-001"] = []cf.CloudControllerUser{
            cf.CloudControllerUser{Guid: "user-123"},
            cf.CloudControllerUser{Guid: "user-456"},
        }

        fakeCC.Spaces = map[string]cf.CloudControllerSpace{
            "space-001": cf.CloudControllerSpace{
                Name:             "production",
                GUID:             "space-001",
                OrganizationGUID: "org-001",
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
        fakeUAA = fakes.FakeUAAClient{
            ClientToken: uaa.Token{
                Access: fakes.BuildToken(tokenHeader, tokenClaims),
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

        fakeReceiptsRepo = fakes.NewFakeReceiptsRepo()
        fakeUnsubscribesRepo = fakes.NewFakeUnsubscribesRepo()

        buffer = bytes.NewBuffer([]byte{})
        id = 1234
        logger = log.New(buffer, "", 0)
        mailClient = fakes.FakeMailClient{}
        env = config.NewEnvironment()
        fs = NewFakeFileSystem(env)

        queue = fakes.NewFakeQueue()
        mailer = fakes.NewFakeMailer()

        tokenLoader = postal.NewTokenLoader(&fakeUAA)
        userLoader = postal.NewUserLoader(&fakeUAA, logger, fakeCC)
        spaceLoader = postal.NewSpaceLoader(fakeCC)
        templateLoader = postal.NewTemplateLoader(&fs, env.RootPath)

        uaaRecipe = postal.NewUAARecipe(tokenLoader, userLoader, spaceLoader, templateLoader, mailer, &fakeReceiptsRepo)
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
                        _, err := uaaRecipe.Dispatch(clientID, postal.SpaceGUID("space-001"), options, conn)

                        Expect(err).To(BeAssignableToTypeOf(postal.CCDownError("")))
                    })
                })

                Context("when Cloud Controller is unavailable to load a space", func() {
                    It("returns a CCDownError error", func() {
                        fakeCC.LoadSpaceError = errors.New("BOOM!")
                        _, err := uaaRecipe.Dispatch(clientID, postal.SpaceGUID("space-000"), options, conn)

                        Expect(err).To(Equal(errors.New("BOOM!")))
                    })
                })

                Context("when UAA cannot be reached", func() {
                    It("returns a UAADownError", func() {
                        fakeUAA.ErrorForUserByID = uaa.NewFailure(404, []byte("Requested route ('uaa.10.244.0.34.xip.io') does not exist"))
                        _, err := uaaRecipe.Dispatch(clientID, postal.UserGUID("user-123"), options, conn)

                        Expect(err).To(BeAssignableToTypeOf(postal.UAADownError("")))
                    })
                })

                Context("when UAA fails for unknown reasons", func() {
                    It("returns a UAAGenericError", func() {
                        fakeUAA.ErrorForUserByID = errors.New("BOOM!")
                        _, err := uaaRecipe.Dispatch(clientID, postal.UserGUID("user-123"), options, conn)

                        Expect(err).To(BeAssignableToTypeOf(postal.UAAGenericError("")))
                    })
                })

                Context("when a template cannot be loaded", func() {
                    It("returns a TemplateLoadError", func() {
                        delete(fs.Files, env.RootPath+"/templates/user_body.text")

                        _, err := uaaRecipe.Dispatch(clientID, postal.UserGUID("user-123"), options, conn)

                        Expect(err).To(BeAssignableToTypeOf(postal.TemplateLoadError("")))
                    })
                })
            })

            Context("When load Users returns multiple users", func() {
                Context("when create receipts call returns an err", func() {
                    It("returns an error", func() {
                        fakeReceiptsRepo.CreateReceiptsError = true

                        _, err := uaaRecipe.Dispatch(clientID, postal.SpaceGUID("space-001"), options, conn)
                        Expect(err).ToNot(BeNil())
                    })
                })

                It("records a receipt for each user", func() {
                    _, err := uaaRecipe.Dispatch(clientID, postal.SpaceGUID("space-001"), options, conn)
                    if err != nil {
                        panic(err)
                    }

                    Expect(fakeReceiptsRepo.CreateUserGUIDs).To(Equal([]string{"user-123", "user-456"}))
                    Expect(fakeReceiptsRepo.ClientID).To(Equal(clientID))
                    Expect(fakeReceiptsRepo.KindID).To(Equal(options.KindID))
                })

                It("logs the UUIDs of all recipients", func() {
                    _, err := uaaRecipe.Dispatch(clientID, postal.SpaceGUID("space-001"), options, conn)
                    if err != nil {
                        panic(err)
                    }

                    lines := strings.Split(buffer.String(), "\n")

                    Expect(lines).To(ContainElement("CloudController user guid: user-123"))
                    Expect(lines).To(ContainElement("CloudController user guid: user-456"))
                })

                It("calls mailer.Deliver with the correct arguments", func() {
                    _, err := uaaRecipe.Dispatch(clientID, postal.SpaceGUID("space-001"), options, conn)
                    if err != nil {
                        panic(err)
                    }

                    user123 := uaa.User{
                        ID:     "user-123",
                        Emails: []string{"user-123@example.com"},
                    }

                    user456 := uaa.User{
                        ID:     "user-456",
                        Emails: []string{"user-456@example.com"},
                    }

                    users := map[string]uaa.User{"user-123": user123, "user-456": user456}

                    templates := postal.Templates{
                        Subject: "default-missing-subject",
                        Text:    "default-space-text",
                        HTML:    "default-space-html",
                    }

                    Expect(mailer.DeliverArguments).To(ContainElement(conn))
                    Expect(mailer.DeliverArguments).To(ContainElement(templates))
                    Expect(mailer.DeliverArguments).To(ContainElement(users))
                    Expect(mailer.DeliverArguments).To(ContainElement(options))
                    Expect(mailer.DeliverArguments).To(ContainElement(cf.CloudControllerOrganization{
                        Name: "pivotaltracker",
                    }))
                    Expect(mailer.DeliverArguments).To(ContainElement(cf.CloudControllerSpace{
                        GUID:             "space-001",
                        Name:             "production",
                        OrganizationGUID: "org-001",
                    }))
                    Expect(mailer.DeliverArguments).To(ContainElement(clientID))
                })
            })
        })
    })

    Describe("Trim", func() {
        Describe("TrimFields", func() {
            It("trims the specified fields from the response object", func() {
                responses, err := json.Marshal([]postal.Response{
                    {
                        Status:         "delivered",
                        Recipient:      "user-123",
                        NotificationID: "123-456",
                    },
                })

                trimmedResponses := uaaRecipe.Trim(responses)

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
