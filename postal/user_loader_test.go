package postal_test

import (
    "bytes"
    "log"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/fakes"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("UserLoader", func() {
    var loader postal.UserLoader
    var token string
    var fakeUAAClient fakes.FakeUAAClient
    var fakeCC *fakes.FakeCloudController

    Describe("Load", func() {
        BeforeEach(func() {
            tokenHeader := map[string]interface{}{
                "alg": "FAST",
            }

            tokenClaims := map[string]interface{}{
                "client_id": "mister-client",
                "exp":       int64(3404281214),
                "scope":     []string{"notifications.write"},
            }

            token = fakes.BuildToken(tokenHeader, tokenClaims)

            fakeCC = fakes.NewFakeCloudController()
            fakeCC.UsersBySpaceGuid["space-001"] = []cf.CloudControllerUser{
                cf.CloudControllerUser{Guid: "user-123"},
                cf.CloudControllerUser{Guid: "user-789"},
            }

            fakeCC.UsersByOrganizationGuid["org-001"] = []cf.CloudControllerUser{
                cf.CloudControllerUser{Guid: "user-123"},
                cf.CloudControllerUser{Guid: "user-456"},
                cf.CloudControllerUser{Guid: "user-789"},
                cf.CloudControllerUser{Guid: "user-999"},
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

            fakeUAAClient = fakes.FakeUAAClient{
                ClientToken: uaa.Token{
                    Access: token,
                },
                UsersByID: map[string]uaa.User{
                    "user-123": uaa.User{
                        Emails: []string{"user-123@example.com"},
                        ID:     "user-123",
                    },
                    "user-456": uaa.User{
                        Emails: []string{"user-456@example.com"},
                        ID:     "user-456",
                    },
                    "user-999": uaa.User{
                        Emails: []string{"user-999@example.com"},
                        ID:     "user-999",
                    },
                },
            }

            logger := log.New(bytes.NewBufferString(""), "", 0)
            loader = postal.NewUserLoader(&fakeUAAClient, logger, fakeCC)
        })

        Context("UAA returns a collection of users", func() {
            It("returns a map of users from GUID to uaa.User using a space guid", func() {
                users, err := loader.Load(postal.SpaceGUID("space-001"), token)
                if err != nil {
                    panic(err)
                }

                Expect(len(users)).To(Equal(2))

                user123 := users["user-123"]
                Expect(user123.Emails[0]).To(Equal("user-123@example.com"))
                Expect(user123.ID).To(Equal("user-123"))

                user789, ok := users["user-789"]
                Expect(ok).To(BeTrue())
                Expect(user789).To(Equal(uaa.User{}))
            })

            It("returns a map of users from GUID to uaa.User using an organization guid", func() {
                users, err := loader.Load(postal.OrganizationGUID("org-001"), token)
                if err != nil {
                    panic(err)
                }

                Expect(len(users)).To(Equal(4))

                user123 := users["user-123"]
                Expect(user123.Emails[0]).To(Equal("user-123@example.com"))
                Expect(user123.ID).To(Equal("user-123"))

                user456 := users["user-456"]
                Expect(user456.Emails[0]).To(Equal("user-456@example.com"))
                Expect(user456.ID).To(Equal("user-456"))

                user789, ok := users["user-789"]
                Expect(ok).To(BeTrue())
                Expect(user789).To(Equal(uaa.User{}))

                user999 := users["user-999"]
                Expect(user999.Emails[0]).To(Equal("user-999@example.com"))
                Expect(user999.ID).To(Equal("user-999"))
            })
        })

        Describe("UAA Error Responses", func() {
            Context("when UAA cannot be reached", func() {
                It("returns a UAADownError", func() {
                    fakeUAAClient.ErrorForUserByID = uaa.NewFailure(404, []byte("Requested route ('uaa.10.244.0.34.xip.io') does not exist"))

                    _, err := loader.Load(postal.UserGUID("user-123"), token)

                    Expect(err).To(BeAssignableToTypeOf(postal.UAADownError("")))
                })
            })

            Context("when UAA returns an unknown UAA 404 error", func() {
                It("returns a UAAGenericError", func() {
                    fakeUAAClient.ErrorForUserByID = uaa.NewFailure(404, []byte("Weird message we haven't seen"))

                    _, err := loader.Load(postal.UserGUID("user-123"), token)

                    Expect(err).To(BeAssignableToTypeOf(postal.UAAGenericError("")))
                })
            })

            Context("when UAA returns an failure code that is not 404", func() {
                It("returns a UAADownError", func() {
                    fakeUAAClient.ErrorForUserByID = uaa.NewFailure(500, []byte("Doesn't matter"))

                    _, err := loader.Load(postal.UserGUID("user-123"), token)

                    Expect(err).To(BeAssignableToTypeOf(postal.UAADownError("")))
                })
            })
        })
    })
})
