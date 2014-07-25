package handlers_test

import (
    "bytes"
    "log"
    "net/http/httptest"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("NotifyHelper", func() {
    var helper handlers.NotifyHelper
    var fakeCC *FakeCloudController
    var logger *log.Logger
    var fakeUAA FakeUAAClient
    var mailClient FakeMailClient
    var writer *httptest.ResponseRecorder
    var token string
    var buffer *bytes.Buffer

    Describe("LoadUaaUser", func() {
        BeforeEach(func() {

            tokenHeader := map[string]interface{}{
                "alg": "FAST",
            }

            tokenClaims := map[string]interface{}{
                "client_id": "mister-client",
                "exp":       3404281214,
                "scope":     []string{"notifications.write"},
            }

            token = BuildToken(tokenHeader, tokenClaims)

            writer = httptest.NewRecorder()

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

            fakeUAA = FakeUAAClient{
                ClientToken: uaa.Token{
                    Access: token,
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

            buffer = bytes.NewBuffer([]byte{})
            logger = log.New(buffer, "", 0)

            mailClient = FakeMailClient{}

            helper = handlers.NewNotifyHelper(fakeCC, logger, &fakeUAA, FakeGuidGenerator, &mailClient)
        })

        Context("UAA returns a user", func() {
            It("returns the uaa.User", func() {
                user, err := helper.LoadUaaUser("user-123", fakeUAA)
                if err != nil {
                    panic(err)
                }

                Expect(user.ID).To(Equal("user-123"))
                Expect(user.Emails[0]).To(Equal("user-123@example.com"))
            })
        })

        Describe("UAA Error Responses", func() {
            Context("when UAA cannot be reached", func() {
                It("returns a UAADownError", func() {
                    fakeUAA.ErrorForUserByID = uaa.NewFailure(404, []byte("Requested route ('uaa.10.244.0.34.xip.io') does not exist"))

                    _, err := helper.LoadUaaUser("user-123", fakeUAA)

                    Expect(err).To(BeAssignableToTypeOf(handlers.UAADownError{}))
                })
            })

            Context("when UAA cannot find the user", func() {
                It("returns a UAAUserNotFoundError", func() {
                    fakeUAA.ErrorForUserByID = uaa.NewFailure(404, []byte("User f3b51aac-866e-4b7a-948c-de31beefc475d does not exist"))

                    _, err := helper.LoadUaaUser("user-123", fakeUAA)

                    Expect(err).To(BeAssignableToTypeOf(handlers.UAAUserNotFoundError{}))
                })
            })

            Context("when UAA returns an unknown UAA 404 error", func() {
                It("returns a UAAGenericError", func() {
                    fakeUAA.ErrorForUserByID = uaa.NewFailure(404, []byte("Weird message we haven't seen"))

                    _, err := helper.LoadUaaUser("user-123", fakeUAA)

                    Expect(err).To(BeAssignableToTypeOf(handlers.UAAGenericError{}))
                })
            })

            Context("when UAA returns an failure code that is not 404", func() {
                It("returns a UAADownError", func() {
                    fakeUAA.ErrorForUserByID = uaa.NewFailure(500, []byte("Doesn't matter"))

                    _, err := helper.LoadUaaUser("user-123", fakeUAA)

                    Expect(err).To(BeAssignableToTypeOf(handlers.UAADownError{}))
                })
            })
        })
    })
})
