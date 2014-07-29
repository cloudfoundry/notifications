package postal_test

import (
    "bytes"
    "log"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("UserLoader", func() {
    var loader postal.UserLoader
    var token string
    var fakeUAAClient FakeUAAClient
    var fakeCC *FakeCloudController

    Describe("Load", func() {
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

            fakeCC = NewFakeCloudController()
            fakeCC.UsersBySpaceGuid["space-001"] = []cf.CloudControllerUser{
                cf.CloudControllerUser{Guid: "user-123"},
                cf.CloudControllerUser{Guid: "user-789"},
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

            fakeUAAClient = FakeUAAClient{
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

            logger := log.New(bytes.NewBufferString(""), "", 0)
            loader = postal.NewUserLoader(&fakeUAAClient, logger, fakeCC)
        })

        Context("UAA returns a collection of users", func() {
            It("returns a map of users from GUID to uaa.User", func() {
                users, err := loader.Load(postal.IsSpace, "space-001", token)
                if err != nil {
                    panic(err)
                }

                Expect(len(users)).To(Equal(2))

                user123 := users["user-123"]
                Expect(user123.ID).To(Equal("user-123"))
                Expect(user123.Emails[0]).To(Equal("user-123@example.com"))

                user789, ok := users["user-789"]
                Expect(ok).To(BeTrue())
                Expect(user789).To(Equal(uaa.User{}))
            })
        })

        Describe("UAA Error Responses", func() {
            Context("when UAA cannot be reached", func() {
                It("returns a UAADownError", func() {
                    fakeUAAClient.ErrorForUserByID = uaa.NewFailure(404, []byte("Requested route ('uaa.10.244.0.34.xip.io') does not exist"))

                    _, err := loader.Load(postal.IsUser, "user-123", token)

                    Expect(err).To(BeAssignableToTypeOf(postal.UAADownError("")))
                })
            })

            Context("when UAA returns an unknown UAA 404 error", func() {
                It("returns a UAAGenericError", func() {
                    fakeUAAClient.ErrorForUserByID = uaa.NewFailure(404, []byte("Weird message we haven't seen"))

                    _, err := loader.Load(postal.IsUser, "user-123", token)

                    Expect(err).To(BeAssignableToTypeOf(postal.UAAGenericError("")))
                })
            })

            Context("when UAA returns an failure code that is not 404", func() {
                It("returns a UAADownError", func() {
                    fakeUAAClient.ErrorForUserByID = uaa.NewFailure(500, []byte("Doesn't matter"))

                    _, err := loader.Load(postal.IsUser, "user-123", token)

                    Expect(err).To(BeAssignableToTypeOf(postal.UAADownError("")))
                })
            })
        })
    })
})
